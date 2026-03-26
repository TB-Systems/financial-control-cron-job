package main

import (
	"auto-service/internal/clients"
	"auto-service/internal/services"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Loads environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: failed to load .env file: %v", err)
		return
	}

	// Get retry max count from environment variable or use default
	retryMaxCount := getRetryMaxCount()

	// Initialize clients and services for user
	userClient := clients.NewUserClient()
	userService := services.NewUserService(userClient)

	// Initialize clients and services for engine
	engineClient := clients.NewEngineClient()
	engineService := services.NewEngineService(engineClient)

	// Set up cron job to run at 2:00 AM on the first day of every month on BRT timezone
	loc, _ := time.LoadLocation("America/Sao_Paulo")

	// Define the cron job function
	c := cron.New(
		cron.WithLocation(loc),
		cron.WithSeconds(),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		),
	)

	// Define the job function to be executed by the cron scheduler
	runJob := func() {
		// Get the total page count for users
		count, err := getCountWithRetries(userService, retryMaxCount)

		if err != nil {
			log.Printf("job error: GetPageCount failed after %d attempts", retryMaxCount)
			return
		}

		// Process each page of users
		for step := 1; step <= count; step++ {
			// Get paginated users with retries
			usersResponse, err := getUsersPaginatedWithRetries(userService, retryMaxCount)
			if err != nil {
				log.Printf("job error: GetUsersPaginated failed after %d attempts", retryMaxCount)
				return
			}

			// Process users in parallel for the current page
			processUsersInParallel(usersResponse, engineService, retryMaxCount)
		}
	}

	// Execute the job immediately on startup, then schedule it to run monthly
	runJob()

	// Schedule the job to run at 2:00 AM on the first day of every month
	_, _ = c.AddFunc("0 0 2 1 * *", runJob)

	// Start the cron scheduler
	c.Start()

	// Keep the main function running indefinitely
	select {}
}

// getRetryMaxCount reads the RETRY_MAX_COUNT environment variable and returns it as an integer.
// If the variable is not set or is invalid, it returns a default value of 4.
func getRetryMaxCount() int {
	// Default retry attempts if environment variable is not set or invalid
	const defaultAttempts = 4

	// Read the RETRY_MAX_COUNT environment variable
	value := os.Getenv("RETRY_MAX_COUNT")

	// If the environment variable is not set, return the default value
	if value == "" {
		return defaultAttempts
	}

	// Convert the environment variable value to an integer
	attempts, err := strconv.Atoi(value)

	// If there was an error converting the value or if the value is less than 1, log a warning and return the default value
	if err != nil || attempts < 1 {
		log.Printf("warning: invalid RETRY_MAX_COUNT=%q, using default=%d", value, defaultAttempts)
		return defaultAttempts
	}

	return attempts
}

func getCountWithRetries(service services.User, retryMaxCount int) (int, error) {
	// Initialize variables to hold the count and error
	var (
		count int
		err   error
	)

	// Attempt to get the page count with retries
	for attempt := 1; attempt <= retryMaxCount; attempt++ {

		// Try to get the page count
		count, err = service.GetPageCount()

		// If there was no error, break out of the loop
		if err == nil {
			break
		}

		// Log the error for this attempt
		log.Printf("job error on GetPageCount attempt %d/%d: %v", attempt, retryMaxCount, err)
	}

	return count, err
}

func getUsersPaginatedWithRetries(service services.User, retryMaxCount int) (commonsmodels.PaginatedResponse[uuid.UUID], error) {
	// Initialize variables to hold the response and error
	var (
		response commonsmodels.PaginatedResponse[uuid.UUID]
		err      error
	)

	// Attempt to get the paginated users with retries
	for attempt := 1; attempt <= retryMaxCount; attempt++ {

		// Try to get the paginated users
		response, err = service.GetUsersPaginated()

		// If there was no error, break out of the loop
		if err == nil {
			break
		}

		// Log the error for this attempt
		log.Printf("job error on GetUsersPaginated attempt %d/%d: %v", attempt, retryMaxCount, err)
	}

	return response, err
}

