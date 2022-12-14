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
	"path/filepath"
)

// addWatchCmd represents the addWatch command
var addWatchCmd = &cobra.Command{
	Use:   "addWatch path index_ext",
	Short: "add the folder path or file path for monitoring upload",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		node, err := cmd.Flags().GetString(FlagNode)
		if err != nil {
			return err
		}
		fPath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		response, err := http.Get(node + fmt.Sprintf("/api/add_watch_file?path=%s&index=%s", fPath, args[1]))
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
	addWatchCmd.Flags().String(FlagNode, "http://127.0.0.1:8082", "node api")
	rootCmd.AddCommand(addWatchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addWatchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addWatchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
