//go:build unit

package testutil

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestContext holds a gin context and response recorder for handler tests.
type TestContext struct {
	Engine  *gin.Engine
	Context *gin.Context
	Recorder *httptest.ResponseRecorder
}

// NewTestContext creates a fresh TestContext for a given method and path.
func NewTestContext(method, path string) *TestContext {
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(method, path, nil)
	return &TestContext{
		Engine:   gin.Default(),
		Context:  context,
		Recorder: recorder,
	}
}

// NewTestContextWithJSON creates a TestContext with a JSON body.
func NewTestContextWithJSON(method, path string, body any) *TestContext {
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			panic("failed to marshal test body: " + err.Error())
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	context.Request = httptest.NewRequest(method, path, reqBody)
	if body != nil {
		context.Request.Header.Set("Content-Type", "application/json")
	}
	return &TestContext{
		Engine:   gin.Default(),
		Context:  context,
		Recorder: recorder,
	}
}

// SetHeader sets a request header.
func (tc *TestContext) SetHeader(key, value string) {
	tc.Context.Request.Header.Set(key, value)
}

// SetQueryParam sets a query parameter.
func (tc *TestContext) SetQueryParam(key, value string) {
	q := tc.Context.Request.URL.Query()
	q.Set(key, value)
	tc.Context.Request.URL.RawQuery = q.Encode()
}

// SetPathParam sets a path parameter (gin style).
func (tc *TestContext) SetPathParam(key, value string) {
	tc.Context.Params = append(tc.Context.Params, gin.Param{Key: key, Value: value})
}

// SetAuthHeader sets the Authorization header with a Bearer token.
func (tc *TestContext) SetAuthHeader(token string) {
	tc.SetHeader("Authorization", "Bearer "+token)
}

// JSONBody decodes the response body as JSON into the target.
func (tc *TestContext) JSONBody(target any) error {
	return json.Unmarshal(tc.Recorder.Body.Bytes(), target)
}

// BodyString returns the response body as a string.
func (tc *TestContext) BodyString() string {
	return tc.Recorder.Body.String()
}

// StatusCode returns the response status code.
func (tc *TestContext) StatusCode() int {
	return tc.Recorder.Code
}

// AssertStatus checks that the response status code matches expected.
func (tc *TestContext) AssertStatus(t *testing.T, expected int) {
	t.Helper()
	if tc.Recorder.Code != expected {
		t.Errorf("expected status %d, got %d. Body: %s", expected, tc.Recorder.Code, tc.Recorder.Body.String())
	}
}

// AssertJSONField checks that a JSON response field matches expected value.
func (tc *TestContext) AssertJSONField(t *testing.T, field string, expected any) {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal(tc.Recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}
	actual, ok := result[field]
	if !ok {
		t.Errorf("expected JSON field %q not found in response", field)
		return
	}
	// Use JSON round-trip for comparison to handle type differences
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)
	if string(expectedJSON) != string(actualJSON) {
		t.Errorf("field %q: expected %s, got %s", field, expectedJSON, actualJSON)
	}
}

// PerformRequest executes a handler function and records the result.
func (tc *TestContext) PerformRequest(handler gin.HandlerFunc) {
	handler(tc.Context)
}
