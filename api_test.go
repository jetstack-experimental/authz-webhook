package main

import (
  "testing"
  "fmt"
  "strings"
  "encoding/json"
  "io"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
)

var (
  server   *httptest.Server
  reader   io.Reader
  index    string
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
  if err != nil { return nil, err }

  respJson := NewAuthorizationResponse(true)
  err = json.Unmarshal(body, &respJson)
  return respJson, err
}

// This should permit access for serviceaccount
func TestNewAuthorizationRequestPermitSA(t *testing.T) {
  reqJson := `{"spec":{ "resourceAttributes": {"namespace":"namespace-dev"}},"spec":{"user":"system:serviceaccount:namespace-sss:default"}}'`
  result, err := postIndex(reqJson)

  if err != nil { t.Error(err) }

  if result.StatusCode != 200 {
    t.Errorf("Success expected: %d", result.StatusCode)
  }
}

// Should permit end-user straight away
func TestNewAuthorizationRequestPermitUser(t *testing.T) {
  reqJson := `{"spec":{ "resourceAttributes": {"namespace":"namespace-dev"}},"spec": {"user":"someuser"}}'`
  result, err := postIndex(reqJson)

  if err != nil { t.Error(err) }

  if result.StatusCode != 200 { t.Errorf("Success expected: %d", result.StatusCode) }
}

// Should deny serviceaccount on namespace mismatch away
func TestNewAuthorizationRequestDenySA(t *testing.T) {
  reqJson := `{"spec":{ "resourceAttributes": {"namespace":"namespace-dev"}},"spec": {"user":"system:serviceaccount:default:default"}}'`
  result, err := postIndex(reqJson)

  if err != nil { t.Error(err) }

  if result.StatusCode != 403 { t.Errorf("Forbidden expected: %d", result.StatusCode) }

  respJson, err := parseResponse(result)
  if err != nil { t.Error(err) }

  if respJson.Status.Allowed != false {
    t.Errorf("Bad response status: %t", respJson.Status.Allowed)
  }

}

// Should generate status 400 if bad request
func TestNewAuthorizationRequestBadRequest(t *testing.T) {
  reqJson := `asd`
  result, err := postIndex(reqJson)

  if err != nil { t.Error(err) }

  if result.StatusCode != 400 { t.Errorf("Bad request expected: %d", result.StatusCode) }
}

// Test serviceaccount '/apis' path
func TestNewAuthorizationRequestApisPath(t *testing.T) {
  reqJson := `{"spec":{"nonResourceAttributes":{"path":"/apis","verb":"get"},"user":"system:serviceaccount:random:default"}}`
  result, err := postIndex(reqJson)

  if err != nil { t.Error(err) }
  if result.StatusCode != 200 { t.Errorf("Success expected: %d", result.StatusCode)}
}

// Test kube-system serviceaccount 'watch' verb
func TestNewAuthorizationRequestWatchVerb(t *testing.T) {
  reqJson := `{"spec":{"resourceAttributes":{"namespace":"kube-system","verb":"watch", "resource": "services"},"user":"system:serviceaccount:kube-system:default"}}`
  result, err := postIndex(reqJson)

  if err != nil { t.Error(err) }
  if result.StatusCode != 200 { t.Errorf("Success expected: %d", result.StatusCode)}
}

// Test kube-system serviceaccount 'watch' verb
func TestNewAuthorizationRequestListVerb(t *testing.T) {
  reqJson := `{"spec":{"resourceAttributes":{"namespace":"kube-system","verb":"list", "resource": "nodes"},"user":"system:serviceaccount:kube-system:default"}}`
  result, err := postIndex(reqJson)

  if err != nil { t.Error(err) }
  if result.StatusCode != 200 { t.Errorf("Success expected: %d", result.StatusCode)}
}
