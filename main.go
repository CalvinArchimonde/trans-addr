package main

import "github.com/calvinarchimonde/trans-addr/cmd"

func main() {

	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		println(err.Error())
	}
}
