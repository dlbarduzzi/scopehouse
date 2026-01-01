package apis

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dlbarduzzi/scopehouse/tests"
)

type apiTestScenario struct {
	// name is the test name.
	name string

	// method is the HTTP method for the test request to use.
	method string

	// url is the url/path of the endpoint to be tested.
	url string

	// body specifies the body to send with the request.
	// i.e. strings.NewReader(`{"foo":"bar"}`)
	body io.Reader

	// headers specifies the headers to send with the request.
	headers map[string]string

	// extraRoute is an route that is not part of the API but that you
	// want to test. An example is a panic middleware where you explicitly
	// call this route for testing purposes.
	extraRoute *route

	// expectedStatus specifies the expected HTTP response status code.
	expectedStatus int

	// expectedContent specifies the list of keywords that must exist
	// in the response body. i.e. `{"foo":"bar"}`
	expectedContent []string

	// beforeTestFunc runs custom functions before running test cases.
	beforeTestFunc func(t testing.TB, app *tests.TestApp)
}

func (s *apiTestScenario) Test(t *testing.T) {
	t.Run(s.name, func(t *testing.T) {
		s.test(t)
	})
}

func (s *apiTestScenario) test(t *testing.T) {
	app, err := tests.NewTestApp()
	if err != nil {
		t.Fatalf("failed to initialize test app instance - %v", err)
	}

	router := newRouter(app)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(s.method, s.url, s.body)

	// Add support for inline testing routes.
	if s.extraRoute != nil {
		router.get(s.extraRoute.pattern, s.extraRoute.handler)
	}

	// Set default header.
	req.Header.Set("Content-Type", "application/json")

	for k, v := range s.headers {
		req.Header.Set(k, v)
	}

	if s.beforeTestFunc != nil {
		s.beforeTestFunc(t, app)
	}

	mux := router.buildMux()
	mux.ServeHTTP(rec, req)

	res := rec.Result()

	if res.StatusCode != s.expectedStatus {
		t.Fatalf(
			"expected status code to be %d, got %d",
			s.expectedStatus, res.StatusCode,
		)
	}

	if len(s.expectedContent) == 0 {
		if len(rec.Body.Bytes()) != 0 {
			t.Fatalf(
				"expected empty content, got \n%v",
				rec.Body.String(),
			)
		}
	} else {
		var body string
		buf := new(bytes.Buffer)

		err := json.Compact(buf, rec.Body.Bytes())
		if err != nil {
			// Not a json payload.
			body = rec.Body.String()
		} else {
			// A valid json payload.
			body = buf.String()
		}

		for _, content := range s.expectedContent {
			if !strings.Contains(body, content) {
				t.Errorf(
					"expected content %v in response body \n%v",
					content, body,
				)
			}
		}
	}
}
