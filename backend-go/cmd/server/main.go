package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/instaagrammeta/crm-real/backend-go/internal/config"
	"github.com/instaagrammeta/crm-real/backend-go/internal/database"
	"github.com/instaagrammeta/crm-real/backend-go/internal/logger"
	"github.com/instaagrammeta/crm-real/backend-go/internal/router"
	"github.com/instaagrammeta/crm-real/backend-go/internal/services"
	"github.com/instaagrammeta/crm-real/backend-go/internal/services/seed"
	"github.com/instaagrammeta/crm-real/backend-go/internal/services/telegram"
	"github.com/instaagrammeta/crm-real/backend-go/internal/ws"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.AppEnv)

	// ensure upload dirs exist
	for _, sub := range []string{"photos", "files", "chat", "posts"} {
		if err := os.MkdirAll(filepath.Join(cfg.UploadDir, sub), 0o755); err != nil {
			log.Fatal().Err(err).Msg("create upload dir")
		}
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("init database")
	}

	if err := seed.EnsureAdmin(db, cfg); err != nil {
		log.Error().Err(err).Msg("seed admin")
	}

	hub := ws.NewHub()
	go hub.Run()

	tgBot := telegram.NewBot(cfg, db)
	if err := tgBot.Start(); err != nil {
		log.Warn().Err(err).Msg("telegram bot disabled")
	}

	notifier := services.NewNotifier(db, hub, tgBot)

	r := router.New(cfg, db, hub, notifier, tgBot)

	srv := &http.Server{
		Addr:              cfg.AppHost + ":" + cfg.AppPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Msgf("HTTP server listening on http://%s:%s", cfg.AppHost, cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("http server")
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info().Msg("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown")
	}
	tgBot.Stop()
	log.Info().Msg("bye")
}
