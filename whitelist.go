package main

// UserAction defines whitelisted action for a particular user.
type UserAction struct {
	Username  string
	Verb      string
	Group     string
	Resource  string
	Namespace string
	Path      string
}

// Whitelist of allowed actions.
// Empty field indicates any request value is allowed. E.g.
//  { Namespace: "a"} will allow everything within "a" namespace
var Whitelist = []UserAction{
	{
		Username: "system:serviceaccount:kube-system:default",
		Verb:     "list",
	},
	{
		Username: "system:serviceaccount:kube-system:default",
		Verb:     "watch",
	},
	{
		Username: "system:serviceaccount:kube-system:default",
		Verb:     "get",
	},
	{
		Verb:     "create",
		Group:    "extensions",
		Resource: "thirdpartyresources",
	},
	{
		Path: "/api",
	},
	{
		Path: "/apis",
	},
	{
		Path: "/version",
	},
	{
		Path: "/swaggerapi/api/v1",
		Verb: "get",
	},
	{
		Path: "/swaggerapi/apis/extensions/v1beta1",
		Verb: "get",
	},
}

// Matches is true if username is allowed to perform action specified in AuthorizationRequest
func (s *UserAction) Matches(username string, request *AuthorizationRequest) bool {
	return match(s.Username, username) &&
		match(s.Verb, request.Action()) &&
		match(s.Group, request.Group()) &&
		match(s.Resource, request.Resource()) &&
		match(s.Namespace, request.Namespace()) &&
		match(s.Path, request.Path())
}

func match(field string, str string) bool {
	return field == "" || field == "*" || field == str
}
