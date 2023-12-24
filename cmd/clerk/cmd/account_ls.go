/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tembleking/myBankSourcing/internal/factory"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists the accounts created",
	Run: func(cmd *cobra.Command, args []string) {
		accounts := factory.NewFactory().NewAccountProjection(cmd.Context()).Accounts()

		for _, account := range accounts {
			account := account
			cmd.Printf("Account ID: %s, Balance: %d\n", account.AccountID, account.Balance)
			if len(account.Movements) != 0 {
				cmd.Println("Movements:")
			}
			for _, movement := range account.Movements {
				movement := movement
				cmd.Printf("\t%s of %d, resulting in %d\n", movement.Type, movement.Amount, movement.ResultingBalance)
			}
		}
	},
}

func init() {
	accountCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
