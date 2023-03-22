/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tembleking/myBankSourcing/internal/factory"
)

// accountOpenCmd represents the open command
var accountOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open an account",
	Run: func(cmd *cobra.Command, args []string) {
		service := factory.NewAccountService()

		account, err := service.OpenAccount(cmd.Context())
		if err != nil {
			panic(err)
		}

		fmt.Println("Account created:", account.ID())

		accounts, err := service.ListAccounts(cmd.Context())
		if err != nil {
			panic(err)
		}

		fmt.Println("Your Accounts:")
		for _, account := range accounts {
			fmt.Println("-", account.ID())
		}
	},
}

func init() {
	accountCmd.AddCommand(accountOpenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// accountOpenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// accountOpenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
