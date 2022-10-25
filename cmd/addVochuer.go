/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/redesblock/uploader/core/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addVoucherCmd represents the addVoucher command
var addVoucherCmd = &cobra.Command{
	Use:   "addVoucher voucherID Host",
	Args:  cobra.ExactArgs(2),
	Short: "add an usable voucher",
	RunE: func(cmd *cobra.Command, args []string) error {
		voucher := args[0]
		host := args[1]

		db, err := model.New(viper.GetString(FlagDBMode), viper.GetString(FlagDBDSN))
		if err != nil {
			return err
		}
		var item model.Voucher
		if res := db.Where("voucher = ?", voucher).Find(&item); res.Error != nil {
			return res.Error
		} else if res.RowsAffected == 0 {
			item.Voucher = voucher
			item.Usable = true
		}
		item.Host = host

		if err := db.Save(&item).Error; err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addVoucherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addVoucherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addVoucherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
