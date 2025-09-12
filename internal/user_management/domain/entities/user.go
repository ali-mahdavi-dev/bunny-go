package entities

import (
	"bunny-go/internal/framwork/adapter"
	cerrors "bunny-go/internal/framwork/cerrors"
	commandeventhandler "bunny-go/internal/framwork/service_layer/command_event_handler"
)

type User struct {
	adapter.BaseEntity
	Age      int
	UserName string
	Amount   int
	Trades   []Trade
	Events   []commandeventhandler.EventHandler `gorm:"-"`
}

func NewUser(userName string, age int, amount int) (*User, error) {
	if userName == "admin" {
		return nil, cerrors.BadRequest("Transaction.Invalid")
	}
	if age < 18 {
		return nil, cerrors.BadRequest("Transaction.AgeInvalid")
	}
	user := &User{}
	user.UserName = userName
	user.Age = age
	user.Amount = amount
	return user, nil
}

func (u *User) Update(userName string, age int, amount int) error {
	if userName == "admin" {
		return cerrors.BadRequest("Transaction.Invalid")
	}
	if age < 18 {
		return cerrors.BadRequest("Transaction.AgeInvalid")
	}

	u.UserName = userName
	u.Age = age
	u.Amount = amount
	return nil
}

func (u *User) Event() []commandeventhandler.EventHandler {
	return u.Events
}
