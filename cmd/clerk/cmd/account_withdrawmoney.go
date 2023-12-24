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

// withdrawMoneyCmd represents the withdrawMoney command
var withdrawMoneyCmd = &cobra.Command{
	Use:   "withdrawMoney",
	Short: "Withdraws money from an account",
	Run: func(cmd *cobra.Command, args []string) {
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			cmd.PrintErrln(fmt.Errorf("invalid amount %s: %w", args[1], err))
		}
		updatedAccount, err := factory.NewFactory().NewAccountService().WithdrawMoneyFromAccount(cmd.Context(), args[0], amount)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		cmd.Printf("Account ID: %s, Withdrawn: %d, Balance: %d\n", updatedAccount.ID(), amount, updatedAccount.Balance())
	},
	Args: cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			accounts := factory.NewFactory().NewAccountProjection(cmd.Context()).Accounts()

			ids := make([]string, len(accounts))
			for i, account := range accounts {
				ids[i] = account.AccountID
			}

			return ids, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	accountCmd.AddCommand(withdrawMoneyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// withdrawMoneyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// withdrawMoneyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
