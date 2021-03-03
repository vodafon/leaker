package leaker

type Validator interface {
	IsValid(string) bool
}
