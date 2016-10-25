package main

type AuthzUser struct {
	context *RequestContext
}

// NewAuthzUser reuturns new AuthzUser struct
func NewAuthzUser(req *AuthorizationRequest) *AuthzUser {
	context := NewRequestContext(req)

	return &AuthzUser{
		context: context,
	}
}

// IsAllowed checks if service account can access resource
// returns true on success, false otherwise
func (r *AuthzUser) IsAllowed() bool {
	for _, entry := range config.Rules {
		if entry.Matches(r.context) {
			return true
		}
	}
	return false
}

// Username returns request's spec.user
func (r *AuthzUser) Username() string {
	return r.context.Username
}

// Request returns full request struct
func (r *AuthzUser) Request() *AuthorizationRequest {
	return r.context.Request
}
