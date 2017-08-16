package model

// RatesObjectType blockchain object type
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
