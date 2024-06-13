package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"
)

const (
	flag_chain_filename = "flag.chain_info_filename"
)

type ChainInfo struct {
	ChainName    string `json:"chain_name"`
	Bech32Prefix string `json:"bech32_prefix"`
}

func listAddressCmd(use string, f func(string, map[string][]ChainInfo) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: "Show different network address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fileName, err := cmd.Flags().GetString(flag_chain_filename)
			if err != nil {
				fileName = "chains.json"
			}

			chains, err := readChainInfoFromFile(fileName)
			if err != nil {
				println(err.Error())
				return
			}

			if err = f(args[0], chains); err != nil {
				println(err.Error())
			}
		},
	}

	cmd.Flags().StringP(flag_chain_filename, "f", "chains.json", "set the filename of chain info")
	return cmd
}

func listComosChainAddressCmd() *cobra.Command {
	return listAddressCmd("cosmos [address]", listComosChainAddress)
}

func listComosChainAddress(address string, chains map[string][]ChainInfo) error {
	_, data, err := bech32.Decode(address)
	if err != nil {
		return fmt.Errorf("decoding bech32 failed: %w", err)
	}

	converted, _ := bech32.ConvertBits(data, 5, 8, false)
	println("- PublicKey:", base64.StdEncoding.EncodeToString(converted))
	println("- PublicKey:", strings.ToUpper(hex.EncodeToString(converted)), "bytes len=", len(converted))
	println("- Address List:")
	chainInfos := chains["CosmosHub"]
	for _, info := range chainInfos {
		addr, err := bech32.Encode(info.Bech32Prefix, data)
		if err != nil {
			println("  -", info.ChainName, "convert err.", err)
		} else {
			println("  -", info.ChainName, ":", addr)
		}
	}

	hash := sha3.NewLegacyKeccak256()
	hash.Write(converted)
	hashedPubKey := hash.Sum(nil)
	// Take the last 20 bytes to form the Ethereum address
	ethAddress := "0x" + hex.EncodeToString(hashedPubKey[len(hashedPubKey)-20:])
	println("ethereum address:", ethAddress)

	return nil
}

func listDymHubAddressCmd() *cobra.Command {
	return listAddressCmd("dym [address]", listDymChainAddress)
}

func listDymChainAddress(address string, chains map[string][]ChainInfo) error {
	var bz []byte
	if strings.HasPrefix(address, "0x") {
		if tmp, err := hexutil.Decode(address); err != nil {
			return err
		} else {
			bz, err = bech32.ConvertBits(tmp, 8, 5, true)
			if err != nil {
				return err
			}
		}
	} else {
		var err error
		if _, bz, err = bech32.Decode(address); err != nil {
			return err
		}
	}

	println("- Address List:")
	converted, _ := bech32.ConvertBits(bz, 5, 8, false)
	println("  - evm address:", hexutil.Encode(converted))

	chainInfos := chains["Dymsion"]
	for _, info := range chainInfos {
		addr, err := bech32.Encode(info.Bech32Prefix, bz)
		if err != nil {
			println("  -", info.ChainName, "convert err.", err)
		} else {
			println("  -", info.ChainName, ":", addr)
		}
	}

	return nil
}

func readChainInfoFromFile(fileName string) (map[string][]ChainInfo, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	chains := make(map[string][]ChainInfo)
	if err = json.Unmarshal(data, &chains); err != nil {
		return nil, err
	}

	return chains, nil
}
