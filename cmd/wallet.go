package cmd

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
)

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
	kb := keyring.NewInMemory(&codec.ProtoCodec{})
	mnemonic, err := NewMnemonic()
	if err != nil {
		return err
	}

	k, err := kb.NewAccount("default", mnemonic, keyring.DefaultBIP39Passphrase, sdk.FullFundraiserPath, hd.Secp256k1)
	if err != nil {
		return err
	}
	address, err := bech32.ConvertAndEncode("cosmos", k.PubKey.Value)
	if err != nil {
		return err
	}
	println("- ", address)
	println("\n**Important** write this mnemonic phrase in a safe place.")
	println("It is the only way to recover your account if you ever forget your password.\n")
	println(mnemonic)
	return nil
}

func NewMnemonic() (string, error) {
	// 12 words mnemonic
	entropySeed, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
