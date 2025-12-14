package cmd

import (
	"fmt"

	"plan/internal/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the webapp server",
	Long:  `Start the HTTP server that serves the webapp and provides API endpoints.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		host, _ := cmd.Flags().GetString("host")
		dev, _ := cmd.Flags().GetBool("dev")
		webDir, _ := cmd.Flags().GetString("web-dir")

		addr := fmt.Sprintf("%s:%d", host, port)
		fmt.Printf("Starting server at http://%s\n", addr)
		fmt.Printf("SSE endpoint: http://%s/api/events\n", addr)

		if dev {
			fmt.Println("Development mode: proxying to Vite dev server at http://localhost:5173")
		}

		return server.Start(addr, dev, webDir)
	},
}

func init() {
	serveCmd.Flags().IntP("port", "p", 8080, "Server port")
	serveCmd.Flags().String("host", "localhost", "Server host")
	serveCmd.Flags().Bool("dev", false, "Development mode (proxy to Vite dev server)")
	serveCmd.Flags().String("web-dir", "../web/dist", "Path to web dist directory")
}
