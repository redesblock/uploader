/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

// addVoucherCmd represents the addVoucher command
var addVoucherCmd = &cobra.Command{
	Use:   "addVoucher voucherID Node",
	Args:  cobra.ExactArgs(2),
	Short: "add an usable voucher id from special node for file uploading",
	RunE: func(cmd *cobra.Command, args []string) error {
		node, err := cmd.Flags().GetString(FlagNode)
		if err != nil {
			return err
		}
		response, err := http.Get(node + fmt.Sprintf("/api/add_voucher?voucher=%s&node=%s", args[0], args[1]))
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
		var out bytes.Buffer
		json.Indent(&out, bts, "", "  ")
		fmt.Println(out.String())
		return nil
	},
}

func init() {
	addVoucherCmd.Flags().String(FlagNode, "http://127.0.0.1:8082", "node api")
	rootCmd.AddCommand(addVoucherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addVoucherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addVoucherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
