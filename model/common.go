package model

// Model for a blockchain table
type Model interface {
	GetObjectType() string
	Validate() error
}

// Entity holds the base type information
type Entity struct {
	ObjectType string `json:"docType"` // docType is used to distinguish the various types of objects
}
