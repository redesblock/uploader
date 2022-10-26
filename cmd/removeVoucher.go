/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

// removeVoucherCmd represents the removeVoucher command
var removeVoucherCmd = &cobra.Command{
	Use:   "removeVoucher path",
	Short: "remove the voucher for uploading",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		node, err := cmd.Flags().GetString(FlagNode)
		if err != nil {
			return err
		}
		response, err := http.Get(node + fmt.Sprintf("/api/remove_voucher?voucher=%s", args[0]))
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			bts, _ := ioutil.ReadAll(response.Body)
			return fmt.Errorf("%s %s", response.Status, string(bts))
		}
		bts, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		fmt.Println("affect rows:", string(bts))
		return nil
	},
}

func init() {
	removeVoucherCmd.Flags().String(FlagNode, "http://127.0.0.1:8082", "node api")
	rootCmd.AddCommand(removeVoucherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeVoucherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeVoucherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
