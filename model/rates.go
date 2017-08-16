package model

const RatesObjectType = "Rates"

// Rates represents exchange rates for a given base at a given date
type Rates struct {
	Entity
	Base       string `json:"base"`
	Date       string `json:"date"`
	Currencies struct {
		HKD float64
		NZD float64
		SGD float64
	} `json:"rates"`
}

// GetObjectType returns the blockchain table name
func (r *Rates) GetObjectType() string {
	return r.ObjectType
}
