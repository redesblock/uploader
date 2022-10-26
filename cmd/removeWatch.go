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

// removeWatchCmd represents the removeWatch command
var removeWatchCmd = &cobra.Command{
	Use:   "removeWatch path",
	Short: "remove the folder path or file path for monitoring upload",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		node, err := cmd.Flags().GetString(FlagNode)
		if err != nil {
			return err
		}
		response, err := http.Get(node + fmt.Sprintf("/api/remove_watch_file?path=%s", args[0]))
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
	removeWatchCmd.Flags().String(FlagNode, "http://127.0.0.1:8082", "node api")
	rootCmd.AddCommand(removeWatchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeWatchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeWatchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
