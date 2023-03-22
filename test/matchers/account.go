package matchers

import (
	"fmt"

	"github.com/onsi/gomega/gcustom"

	"github.com/tembleking/myBankSourcing/pkg/domain/account"
)

func BeAnAccountEqualsTo(expected *account.Account) gcustom.CustomGomegaMatcher {
	return gcustom.MakeMatcher(func(actual interface{}) (success bool, err error) {
		if expected == nil {
			return actual == nil, nil
		}

		if actual == nil {
			return false, nil
		}

		actualAccount, ok := actual.(*account.Account)
		if !ok {
			return false, fmt.Errorf("expected an account, got %T", actual)
		}

		return actualAccount.ID() == expected.ID() &&
			actualAccount.Balance() == expected.Balance() &&
			actualAccount.IsOpen() == expected.IsOpen(), nil
	}).WithMessage(fmt.Sprintf("expected account to be equal to %#v", expected))
}
