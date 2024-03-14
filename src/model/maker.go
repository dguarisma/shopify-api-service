package model

type Maker struct {
	CustomModel
	Name       string
	Status     bool
	Trademarks []Trademark
}
