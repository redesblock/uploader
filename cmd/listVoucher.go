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
	"strconv"
)

// listVoucherCmd represents the listVoucher command
var listVoucherCmd = &cobra.Command{
	Use:   "listVoucher [page_size] [page_num]",
	Short: "list vouchers for uploading file",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		page_num := int64(1)
		page_size := int64(10)
		if len(args) > 1 {
			n, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid page_size")
			}
			page_size = n
		}
		if len(args) > 2 {
			n, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid page_num")
			}
			page_num = n
		}
		node, err := cmd.Flags().GetString(FlagNode)
		if err != nil {
			return err
		}
		response, err := http.Get(node + fmt.Sprintf("/api/vouchers?page_num=%d&page_size=%d", page_num, page_size))
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
	listVoucherCmd.Flags().String(FlagNode, "http://127.0.0.1:8082", "node api")
	rootCmd.AddCommand(listVoucherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listVoucherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listVoucherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
