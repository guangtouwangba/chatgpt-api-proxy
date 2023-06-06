/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"chatgpt-api-proxy/pkg/server"

	"github.com/spf13/cobra"
)

const defaultPort = 8080

// serverCmd represents the server command.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launches the server",
	Long: `Launches the server. For example:
	proxy server -p 8080
 `,

	Run: runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().IntP("port", "p", defaultPort, "Port to listen on")
	// path to config file
	serverCmd.Flags().StringP("config", "c", "", "Path to config file")
}

func runServer(_ *cobra.Command, _ []string) {
	server.RunServer()
}
