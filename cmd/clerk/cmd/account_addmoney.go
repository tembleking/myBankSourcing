/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/tembleking/myBankSourcing/internal/factory"
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
		account, err := factory.NewFactory().NewAccountService().AddMoneyToAccount(cmd.Context(), args[0], amount)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		cmd.Printf("Account ID: %s, Added: %d, Balance: %d\n", account.ID(), amount, account.Balance())
	},
	Args: cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			accounts := factory.NewFactory().NewAccountView().Accounts()

			ids := make([]string, len(accounts))
			for i, account := range accounts {
				ids[i] = string(account.ID())
			}

			return ids, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
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
