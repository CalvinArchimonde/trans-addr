package main

import "github.com/calvinarchimonde/trans-addr/cmd"

func main() {

	rootCmd := cmd.NewRootCmd()
	rootCmd.Execute()
}
