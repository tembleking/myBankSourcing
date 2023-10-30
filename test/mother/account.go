package mother

import "github.com/tembleking/myBankSourcing/pkg/account"

// AccountOpenWithMovements returns an open Account that:
// - is open
// - has 3 movements:
//   - deposit 50
//   - withdraw 30
//   - withdraw 15
//
// - has 5 euros remaining
// - has version 4
func AccountOpenWithMovements() *account.Account {
	acc, _ := account.OpenAccount("some-account")
	_ = acc.DepositMoney(50)
	_ = acc.WithdrawMoney(30)
	_ = acc.WithdrawMoney(15)
	// remaining money: 5
	// version: 4
	return acc
}
