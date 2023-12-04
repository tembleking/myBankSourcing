package matchers

import (
	"fmt"

	"github.com/tembleking/myBankSourcing/pkg/account"

	"github.com/onsi/gomega/gcustom"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

func BeAnAccountEqualsTo(expected *account.Account) gcustom.CustomGomegaMatcher {
	return gcustom.MakeMatcher(func(actual interface{}) (success bool, err error) {
		if expected == nil {
			return actual == nil, nil
		}

		if actual == nil {
			return false, nil
		}

		if actualAccount, ok := actual.(*account.Account); ok {
			return expected.SameEntityAs(actualAccount), nil
		}

		return false, nil
	}).WithMessage(fmt.Sprintf("expected account to be equal to %#v", expected))
}

func BeAggregateWithTheSameVersionAs(expected domain.Aggregate) gcustom.CustomGomegaMatcher {
	return gcustom.MakeMatcher(func(actual interface{}) (success bool, err error) {
		actualAggregate, ok := actual.(domain.Aggregate)
		if !ok {
			return false, fmt.Errorf("expected an aggregate, got %T", actual)
		}

		return actualAggregate.Version() == expected.Version(), nil
	}).WithMessage(fmt.Sprintf("be an aggregate with version %d", expected.Version()))
}
