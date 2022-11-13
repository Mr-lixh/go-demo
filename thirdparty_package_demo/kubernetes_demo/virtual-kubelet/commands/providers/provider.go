package providers

import (
	"fmt"
	"github.com/Mr-lixh/go-demo/thirdparty_package_demo/kubernetes_demo/virtual-kubelet/provider"
	"github.com/spf13/cobra"
	"os"
)

// NewCommand creates a new providers subcommand
// This subcommand is used to determine which providers are registered.
func NewCommand(s *provider.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "providers",
		Short: "Show the list of supported providers",
		Long:  "Show the list of supported providers",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				for _, p := range s.List() {
					fmt.Fprintln(cmd.OutOrStdout(), p)
				}
			case 1:
				if !s.Exists(args[0]) {
					fmt.Fprintln(cmd.OutOrStderr(), "no such providers", args[0])

					// TODO(@cpuuy83): would be nice to not short-circuit the exit here
					// But at the moment this seems to be the only way to exit non-zero and
					// handle our own error output
					os.Exit(1)
				}
				fmt.Fprintln(cmd.OutOrStdout(), args[0])
			}
		},
	}
}
