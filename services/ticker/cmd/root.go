package cmd

import (
	"fmt"
	"os"

	horizonclient "github.com/paydex-core/paydex-go/clients/horizonclient"
	hlog "github.com/paydex-core/paydex-go/support/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var DatabaseURL string
var Client *horizonclient.Client
var UseTestNet bool
var Logger = hlog.New()

var rootCmd = &cobra.Command{
	Use:   "ticker",
	Short: "Paydex Development Foundation Ticker.",
	Long:  `A tool to provide Paydex Asset and Market data.`,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(
		&DatabaseURL,
		"db-url",
		"d",
		"postgres://localhost:5432/paydexticker01?sslmode=disable",
		"database URL, such as: postgres://user:pass@localhost:5432/ticker",
	)
	rootCmd.PersistentFlags().BoolVar(
		&UseTestNet,
		"testnet",
		false,
		"use the Paydex Test Network, instead of the Paydex Public Network",
	)

	Logger.SetLevel(logrus.DebugLevel)
}

func initConfig() {
	if UseTestNet {
		Logger.Debug("Using Paydex Default Test Network")
		Client = horizonclient.DefaultTestNetClient
	} else {
		Logger.Debug("Using Paydex Default Public Network")
		Client = horizonclient.DefaultPublicNetClient
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
