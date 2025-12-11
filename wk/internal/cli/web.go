package cli

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func newWebCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "Run a simple web server with a hello endpoint",
		Long: `Start a minimal HTTP server useful for quick checks or integrations.

It exposes:
  - GET /hello  -> returns "hello, world"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}

			mux := http.NewServeMux()
			mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "hello, world")
			})

			addr := ":" + port
			log.Printf("Starting wk web server on %s (GET /hello)", addr)
			if err := http.ListenAndServe(addr, mux); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}
