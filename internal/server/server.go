package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bombsimon/logrusr/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/shemanaev/updtr/internal/config"
	"github.com/shemanaev/updtr/internal/container"
	"github.com/shemanaev/updtr/internal/meta"
	"github.com/shemanaev/updtr/internal/static"
	mw "github.com/shemanaev/updtr/pkg/middleware"
)

func RunServer() error {
	cfg := config.LoadConfig()

	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if cfg.Trace {
		log.SetLevel(log.TraceLevel)
	}

	mapping, err := config.LoadMappings("/data/mapping.yml")
	if err != nil {
		return err
	}
	log.Infof("Loaded mappings: %d", len(mapping))
	log.Trace(mapping)

	logger := log.New()

	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)

	cacheControlDuration, _ := time.ParseDuration("24h")
	staticFiles := mw.FSCacheControl(static.Files, meta.BuildDate, cacheControlDuration)
	r.Handle("/static/*", http.StripPrefix("/static/", staticFiles))

	client := container.NewClient(cfg, mapping)

	web := NewHomeHandler(client)
	r.Get("/", web.Home)

	api := NewApiHandler(client)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/containers", api.ListContainers)
		r.Post("/containers/update", api.UpdateContainers)
	})
	r.Get("/v1/update", api.RefreshContainers)

	var scheduleSpec string
	if len(cfg.Schedule) > 0 {
		scheduleSpec = cfg.Schedule
	} else {
		scheduleSpec = fmt.Sprintf("@every %ds", cfg.PollInterval)
	}

	scheduler := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(logrusr.New(logger)),
	))
	_, err = scheduler.AddFunc("@every 3s", func() { client.DoUpdate() })
	if err != nil {
		return err
	}
	_, err = scheduler.AddFunc(scheduleSpec, func() { client.RefreshContainers() })
	if err != nil {
		return err
	}

	scheduler.Start()
	defer scheduler.Stop()

	return http.ListenAndServe(":8080", r)
}
