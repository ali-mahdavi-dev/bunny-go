package commands

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	config "github.com/ali-mahdavi-dev/bunny-go/config"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account"
	mwF "github.com/ali-mahdavi-dev/bunny-go/internal/framework/api/middleware"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/databases"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/unit_of_work"
)

func runHTTPServerCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "start http server",
		RunE: func(_ *cobra.Command, _ []string) error {
			initializeConfigs()

			log.Println("starting http server")

			return startServer(&cfg)
		},
	}
}

func startServer(cfg *config.Config) error {
	db, err := databases.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	uow := unit_of_work.New(db, LogInstans)

	server := fiber.New()

	// Middleware
	middlewareF := mwF.NewMiddleware(cfg, uow)
	middlewareF.Register(server)

	// Swagger
	registerSwagger(server, cfg)

	// Metrics (Prometheus)
	server.Get("/metrics", func(c *fiber.Ctx) error {
		// Convert http.Handler to fasthttp.Handler
		h := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
		h(c.Context())
		return nil
	})

	// Bootstrap application routes
	account.Bootstrap(server, db, cfg, LogInstans)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Domain, cfg.Server.InternalPort)
	if err := server.Listen(addr); err != nil {
		panic(err)
	}

	return nil
}

func registerSwagger(app *fiber.App, cfg *config.Config) {

	swCfg := swagger.Config{
		Title:    "golang web api",
		BasePath: fmt.Sprintf("localhost:%s/api", cfg.Server.ExternalPort),
		FilePath: "docs/swagger.json",
	}

	app.Use(swagger.New(swCfg))
}
