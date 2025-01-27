package cmd

import (
	"context"

	"github.com/zeropsio/zcli/src/proto/daemon"

	"github.com/zeropsio/zcli/src/i18n"

	"github.com/zeropsio/zcli/src/cliAction/stopVpn"

	"github.com/spf13/cobra"
)

func vpnStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "stop",
		Short:        i18n.CmdVpnStop,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			regSignals(cancel)

			daemonClient, daemonCloseFunc, err := daemon.CreateClient(ctx)
			if err != nil {
				return err
			}
			defer daemonCloseFunc()

			return stopVpn.New(
				stopVpn.Config{},
				daemonClient,
			).Run(ctx, stopVpn.RunConfig{})
		},
	}
	cmd.Flags().BoolP("help", "h", false, helpText(i18n.VpnStopHelp))
	return cmd
}
