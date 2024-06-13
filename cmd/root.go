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
		listComosChainAddressCmd(),
		listDymHubAddressCmd(),
	)

	cmd.CompletionOptions.DisableDefaultCmd = true
	return cmd
}

func createWalletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new wallet",
		Run: func(cmd *cobra.Command, args []string) {
			if err := createWallet(); err != nil {
				println(err.Error())
			}
		},
	}
	return cmd
}

func createWallet() error {
	if mnemonic, err := NewMnemonic(); err != nil {
		return err
	} else {
		// create cosmos wallet
		if err = CreateCmsWallet(mnemonic); err != nil {
			return err
		}

		// create dymsion wallet
		println()
		CreateDymsionWallet(mnemonic)

		// create ethereum wallet
		println()
		CreateEtherWallet()

		println()
		println("**Important** write this mnemonic phrase in a safe place.")
		println("It is the only way to recover your account if you ever forget your password.")
		println()
		println(mnemonic)
		println()
	}

	return nil
}
