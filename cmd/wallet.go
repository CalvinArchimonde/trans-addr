package cmd

import (
	"github.com/spf13/cobra"
)

func createWalletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new wallet",
		Run: func(cmd *cobra.Command, args []string) {
			println("Coding....")
		},
	}
	return cmd
}
