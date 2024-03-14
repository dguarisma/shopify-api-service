package utils

type GeneralCase struct {
	Title          string
	Method         string
	ExpectedStatus int
	BodyRq         interface{}
	ExpectedBodyRs interface{}
}
