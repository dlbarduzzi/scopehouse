package security

import (
	"regexp"
	"slices"
	"testing"
)

func TestRandomStringGenerator(t *testing.T) {
	testCases := []struct {
		name            string
		alphabets       []Alphabet
		expectedPattern string
	}{
		{
			name:            "lowercase",
			alphabets:       []Alphabet{AZLowercase},
			expectedPattern: `[a-z]+`,
		},
		{
			name:            "uppercase",
			alphabets:       []Alphabet{AZUppercase},
			expectedPattern: `[A-Z]+`,
		},
		{
			name:            "lowercase and numbers",
			alphabets:       []Alphabet{AZLowercase, Digits},
			expectedPattern: `[a-z0-9]+`,
		},
		{
			name:            "uppercase and dashes",
			alphabets:       []Alphabet{AZUppercase, Dashes},
			expectedPattern: `[A-Z_\-]+`,
		},
		{
			name:            "default alphabet",
			alphabets:       []Alphabet{},
			expectedPattern: `[a-zA-Z0-9_\-]+`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			length := 10
			attempt := 1

			generated := make([]string, 0, 500)

			for j := 0; j < 500; j++ {
				result := RandomStringGenerator(length, tc.alphabets)
				if len(result) != length {
					t.Fatalf("expected length to be %d, got %d", length, len(result))
				}

				reg := regexp.MustCompile(tc.expectedPattern)
				if match := reg.MatchString(result); !match {
					t.Fatalf(
						"expected result to have only characters %s, got %s",
						tc.expectedPattern, result,
					)
				}

				if slices.Contains(generated, result) {
					if attempt > 3 {
						t.Fatalf(
							"expected not to repeat random string - test (%d) - found %q in %q",
							j, result, generated,
						)
					}
					// rerun
					continue
				}

				generated = append(generated, result)
			}
		})
	}
}
