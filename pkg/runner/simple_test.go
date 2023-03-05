package runner

import (
	"errors"
	"net/http"
	"testing"

	_ "embed"
	"github.com/h2non/gock"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestTestCase(t *testing.T) {
	tests := []struct {
		name     string
		testCase *atest.TestCase
		ctx      interface{}
		prepare  func()
		verify   func(t *testing.T, output interface{}, err error)
	}{{
		name: "normal, response is map",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
				Header: map[string]string{
					"key": "value",
				},
				Body: `{"foo":"bar"}`,
			},
			Expect: atest.Response{
				StatusCode: http.StatusOK,
				BodyFieldsExpect: map[string]string{
					"name": "linuxsuren",
				},
				Header: map[string]string{
					"type": "generic",
				},
				Verify: []string{
					`name == "linuxsuren"`,
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").
				MatchHeader("key", "value").
				Reply(http.StatusOK).
				SetHeader("type", "generic").
				File("testdata/generic_response.json")
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.Nil(t, err)
			assert.Equal(t, map[string]interface{}{"name": "linuxsuren"}, output)
		},
	}, {
		name: "normal, response is slice",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				StatusCode: http.StatusOK,
				Body:       `["foo", "bar"]`,
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").
				Reply(http.StatusOK).
				BodyString(`["foo", "bar"]`)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.Nil(t, err)
			assert.Equal(t, []interface{}{"foo", "bar"}, output)
		},
	}, {
		name: "normal, response from file",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:          "http://localhost/foo",
				Method:       http.MethodPost,
				BodyFromFile: "testdata/generic_response.json",
			},
			Expect: atest.Response{
				StatusCode: http.StatusOK,
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Post("/foo").BodyString(genericBody).
				Reply(http.StatusOK).BodyString("123")
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "response from a not found file",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:          "http://localhost/foo",
				Method:       http.MethodPost,
				BodyFromFile: "testdata/fake.json",
			},
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "bad request",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				StatusCode: http.StatusOK,
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusBadRequest)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "error with request",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").ReplyError(errors.New("error"))
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "not match with body",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				Body: "bar",
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString("foo")
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "not match with header",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				Header: map[string]string{
					"foo": "bar",
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).SetHeader("foo", "value")
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "not found from fields",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				BodyFieldsExpect: map[string]string{
					"foo": "bar",
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString(genericBody)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "body filed not match",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				BodyFieldsExpect: map[string]string{
					"name": "bar",
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString(genericBody)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "invalid filed finding",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				BodyFieldsExpect: map[string]string{
					"items[1]": "bar",
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "failed to get field")
		},
	}, {
		name: "verify failed",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				Verify: []string{
					"len(items) > 0",
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "failed to verify")
		},
	}, {
		name: "failed to compile",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				Verify: []string{
					`println("12")`,
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "unknown name println")
		},
	}, {
		name: "failed to compile",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				Verify: []string{
					`1 + 1`,
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "expected bool, but got int")
		},
	}, {
		name: "wrong API format",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:    "ssh://localhost/foo",
				Method: "fake,fake",
			},
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "invalid method")
		},
	}, {
		name: "failed to render API",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo/{{.abc}",
			},
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "template: api:1:")
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()
			if tt.prepare != nil {
				tt.prepare()
			}
			output, err := RunTestCase(tt.testCase, tt.ctx)
			tt.verify(t, output, err)
		})
	}
}

//go:embed testdata/generic_response.json
var genericBody string