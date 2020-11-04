package main

import (
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"istio.io/istio/pilot/pkg/request"
)

var (
	requestCmd = &cobra.Command{
		Use:   "request <method> <path> [<body>]",
		Short: "Makes an HTTP request to Pilot metrics/debug endpoint",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			command := &request.Command{
				Address: "127.0.0.1:15014",
				Client: &http.Client{
					Timeout: 60 * time.Second,
				},
			}
			body := ""
			if len(args) > 2 {
				body = args[2]
			}
			return command.Do(args[0], args[1], body)
		},
	}
)

func init() {
	rootCmd.AddCommand(requestCmd)
}
