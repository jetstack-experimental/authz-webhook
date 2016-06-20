package main

import "strings"

type saccountError struct{}

func (*saccountError) Error() string {
	return "Not a Service Account"
}

type rtypeError struct{}

func (*rtypeError) Error() string {
	return "Wrong request type"
}

type AuthzUser struct {
	saccount []string
	request  *AuthorizationRequest
}

func stripLastPart(str string, sep string) string {
	splitted := strings.Split(str, sep)
	if len(splitted) > 1 {
		return strings.Join(splitted[:len(splitted)-1], sep)
	}
	return str
}

// NewAuthzUser reuturns new AuthzUser struct
func NewAuthzUser(req *AuthorizationRequest) *AuthzUser {
	saccountData := strings.Split(req.Spec.User, ":")
	return &AuthzUser{
		saccount: saccountData,
		request:  req,
	}
}

// IsAllowed checks if service account can access resource
// returns true on success, false otherwise
func (r *AuthzUser) IsAllowed() bool {
	// allow kubectl auto-detect requests:
	// 'get' to /apis
	if (r.request.Path() == "/apis" || r.request.Path() == "/api") && r.request.Action() == "get" {
		return true
	}

	// allow all 'list', 'watch' actions to kube-system service account
	if (r.Username() == "system:serviceaccount:kube-system:default") &&
		(r.request.Action() == "list" || r.request.Action() == "watch") {
		return true
	}

	userNamespace, err := r.serviceAccountNamespace()
	if err != nil {
		return false
	}

	actionNamespace := r.request.Namespace()
	if actionNamespace == "" {
		return false
	}
	// We allow access for namespace-${anything} user to namespace-${anything}
	strippedUserNamespace := stripLastPart(userNamespace, "-")
	strippedActionNamespace := stripLastPart(actionNamespace, "-")
	return strippedUserNamespace == strippedActionNamespace
}

// IsServiceAccount returns true if user is service account
func (r *AuthzUser) IsServiceAccount() bool {
	if len(r.saccount) == 4 && r.saccount[0] == "system" && r.saccount[1] == "serviceaccount" {
		return true
	}
	return false
}

// Username returns request's spec.user
func (r *AuthzUser) Username() string {
	return r.request.Spec.User
}

// Namespace returns request's namespace
func (r *AuthzUser) Namespace() string {
	return r.request.Namespace()
}

// Request returns full request struct
func (r *AuthzUser) Request() *AuthorizationRequest {
	return r.request
}

// NamespaceRequest returns true if request is made for namespace
func (r *AuthzUser) NamespaceRequest() bool {
	return r.request.IsResourceRequest()
}

func (r *AuthzUser) serviceAccountNamespace() (string, error) {
	if !r.IsServiceAccount() {
		return "", &saccountError{}
	}
	return r.saccount[2], nil
}
