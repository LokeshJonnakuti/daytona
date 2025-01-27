// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"fmt"
	"time"

	"github.com/daytonaio/daytona/pkg/server"
	"github.com/daytonaio/daytona/pkg/server/config"
	"github.com/daytonaio/daytona/pkg/server/frpc"
	"github.com/daytonaio/daytona/pkg/views/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server process",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if log.GetLevel() < log.InfoLevel {
			//	for now, force the log level to info when running the server
			log.SetLevel(log.InfoLevel)
		}
		errCh := make(chan error)

		err := server.Start(errCh)
		if err != nil {
			log.Fatal(err)
		}

		c, err := config.GetConfig()
		if err != nil {
			log.Fatal(err)
		}

		select {
		case err := <-errCh:
			log.Fatal(err)
		// TODO: This is an optimistic check. We should check if the server is actually running
		case <-time.After(5 * time.Second):
			util.RenderBorderedMessage(fmt.Sprintf("Daytona Server running on port: %d.\nTo connect to the server remotely, use the following command on the client machine:\n\ndaytona profile add -a %s", c.ApiPort, frpc.GetApiUrl(c)))
		}

		err = <-errCh
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	ServerCmd.AddCommand(configureCmd)
	ServerCmd.AddCommand(configCmd)
	ServerCmd.AddCommand(logsCmd)
}
