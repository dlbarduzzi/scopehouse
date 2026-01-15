package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dlbarduzzi/scopehouse/internal/tests"
)

type apiTestRoute struct {
	pattern string
	handler func(http.ResponseWriter, *http.Request)
}

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
	testRoute *apiTestRoute

	// status specifies the expected response HTTP status code.
	status int

	// content specifies the list of keywords that must exist
	// in the response body. i.e. `{"foo":"bar"}` - leave empty if you
	// want to ensure that the response didn't have any body (e.g. 204).
	content []string

	// beforeTestFunc runs custom functions before running test cases.
	beforeTestFunc func(t testing.TB, app *tests.TestApp)
}

func (s *apiTestScenario) test(t *testing.T) {
	t.Run(s.normalizeName(), func(t *testing.T) {
		s._test(t)
	})
}

func (s *apiTestScenario) _test(t *testing.T) {
	app, err := tests.NewTestApp()
	if err != nil {
		t.Fatalf("failed to initialize test app instance; %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(s.method, s.url, s.body)

	// Set default header.
	req.Header.Set("Content-Type", "application/json")

	for k, v := range s.headers {
		req.Header.Set(k, v)
	}

	if s.beforeTestFunc != nil {
		s.beforeTestFunc(t, app)
	}

	svc := newService(app)
	mux := svc.routes()

	if (s.testRoute) != nil {
		mux.HandleFunc(s.testRoute.pattern, s.testRoute.handler)
	}

	handler := svc.handler(mux)
	handler.ServeHTTP(rec, req)

	res := rec.Result()

	if res.StatusCode != s.status {
		t.Fatalf(
			"expected status code to be %d, got %d",
			s.status, res.StatusCode,
		)
	}

	testBodyContent(t, rec, s.content)
}

func (s *apiTestScenario) normalizeName() string {
	name := strings.TrimSpace(s.name)

	if name == "" {
		name = fmt.Sprintf("%s:%s", s.method, s.url)
	}

	return name
}

func testBodyContent(t *testing.T, rec *httptest.ResponseRecorder, content []string) {
	bodyBytes := rec.Body.Bytes()

	if len(content) == 0 {
		if len(bodyBytes) != 0 {
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

		for _, item := range content {
			if !strings.Contains(body, item) {
				t.Fatalf(
					"expected content `%v` in response body \n%v",
					item, body,
				)
			}
		}
	}
}
