package command

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	config "github.com/ali-mahdavi-dev/bunny-go/configs"
)

var (
	cfg     config.Config
	envFile string
	rootCmd = &cobra.Command{
		Use: "",
		Run: func(cmd *cobra.Command, args []string) {
			initializeConfigs()
		},
	}
)

func initializeConfigs() {
	err := godotenv.Load(envFile)
	if err != nil {
		panic(err)
	}

	c, err := config.Load()
	if err != nil {
		log.Fatalf("could not load configuration %s\n", err.Error())
	}

	cfg = *c
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
