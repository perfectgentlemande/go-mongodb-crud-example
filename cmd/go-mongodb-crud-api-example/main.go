package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/api"
	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/database/dbuser"
	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/logger"
	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	log := logger.DefaultLogger()

	configPath := flag.String("c", "config.yaml", "path to your config")
	flag.Parse()

	conf, err := readConfig(*configPath)
	if err != nil {
		log.Fatal("failed to read config", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		}, zapcore.Field{
			Key:    "config_path",
			Type:   zapcore.StringType,
			String: *configPath,
		})
	}

	dbUser, err := dbuser.NewDatabase(ctx, conf.DBUser)
	if err != nil {
		log.Fatal("cannot create db", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
	}

	defer dbUser.Close(ctx)

	err = dbUser.Ping(ctx)
	if err != nil {
		log.Fatal("cannot ping db", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		}, zapcore.Field{
			Key:    "conn_string",
			Type:   zapcore.StringType,
			String: conf.DBUser.ConnStr,
		})
	}

	serverParams := api.ServerParams{
		Cfg:  conf.Server,
		Srvc: service.NewService(dbUser),
		Log:  log,
	}
	srv := api.NewServer(&serverParams)

	rungroup, ctx := errgroup.WithContext(ctx)

	log.Info("starting server", zap.Field{
		Key:    "address",
		Type:   zapcore.ErrorType,
		String: srv.Addr,
	})
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
		log.Error("run group exited because of error", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
		return
	}

	log.Info("server exited properly")
}
