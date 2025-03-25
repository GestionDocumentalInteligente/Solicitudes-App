package entities

type Person struct {
	ID                   int64  // Unique identifier for the person
	Cuil                 string // Person's CUIL
	Dni                  string // Person's DNI (optional)
	FirstName            string // Person's first name
	LastName             string // Person's last name
	Email                string // Person's email
	Phone                string // Person's phone number
	IsVerified           bool
	AcceptsNotifications bool
}
