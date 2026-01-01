package apis

import (
	"regexp"
	"slices"
	"testing"
)

func TestGenerateId(t *testing.T) {
	size := 10
	pattern := `[a-zA-Z0-9]+`

	attempt := 1
	generated := make([]string, 0, 500)

	for j := 0; j < 500; j++ {
		result := generateId(10)
		if len(result) != size {
			t.Fatalf("expected size to be %d, got %d", size, len(result))
		}

		reg := regexp.MustCompile(pattern)
		if match := reg.MatchString(result); !match {
			t.Fatalf(
				"expected result to have only characters %s, got %s",
				pattern, result,
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
}
