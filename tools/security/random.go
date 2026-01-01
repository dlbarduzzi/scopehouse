package security

import "math/rand"

type Alphabet string

//nolint:lll
const defaultAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"

const (
	Dashes          Alphabet = "_-"
	Digits          Alphabet = "0-9"
	AZLowercase     Alphabet = "a-z"
	AZUppercase     Alphabet = "A-Z"
	DefaultAlphabet Alphabet = defaultAlphabet
)

func getAlphabet(alphabet Alphabet) string {
	switch alphabet {
	case Dashes:
		return "-_"
	case Digits:
		return "0123456789"
	case AZLowercase:
		return "abcdefghijklmnopqrstuvwxyz"
	case AZUppercase:
		return "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	return defaultAlphabet
}

func RandomStringGenerator(size int, alphabets []Alphabet) string {
	alphabet := ""

	if len(alphabets) < 1 {
		alphabets = append(alphabets, DefaultAlphabet)
	}

	for _, a := range alphabets {
		alphabet += getAlphabet(a)
	}

	b := make([]byte, size)
	max := len(alphabet)

	for i := range b {
		b[i] = alphabet[rand.Intn(max)]
	}

	return string(b)
}
