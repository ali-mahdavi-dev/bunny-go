package command

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	config "github.com/ali-mahdavi-dev/bunny-go/configs"
	"github.com/ali-mahdavi-dev/bunny-go/docs"
	"github.com/ali-mahdavi-dev/bunny-go/internal/account"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/databases"
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

	db, err := databases.New(databases.Config{
		Debug:        cfg.Debug,
		DBType:       cfg.Database.Type,
		DSN:          cfg.Database.Dns,
		MaxLifetime:  cfg.Database.MaxLifeTime,
		MaxIdleTime:  cfg.Database.MaxIdleTime,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxOpenConns: cfg.Database.MaxOpenConns,
	})
	if err != nil {
		panic(err)
	}

	server := gin.Default()
	registerSwagger(server, cfg)

	// Bootstrap
	account.Bootstrap(server, db)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
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
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.Host)
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
