package cmd

import "testing"

func TestReadChainInfoFromFile(t *testing.T) {
	_, err := readChainInfoFromFile("../chains.json")
	if err != nil {
		println(err.Error())
	}
}
