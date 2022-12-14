/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/redesblock/uploader/core/model"
	"github.com/redesblock/uploader/core/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagInterval = "interval"
	FlagPort     = "port"
	FlagDBMode   = "database_mode"
	FlagDBDSN    = "database_dsn"
	FlagLevel    = "log_level"
	FlagGateWay  = "gateway"
	FlagNode     = "node"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "automatically upload files or folders to mop cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		level, err := log.ParseLevel(viper.GetString(FlagLevel))
		if err != nil {
			return err
		}
		log.SetLevel(level)

		db, err := model.New(viper.GetString(FlagDBMode), viper.GetString(FlagDBDSN))
		if err != nil {
			return err
		}
		return server.Start(":"+viper.GetString(FlagPort), db, viper.GetString(FlagInterval), viper.GetString(FlagGateWay))
	},
}

func init() {
	serveCmd.Flags().String(FlagPort, "8082", "listen port")
	serveCmd.Flags().String(FlagInterval, "10m", "watcher poll interval")
	serveCmd.Flags().String(FlagGateWay, "https://gateway.mopweb3.cn", "mop gateway")

	viper.BindPFlag(FlagPort, serveCmd.Flags().Lookup(FlagPort))
	viper.BindPFlag(FlagInterval, serveCmd.Flags().Lookup(FlagInterval))
	viper.BindPFlag(FlagGateWay, serveCmd.Flags().Lookup(FlagGateWay))

	serveCmd.PersistentFlags().String(FlagDBMode, "sqlite", "database mode, sqlite、mysql、postgres")
	serveCmd.PersistentFlags().String(FlagDBDSN, "sqlite.db", "database source name")
	serveCmd.PersistentFlags().String(FlagLevel, "info", "log level")

	viper.BindPFlag(FlagDBMode, serveCmd.PersistentFlags().Lookup(FlagDBMode))
	viper.BindPFlag(FlagDBDSN, serveCmd.PersistentFlags().Lookup(FlagDBDSN))
	viper.BindPFlag(FlagLevel, serveCmd.PersistentFlags().Lookup(FlagLevel))

	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
