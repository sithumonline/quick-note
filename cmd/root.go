package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sithumonline/quick-note/cmd/server"
)

var rootCmd = &cobra.Command{
	Use:   "quick-note",
	Short: "template repo",
	Long:  `Golang template repo form Golang Sri Lanka `,
}

func init() {
	rootCmd.AddCommand(server.RunServerCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
