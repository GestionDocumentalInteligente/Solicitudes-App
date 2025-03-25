package user

import (
	"errors"
	"fmt"
)

type UserType string

const (
	Admin    UserType = "admin"
	Owner    UserType = "owner"
	Occupant UserType = "occupant"
)

type User struct {
	Type           UserType
	Cuil           int64
	DocumentNumber int64
	FirstName      string
	LastName       string
	Email          string
	Phone          int64
	Address        Address
}

func GetFullName(user User) string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

type Address struct {
	ABLNumber int64
	ZipCode   int64
	Street    string
	Number    string
}

func (u UserType) IsValid() error {
	switch u {
	case Admin, Owner, Occupant:
		return nil
	default:
		return errors.New("invalid user_type: must be 'admin', 'owner', or 'occupant'")
	}
}
