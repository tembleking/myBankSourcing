package account

import "errors"

var ErrAddMoneyQuantityCannotBeNegative = errors.New("the quantity cannot be negative")
var ErrBalanceIsNotEnoughForWithdrawal = errors.New("balance is not enough for a withdrawal")
