package apis

import "github.com/dlbarduzzi/scopehouse/tools/security"

func generateId(size int) string {
	return security.RandomStringGenerator(size, []security.Alphabet{
		security.Digits,
		security.AZLowercase,
		security.AZUppercase,
	})
}
