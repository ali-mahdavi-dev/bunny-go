package command

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	config "github.com/ali-mahdavi-dev/bunny-go/config"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/logging"
)

var (
	cfg        config.Config
	envFile    string
	LogInstans logging.Logger
	rootCmd    = &cobra.Command{
		Use: "",
		Run: func(cmd *cobra.Command, args []string) {
			initializeConfigs()
		},
	}
)

func initializeConfigs() {
	cfg = *config.GetConfig()
	LogInstans = logging.NewLogger(&cfg)
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVarP(&envFile, "env-file", "e", ".env", ".env file")

	rootCmd.AddCommand(runHTTPServerCMD())
	rootCmd.AddCommand(migrateCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
