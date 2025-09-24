package command

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	config "github.com/ali-mahdavi-dev/bunny-go/config"
	"github.com/ali-mahdavi-dev/bunny-go/docs"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/databases"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/websocket"
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
	serverWs := socketio.NewServer(nil)

	ws := websocket.NewWebsocket(serverWs, LogInstans, cfg)
	ws.AddWsRoutes()

	go func() {
		if err := serverWs.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer serverWs.Close()

	db, err := databases.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	server := gin.Default()

	// init ws
	// برای upgrade کردن به websocket
	server.GET("/socket.io/*any", gin.WrapH(serverWs))
	server.POST("/socket.io/*any", gin.WrapH(serverWs))

	// middleware
	// server.Use(middleware.DefaultStructuredLogger(cfg))

	// swagger
	registerSwagger(server, cfg)

	// metrics
	server.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Bootstrap
	account.Bootstrap(server, db, cfg, LogInstans)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Domain, cfg.Server.InternalPort)
	err = server.Run(addr)
	if err != nil {
		panic(err)
	}

	return nil
}

func registerSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "golang web api"
	docs.SwaggerInfo.Description = "golang web api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.ExternalPort)
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
