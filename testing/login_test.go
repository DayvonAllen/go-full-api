package testing

import (
	"example.com/app/router"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginRoute(t *testing.T) {
	// Define a structure for specifying input and output
	// data of a single test case. This structure is then used
	// to create a so called test map, which contains all test
	// cases, that should be run for testing this function
	tests := []struct {
		description string

		// Test input
		route string
		reqBody *strings.Reader

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "login route",
			route:         "/auth/login",
			reqBody: strings.NewReader(`{"email": "jdoedddd25455@gmail.com", "password": "password"}`),
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "OK",
		},
		{
			description:   "Unauthorized user",
			route:         "/auth/login",
			reqBody: strings.NewReader(`{"email": "jdoedddd255@gmail.com", "password": "password"}`),
			expectedError: false,
			expectedCode:  401,
			expectedBody:  "Unauthorized",
		},
	}

	// Setup the app as it is done in the main function
	app := router.Setup()

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"POST",
			test.route,
			test.reqBody,
			)

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the response body equals the expected body
		assert.Equalf(t, test.expectedBody, string(body), test.description)
	}
}
