package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {

	// Create a new instance of our application struct. For now, this just
	// contains a couple of mock loggers (which discard anything written to
	// them).
	app := newTestApplication(t)

	// We then use the httptest.NewTLSServer() function to create a new test
	// server, passing in the value returned by our app.routes() method as the
	// handler for the server. This starts up a HTTPS server which listens on a
	// randomly-chosen port of your local machine for the duration of the test.
	// Notice that we defer a call to ts.Close() to shutdown the server when
	// the test finishes.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// The network address that the test server is listening on is contained
	// in the ts.URL field. We can use this along with the ts.Client().Get()
	// method to make a GET /ping request against the test server. This
	// returns a http.Response struct containing the response.
	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "Ok" {
		t.Errorf("want body to equal %q", "Ok")
	}
}

func TestShowSnippet(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	testCases := []struct {
		desc     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			code, _, body := ts.get(t, tC.urlPath)
			if code != tC.wantCode {
				t.Errorf("want %d; got %d", tC.wantCode, code)
			}
			if !bytes.Contains(body, tC.wantBody) {
				t.Errorf("want body to contain %q", tC.wantBody)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	// Create the application struct containing our mocked dependencies and set
	// up the test server for running and end-to-end test.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET /user/signup request and then extract the CSRF token from the
	// response body.
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)
	// Log the CSRF token value in our test output.
	// t.Log(csrfToken)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, []byte("")},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("")},
		{"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusOK, []byte("")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$w", csrfToken, http.StatusOK, []byte("")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$", csrfToken, http.StatusOK, []byte("")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("")},
		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, []byte("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)
			code, _, body := ts.postForm(t, "/user/signup", form)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}
