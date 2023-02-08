package cmd

import (
	"keycloakslackbot/logs"
	"keycloakslackbot/proc"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const VERSION = "0.1.0"

var (
	SLACK_URL         string
	KEYCLOAK_URL      string
	KEYCLOAK_USER     string
	KEYCLOAK_PASSWORD string
	INTERVAL          string
	KEYCLOAK_REALM    string
)

var rootCmd = &cobra.Command{
	Version: VERSION,
	Use:     "",
	Short:   "",
	Run: func(cmd *cobra.Command, args []string) {

		interval, err := strconv.Atoi(INTERVAL)
		if err != nil {
			panic("Interval time must be a integer in seconds")
		}

		server := proc.NewServer(
			SLACK_URL,
			KEYCLOAK_URL,
			KEYCLOAK_USER,
			KEYCLOAK_PASSWORD,
			interval,
			KEYCLOAK_REALM,
		)
		logs.Logger.Info("Starting server...")
		server.Run()
	},
}

// Execute does the thing
func Execute() {
	logs.CreateLogger()

	if err := rootCmd.Execute(); err != nil {
		logs.Logger.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	viper.AutomaticEnv()
	rootCmd.PersistentFlags().StringVarP(&SLACK_URL, "slack-url", "s", get("SLACK_URL"), "slack url")
	rootCmd.PersistentFlags().StringVarP(&KEYCLOAK_URL, "keycloak-url", "k", get("KEYCLOAK_URL"), "keycloak root, no auth paths")
	rootCmd.PersistentFlags().StringVarP(&KEYCLOAK_USER, "keycloak-user", "u", get("KEYCLOAK_USER"), "keycloak admin user")
	rootCmd.PersistentFlags().StringVarP(&KEYCLOAK_PASSWORD, "keycloak-password", "p", get("KEYCLOAK_PASSWORD"), "keycloak admin password")
	rootCmd.PersistentFlags().StringVarP(&INTERVAL, "interval", "i", get("INTERVAL"), "interval to check in seconds")
	rootCmd.PersistentFlags().StringVarP(&KEYCLOAK_REALM, "keycloak-realm", "r", get("KEYCLOAK_REALM"), "keycloak realm")
}

func get(key string) string {
	if x := viper.Get(key); x != nil {
		return x.(string)
	}
	return ""
}
