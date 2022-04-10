package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/api"
	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/database/dbuser"
	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/service"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	configPath := flag.String("c", "config.yaml", "path to your config")
	flag.Parse()

	conf, err := readConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	dbUser, err := dbuser.NewDatabase(ctx, conf.DBUser)
	if err != nil {
		log.Fatalf("cannot create db: %v", err)
	}

	defer dbUser.Close(ctx)

	err = dbUser.Ping(ctx)
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	serverParams := api.ServerParams{
		Cfg:  conf.Server,
		Srvc: service.NewService(dbUser),
	}
	srv := api.NewServer(&serverParams)

	rungroup, ctx := errgroup.WithContext(ctx)

	log.Printf("starting server on address: %s", srv.Addr)
	rungroup.Go(func() error {
		if er := srv.ListenAndServe(); er != nil && !errors.Is(er, http.ErrServerClosed) {
			return fmt.Errorf("listen and server %w", er)
		}

		return nil
	})

	rungroup.Go(func() error {
		<-ctx.Done()

		if er := srv.Shutdown(context.Background()); er != nil {
			return fmt.Errorf("shutdown http server %w", er)
		}

		return nil
	})

	err = rungroup.Wait()
	if err != nil {
		log.Printf("run group exited because of error: %v", err)
		return
	}

	log.Println("Server Exited Properly")
}
