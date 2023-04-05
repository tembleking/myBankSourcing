package account

import "errors"

var (
	ErrAddMoneyQuantityCannotBeNegative = errors.New("the quantity cannot be negative")
	ErrBalanceIsNotEnoughForWithdrawal  = errors.New("balance is not enough for a withdrawal")
	ErrBalanceIsNotEnoughForTransfer    = errors.New("balance is not enough for a transfer")
	ErrAccountIsClosed                  = errors.New("account is closed")
	ErrAccountCannotBeClosedWithBalance = errors.New("account cannot be closed with balance")
)
