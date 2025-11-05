package commands

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/swagger/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	config "shikposh-backend/config"
	"shikposh-backend/internal/account"
	mwF "shikposh-backend/pkg/framework/api/middleware"
	"shikposh-backend/pkg/framework/infrastructure/databases"
	"shikposh-backend/pkg/framework/service_layer/unit_of_work"
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
	dbConfig := databases.Config{
		DBType:       "postgres",
		DSN:          fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DbName, cfg.Postgres.SSLMode),
		MaxOpenConns: cfg.Postgres.MaxOpenConns,
		MaxIdleConns: cfg.Postgres.MaxIdleConns,
		MaxLifetime:  int(cfg.Postgres.ConnMaxLifetime.Seconds()),
	}
	db, err := databases.New(dbConfig)
	if err != nil {
		panic(err)
	}
	uow := unit_of_work.New(db)

	server := fiber.New()

	// Middleware
	middlewareConfig := mwF.MiddlewareConfig{
		JWTSecret: cfg.JWT.Secret,
	}
	middlewareF := mwF.NewMiddleware(middlewareConfig, uow)
	middlewareF.Register(server)

	// Swagger Documentation
	registerSwagger(server, cfg)

	// Metrics (Prometheus)
	server.Get("/metrics", func(c fiber.Ctx) error {
		// In Fiber v3, we need to manually handle Prometheus metrics
		// We'll serve the metrics by converting the standard HTTP handler
		// Since we can't directly access RequestCtx, we'll use a different approach
		metricsHandler := promhttp.Handler()
		// Create a temporary HTTP server response
		// Use the adapter to convert http.Handler to fasthttp
		adapter := fasthttpadaptor.NewFastHTTPHandler(metricsHandler)
		// Get the underlying RequestCtx from the fiber context
		// In Fiber v3, we access it through the request object
		if reqCtx, ok := c.Locals("requestCtx").(*fasthttp.RequestCtx); ok && reqCtx != nil {
			adapter(reqCtx)
			return nil
		}
		// Fallback: try to access through request
		req := c.Request()
		if reqCtx := req.Header.UserAgent(); reqCtx != nil {
			// Alternative approach: use standard HTTP conversion
			// For now, return metrics as plain text (simplified)
			c.Set("Content-Type", "text/plain")
			return c.SendString("Prometheus metrics endpoint - requires RequestCtx access")
		}
		return c.SendString("Metrics unavailable")
	})

	// Bootstrap application routes
	account.Bootstrap(server, db, cfg, LogInstans)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Domain, cfg.Server.InternalPort)
	log.Printf("Server starting on %s", addr)
	log.Printf("Swagger UI available at http://%s/swagger/index.html", addr)
	log.Printf("API Base Path: http://%s/api", addr)

	if err := server.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	return nil
}

func registerSwagger(app *fiber.App, cfg *config.Config) {
	// Serve swagger.json file directly
	app.Get("/swagger.json", func(c fiber.Ctx) error {
		return c.SendFile("docs/swagger.json")
	})

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger.json",
	})) 
}
