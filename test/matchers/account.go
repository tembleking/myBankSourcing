package matchers

import (
	"fmt"

	"github.com/onsi/gomega/gcustom"

	"github.com/tembleking/myBankSourcing/pkg/domain"
)

func BeAnEntityEqualTo(expected domain.Entity) gcustom.CustomGomegaMatcher {
	return gcustom.MakeMatcher(func(actual interface{}) (success bool, err error) {
		if expected == nil {
			return actual == nil, nil
		}

		if actual == nil {
			return false, nil
		}

		if actualEntity, ok := actual.(domain.Entity); ok {
			return expected.SameEntityAs(actualEntity), nil
		}

		return false, nil
	}).WithMessage(fmt.Sprintf("expected entity to be equal to %#v", expected))
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
