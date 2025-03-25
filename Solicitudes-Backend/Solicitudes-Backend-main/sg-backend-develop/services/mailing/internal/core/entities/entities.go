package entities

type User struct {
	ID             int64
	PersonID       *int64
	EmailValidated bool
}
