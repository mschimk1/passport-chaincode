package model

const UserObjectType = "User"

// User participant
type User struct {
	Entity
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetObjectType returns the blockchain object type
func (u *User) GetObjectType() string {
	return UserObjectType
}
