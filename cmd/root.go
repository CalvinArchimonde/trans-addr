package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trans-addr",
		Short: "Create a new wallet. Translate address between chains.",
	}

	cmd.AddCommand(
		createWalletCmd(),
		showWalletAddrCmd(),
	)

	cmd.CompletionOptions.DisableDefaultCmd = true
	return cmd
}
