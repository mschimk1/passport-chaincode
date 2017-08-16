package model

// Model for a blockchain table
type Model interface {
	GetObjectType() string
}

// Entity holds the base type information
type Entity struct {
	ObjectType string `json:"docType"` // docType is used to distinguish the various types of objects
}

// GetObjectType returns the blockchain object type
func (e *Entity) GetObjectType() string {
	return e.ObjectType
}