func processUsersInParallel(usersResponse commonsmodels.PaginatedResponse[uuid.UUID], engineService services.Engine, retryMaxCount int) {
	// If there are no users to process, return early
	if len(usersResponse.Items) == 0 {
		return
	}

	// Use a WaitGroup to wait for all goroutines to finish and a channel to collect errors
	var wg sync.WaitGroup

	// Create a buffered channel to collect errors from goroutines
	errCh := make(chan error, len(usersResponse.Items))

	// Process each user in parallel
	for _, userID := range usersResponse.Items {
		// Increment the WaitGroup counter for each user being processed
		wg.Add(1)

		// Start a new goroutine to process the user
		go func(id uuid.UUID) {
			// Decrement the WaitGroup counter when the goroutine completes
			defer wg.Done()

			// Process the user block and send any errors to the error channel
			if err := processUserBlock(id, engineService, retryMaxCount); err != nil {
				errCh <- err
			}
		}(userID)
	}

	// Wait for all goroutines to finish and then close the error channel
	wg.Wait()

	// Close the error channel after all goroutines have finished
	close(errCh)

	// Log any errors that were collected from the goroutines
	for err := range errCh {
		log.Printf("job error while processing user in parallel: %v", err)
	}
}

func processUserBlock(userID uuid.UUID, engineService services.Engine, retryMaxCount int) error {
	// Initialize variables to hold the response lists and error
	var (
		monthlyIDs      commonsmodels.ResponseList[uuid.UUID]
		annualIDs       commonsmodels.ResponseList[uuid.UUID]
		installmentsIDs commonsmodels.ResponseList[uuid.UUID]
		err             error
	)

	// Attempt to get monthly transactions with retries
	for attempt := 1; attempt <= retryMaxCount; attempt++ {
		monthlyIDs, err = engineService.GetMonthlyTransactions()
		if err == nil {
			break
		}

		log.Printf("job error on GetMonthlyTransactions attempt %d/%d: %v", attempt, retryMaxCount, err)
	}

	// Process monthly transactions with retries
	for _, monthlyID := range monthlyIDs.Items {
		for attempt := 1; attempt <= retryMaxCount; attempt++ {
			err = engineService.PostMonthlyTransaction(monthlyID)
			if err == nil {
				break
			}

			log.Printf("job error on PostMonthlyTransaction attempt %d/%d: %v", attempt, retryMaxCount, err)
		}
	}

	// Attempt to get annual transactions with retries
	for attempt := 1; attempt <= retryMaxCount; attempt++ {
		annualIDs, err = engineService.GetAnnualTransactions()
		if err == nil {
			break
		}

		log.Printf("job error on GetAnnualTransactions attempt %d/%d: %v", attempt, retryMaxCount, err)
	}

	// Process annual transactions with retries
	for _, annualID := range annualIDs.Items {
		for attempt := 1; attempt <= retryMaxCount; attempt++ {
			err = engineService.PostAnnualTransaction(annualID)
			if err == nil {
				break
			}

			log.Printf("job error on PostAnnualTransaction attempt %d/%d: %v", attempt, retryMaxCount, err)
		}
	}

	// Attempt to get installments transactions with retries
	for attempt := 1; attempt <= retryMaxCount; attempt++ {
		installmentsIDs, err = engineService.GetInstallmentsTransactions()
		if err == nil {
			break
		}

		log.Printf("job error on GetInstallmentsTransactions attempt %d/%d: %v", attempt, retryMaxCount, err)
	}

	// Process installments transactions with retries
	for _, installmentsID := range installmentsIDs.Items {
		for attempt := 1; attempt <= retryMaxCount; attempt++ {
			err = engineService.PostInstallmentsTransaction(installmentsID)
			if err == nil {
				break
			}

			log.Printf("job error on PostInstallmentsTransaction attempt %d/%d: %v", attempt, retryMaxCount, err)
		}
	}

	// If we reach this point, it means all transactions for the user were processed (or attempted) and we can log the success
	log.Printf("job success: finished processing user %s", userID)
	return nil
}
