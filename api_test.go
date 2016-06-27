package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server *httptest.Server
	reader io.Reader
	index  string
)

func init() {
	server = httptest.NewServer(Handlers())
	index = fmt.Sprintf("%s/", server.URL)
}

func postIndex(data string) (*http.Response, error) {
	reader = strings.NewReader(data)

	request, err := http.NewRequest("POST", index, reader)
	res, err := http.DefaultClient.Do(request)
	return res, err

}

func parseResponse(r *http.Response) (*AuthorizationResponse, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	respJSON := NewAuthorizationResponse(true)
	err = json.Unmarshal(body, &respJSON)
	return respJSON, err
}

func TestNewAuthorizationByNamespace(t *testing.T) {
	var saTests = []struct {
		namespace string
		username  string
		status    int
		allowed   bool
	}{
		{"namespace-dev", "system:serviceccount:namespace-sss:default", 200, true},
		{"namespace", "system:serviceaccount:namespace-sss:default", 200, true},
		{"namespace-dev", "someuser", 200, true},
		{"namespace-dev", "system:serviceaccount:default-ns:default", 403, false},
	}

	for _, tst := range saTests {

		reqJSON := fmt.Sprintf(`
    {
      "spec":{
        "resourceAttributes": {
          "namespace":"%s",
          "verb": "get"
        },
        "user":"%s"
      }
    }`, tst.namespace, tst.username)

		result, err := postIndex(reqJSON)

		if err != nil {
			t.Error(err)
		}

		if result.StatusCode != tst.status {
			t.Errorf("Expected status %d, got: %d (%s)", tst.status, result.StatusCode, reqJSON)
		}

		respJSON, err := parseResponse(result)
		if err != nil {
			t.Error(err)
		}

		if respJSON.Status.Allowed != tst.allowed {
			t.Errorf("Bad response status, expected %t got: %t (%s)", tst.allowed, respJSON.Status.Allowed, reqJSON)
		}
	}
}

func TestNewAuthorizationBadRequest(t *testing.T) {
	reqJSON := `asd`
	result, err := postIndex(reqJSON)

	if err != nil {
		t.Error(err)
	}

	if result.StatusCode != 400 {
		t.Errorf("Bad request expected: %d", result.StatusCode)
	}
}

func TestNewAuthorizationByPath(t *testing.T) {
	var saTests = []struct {
		path     string
		verb     string
		username string
		status   int
	}{
		{"/apis", "get", "system:serviceaccount:random:default", 200},
		{"/api", "get", "system:serviceaccount:random:default", 200},
		{"/version", "get", "system:serviceaccount:random:default", 200},
		{"/api/v1", "get", "system:serviceaccount:random:default", 403},
	}

	for _, tst := range saTests {
		reqJSON := fmt.Sprintf(`
    {
      "spec":{
        "nonResourceAttributes":{
          "path": "%s",
          "verb": "%s"
        },
        "user": "%s"
      }
    }`, tst.path, tst.verb, tst.username)
		result, err := postIndex(reqJSON)

		if err != nil {
			t.Error(err)
		}
		if result.StatusCode != tst.status {
			t.Errorf("Expected %d, got: %d, (%s)", tst.status, result.StatusCode, reqJSON)
		}
	}
}

func TestNewAuthorizationByVerb(t *testing.T) {

	var saTests = []struct {
		namespace string
		verb      string
		username  string
		status    int
	}{
		{"default", "watch", "system:serviceaccount:kube-system:default", 200},
		{"default", "watch", "system:serviceaccount:random:default", 403},
		{"default", "list", "system:serviceaccount:kube-system:default", 200},
		{"default", "list", "system:serviceaccount:random:default", 403},
	}

	for _, tst := range saTests {

		reqJSON := fmt.Sprintf(`
    {
      "spec":{
        "resourceAttributes":{
          "namespace":"%s",
          "verb":"%s",
          "resource": "services"
        },
        "user":"%s"
      }
    }`, tst.namespace, tst.verb, tst.username)
		result, err := postIndex(reqJSON)

		if err != nil {
			t.Error(err)
		}
		if result.StatusCode != tst.status {
			t.Errorf("Expected status %d, got: %d (%s)", tst.status, result.StatusCode, reqJSON)
		}
	}
}
