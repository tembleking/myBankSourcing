/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/tembleking/myBankSourcing/internal/factory"
	"github.com/tembleking/myBankSourcing/pkg/domain/account"
)

// addMoneyCmd represents the addMoney command
var addMoneyCmd = &cobra.Command{
	Use:   "addMoney",
	Short: "Adds money to an account",
	Run: func(cmd *cobra.Command, args []string) {
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			panic(fmt.Errorf("invalid amount %s: %w", args[1], err))
		}
		account, err := factory.NewFactory().NewAccountService().AddMoneyToAccount(cmd.Context(), account.ID(args[0]), amount)
		if err != nil {
			panic(err)
		}

		cmd.Printf("Account ID: %s, Added: %d, Balance: %d\n", account.ID(), amount, account.Balance())
	},
	Args: cobra.MinimumNArgs(2),
}

func init() {
	accountCmd.AddCommand(addMoneyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addMoneyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addMoneyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
