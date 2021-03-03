package leaker

import "github.com/nbutton23/zxcvbn-go"

type ZxcvbnValidator struct {
	score float64
}

func NewZxcvbnValidator(score float64) ZxcvbnValidator {
	return ZxcvbnValidator{
		score: score,
	}
}

func (obj ZxcvbnValidator) IsValid(tok string) bool {
	return zxcvbn.PasswordStrength(tok, nil).Entropy >= obj.score
}
