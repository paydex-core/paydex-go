package main

import (
	"log"

	"github.com/paydex-core/paydex-go/tools/paydex-hd-wallet/commands"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:   "paydex-hd-wallet",
	Short: "Simple HD wallet for Paydex Lumens. THIS PROGRAM IS STILL EXPERIMENTAL. USE AT YOUR OWN RISK.",
}

func init() {
	mainCmd.AddCommand(commands.NewCmd)
	mainCmd.AddCommand(commands.AccountsCmd)
}

func main() {
	if err := mainCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
