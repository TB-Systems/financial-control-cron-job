package main

import (
	"context"
)

func main() {
	// loc, _ := time.LoadLocation("America/Sao_Paulo")

	// c := cron.New(
	// 	cron.WithLocation(loc),
	// 	cron.WithSeconds(),
	// 	cron.WithChain(
	// 		cron.SkipIfStillRunning(cron.DefaultLogger), // evita overlap do mesmo job
	// 		cron.Recover(cron.DefaultLogger),
	// 	),
	// )

	// _, _ = c.AddFunc("0 0 2 * * *", func() {
	// 	if err := runNightJob(context.Background()); err != nil {
	// 		log.Printf("job erro: %v", err)
	// 	}
	// })

	// c.Start()
	// select {}
}

func runNightJob(ctx context.Context) error {
	// // I/O bound: comece com algo entre 16-64 e ajuste por métrica
	// workers := max(16, runtime.NumCPU()*4)

	// g, ctx := errgroup.WithContext(ctx)
	// usersCh := make(chan User, workers*2)

	// // produtor: busca usuários paginados
	// g.Go(func() error {
	// 	defer close(usersCh)
	// 	page := 1
	// 	for {
	// 		users, err := fetchUsersPage(ctx, page, 500) // lote de 500 (ajuste)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if len(users) == 0 {
	// 			return nil
	// 		}
	// 		for _, u := range users {
	// 			select {
	// 			case usersCh <- u:
	// 			case <-ctx.Done():
	// 				return ctx.Err()
	// 			}
	// 		}
	// 		page++
	// 	}
	// })

	// // consumidores
	// for i := 0; i < workers; i++ {
	// 	g.Go(func() error {
	// 		for u := range usersCh {
	// 			// paraleliza serviços por usuário
	// 			ug, uctx := errgroup.WithContext(ctx)
	// 			ug.Go(func() error { return callServiceA(uctx, u) })
	// 			ug.Go(func() error { return callServiceB(uctx, u) })
	// 			ug.Go(func() error { return callServiceC(uctx, u) })
	// 			if err := ug.Wait(); err != nil {
	// 				// política: logar e seguir ou abortar
	// 				log.Printf("user %s falhou: %v", u.ID, err)
	// 			}
	// 		}
	// 		return nil
	// 	})
	// }

	// return g.Wait()
	return nil
}
