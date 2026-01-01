package apis

import (
	"testing"

	"github.com/dlbarduzzi/scopehouse/core"
)

func TestMiddlewaresPanicRecover(t *testing.T) {
	t.Parallel()

	testCases := []apiTestScenario{
		{
			name:           "panic recover middleware",
			url:            "/force/panic",
			expectedStatus: 500,
			expectedContent: []string{
				`"status":500`,
				`"message":"Something went wrong while processing this request."`,
			},
			extraRoute: &route{
				pattern: "/force/panic",
				handler: func(_ *core.EventRequest) {
					panic(123)
				},
			},
		},
	}

	for _, tc := range testCases {
		tc.Test(t)
	}
}
