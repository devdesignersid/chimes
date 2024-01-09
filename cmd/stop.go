package cmd

import (
	"time"

	"github.com/devdesignersid/chimes/pkg/daemon"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		d := daemon.NewDaemon("chimes.pid", "chimes.log", 1*time.Second)
		d.Kill()

	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
