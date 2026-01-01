package event

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewApiError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		apiErr  *ApiError
		content string
		message string
	}{
		{
			name:    "status 300",
			apiErr:  NewApiError(300, "hello world"),
			content: `{"status":300,"message":"Hello world."}`,
			message: "Hello world.",
		},
		{
			name:    "status 400",
			apiErr:  NewApiError(400, ""),
			content: `{"status":400,"message":"Bad Request."}`,
			message: "Bad Request.",
		},
		{
			name:    "status 500",
			apiErr:  NewApiError(500, "Hello world!"),
			content: `{"status":500,"message":"Hello world!"}`,
			message: "Hello world!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := tc.apiErr

			res, err := json.Marshal(e)
			if err != nil {
				t.Fatal(err)
			}

			resStr := string(res)

			if resStr != tc.content {
				t.Fatalf("expected content to be \n%v \ngot \n%v", tc.content, resStr)
			}

			if e.Error() != tc.message {
				t.Fatalf("expected error message to be %q, got %q", tc.message, e.Error())
			}
		})
	}
}

func TestNewInternalServerError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		apiErr  *ApiError
		content []string
		message string
	}{
		{
			name:   "empty message",
			apiErr: NewInternalServerError(""),
			content: []string{
				`"status":500`,
				`Something went wrong while processing this request.`,
			},
			message: "Something went wrong while processing this request.",
		},
		{
			name:   "custom message",
			apiErr: NewInternalServerError("hello world"),
			content: []string{
				`"status":500`,
				`Hello world.`,
			},
			message: "Hello world.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := tc.apiErr

			res, err := json.Marshal(e)
			if err != nil {
				t.Fatal(err)
			}

			resStr := string(res)

			for _, content := range tc.content {
				if !strings.Contains(resStr, content) {
					t.Errorf("expected content %v in response body \n%v", content, resStr)
				}
			}

			if e.Error() != tc.message {
				t.Fatalf("expected error message to be %q, got %q", tc.message, e.Error())
			}
		})
	}
}
