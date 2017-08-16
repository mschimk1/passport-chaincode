package model

// UserObjectType blockchain object type
const UserObjectType = "User"

// User participant
type User struct {
	Entity
	ID   string `json:"id"`
	Name string `json:"name"`
}
