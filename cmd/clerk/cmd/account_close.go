/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tembleking/myBankSourcing/internal/factory"
)

// closeCmd represents the close command
var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "Closes an account",
	Run: func(cmd *cobra.Command, args []string) {
		accountClosed, err := factory.NewFactory().NewAccountService().CloseAccount(cmd.Context(), args[0])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		cmd.Printf("Closed Account ID: %s\n", accountClosed.ID())
	},
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			accounts := factory.NewFactory().NewAccountProjection(cmd.Context()).Accounts()

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
	accountCmd.AddCommand(closeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// closeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// closeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
