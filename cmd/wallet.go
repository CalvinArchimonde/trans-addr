package cmd

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"golang.org/x/crypto/sha3"
)

func CreateEtherWallet() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 0x96216849c49358b10257cb55b28ea603c874b05e
}

func CreateDymsionWallet(mnemonic string) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		log.Fatal(err)
	}

	publicKey, _ := wallet.PublicKeyHex(account)
	if err != nil {
		log.Fatal(err)
	}

	bz, err := bech32.ConvertBits(account.Address.Bytes(), 8, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	dmyAddr, err := bech32.Encode("dym", bz)
	if err != nil {
		log.Fatal(err)
	}

	println("- dymsion:")
	println("  - Private key in hex:", privateKey)
	println("  - Public key in hex:", publicKey)
	println("  - EVM address:", account.Address.Hex())
	println("  - Dym address:", dmyAddr)
}

func CreateCmsWallet(mnemonic string) error {
	kb := keyring.NewInMemory(&codec.ProtoCodec{})
	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return err
	}
	// hd.Secp256k1
	hdPath := hd.CreateHDPath(sdk.CoinType, 0, 0).String()
	k, err := kb.NewAccount("default", mnemonic, keyring.DefaultBIP39Passphrase, hdPath, algo)
	if err != nil {
		return err
	}

	address, err := k.GetAddress()
	if err != nil {
		return err
	}
	pk, _ := k.GetPubKey()
	println("- cosmos hub:")
	println("  - publicKey:", pk.String())
	println("  - address", address.String())
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
