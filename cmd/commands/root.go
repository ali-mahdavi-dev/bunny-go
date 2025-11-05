package commands

import (
	"os"

	"github.com/spf13/cobra"

	config "shikposh-backend/config"
	"shikposh-backend/pkg/framework/infrastructure/logging"
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
	loggerConfig := logging.LoggerConfig{
		Type:   logging.LoggerTypeZerolog,
		Level:  logging.LogLevel(cfg.Logger.Level),
		Format: logging.LogFormatJSON,
	}
	var err error
	LogInstans, err = logging.NewLogger(loggerConfig)
	if err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVarP(&envFile, "env-file", "e", ".env", ".env file")

	rootCmd.AddCommand(runHTTPServerCMD())
	rootCmd.AddCommand(migrateCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Logger might not be initialized, check if available
		if logging.GetLogger() != nil {
			logging.Error("Command execution failed").WithError(err).Log()
		}
		os.Exit(1)
	}
}
