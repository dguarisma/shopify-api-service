package model

type CategoryOne struct {
	CustomModel
	Name   string
	Status bool
}

type CategoryTwo struct {
	CustomModel
	Name          string
	Status        bool
	CategoryOneID uint
}

type CategoryThree struct {
	CustomModel
	Name          string
	Status        bool
	CategoryOneID uint
	CategoryTwoID uint
}
