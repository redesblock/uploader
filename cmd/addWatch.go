/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/redesblock/uploader/core/model"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// addWatchCmd represents the addWatch command
var addWatchCmd = &cobra.Command{
	Use:   "addWatch path index_ext",
	Short: "add a watch file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("invalid path (%s) %v", path, err)
		}
		fullPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("invalid path (%s) %v", path, err)
		}

		index := args[1]
		if len(index) == 0 || strings.ContainsRune(index, os.PathSeparator) {
			return fmt.Errorf("invalid index ext (%s)", index)
		}

		db, err := model.New(viper.GetString(FlagDBMode), viper.GetString(FlagDBDSN))
		if err != nil {
			return err
		}

		var item model.WatchFile
		if res := db.Where("path = ?", fullPath).Find(&item); res.Error != nil {
			return res.Error
		} else if res.RowsAffected == 0 {
			item.Path = fullPath
		}
		item.IndexExt = index

		if err := db.Save(&item).Error; err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addWatchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addWatchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addWatchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
