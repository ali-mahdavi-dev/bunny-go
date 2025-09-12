package command

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	config "bunny-go/configs"
	"bunny-go/internal/framwork/infrastructure/databases"
	"bunny-go/internal/user_management"
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
		TablePrefix:  cfg.Database.TablePrefox,
	})
	if err != nil {
		panic(err)
	}

	server := gin.Default()

	// Bootstrap
	user_management.Bootstrap(server, db)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	err = server.Run(addr)
	if err != nil {
		panic(err)
	}

	return nil
}
