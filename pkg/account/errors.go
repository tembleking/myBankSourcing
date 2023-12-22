package account

import "errors"

var (
	ErrQuantityCannotBeNegative         = errors.New("the quantity cannot be negative")
	ErrBalanceIsNotEnough               = errors.New("balance is not enough")
	ErrAccountIsClosed                  = errors.New("account is closed")
	ErrAccountCannotBeClosedWithBalance = errors.New("account cannot be closed with balance")
	ErrCannotTransferToSameAccount      = errors.New("cannot transfer to the same account")
)
