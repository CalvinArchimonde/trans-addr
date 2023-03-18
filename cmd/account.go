package cmd

import (
	"encoding/json"
	"os"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/spf13/cobra"
)

const (
	flag_chain_filename = "flag.chain_info_filename"
)

type ChainInfo struct {
	ChainName    string `json:"chain_name"`
	Bech32Prefix string `json:"bech32_prefix"`
}

func showWalletAddrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [address]",
		Short: "Show different network address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fileName, err := cmd.Flags().GetString(flag_chain_filename)
			if err != nil {
				fileName = "chain.json"
			}

			chainInfos, err := readChainInfoFromFile(fileName)
			if err != nil {
				println(err.Error())
				return
			}

			if err = listAddressOnChains(args[0], chainInfos); err != nil {
				println(err.Error())
			}
		},
	}

	cmd.Flags().StringP(flag_chain_filename, "f", "chain.json", "set the filename of chain info")
	return cmd
}

func listAddressOnChains(address string, chainInfos []ChainInfo) error {
	_, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return err
	}

	for _, info := range chainInfos {
		addr, err := bech32.ConvertAndEncode(info.Bech32Prefix, bz)
		if err != nil {
			println(info.ChainName, " -- convert err.", err)
		} else {
			println(info.ChainName, " -- ", addr)
		}
	}
	return nil
}

func readChainInfoFromFile(fileName string) ([]ChainInfo, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	info := make([]ChainInfo, 0)
	if err = json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	return info, nil
}
