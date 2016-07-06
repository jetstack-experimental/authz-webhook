package main

import "encoding/json"
import "io"

type ResourceAttributesSpec struct {
	Namespace string `json:"namespace,omitempty"`
	Verb      string `json:"verb"`
	Group     string `json:"group,omitempty"`
	Resource  string `json:"resource"`
}

type NonResourceAttributesSpec struct {
	Path string `json:"path"`
	Verb string `json:"verb"`
}

type AuthorizationRequestSpec struct {
	NonResourceAttributes *NonResourceAttributesSpec `json:"nonResourceAttributes,omitempty"`
	ResourceAttributes    *ResourceAttributesSpec    `json:"resourceAttributes,omitempty"`
	User                  string                     `json:"user"`
	Group                 []string                   `json:"group,omitempty"`
}

type AuthorizationRequest struct {
	ApiVersion string                   `json:"apiVersion"`
	Kind       string                   `json:"kind"`
	Spec       AuthorizationRequestSpec `json:"spec"`
}

func NewAuthorizationRequest(body io.Reader) (*AuthorizationRequest, error) {
	var req *AuthorizationRequest

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&req)

	return req, err
}

func (r *AuthorizationRequest) Namespace() string {
	if !r.IsResourceRequest() {
		return ""
	}
	return r.Spec.ResourceAttributes.Namespace
}

func (r *AuthorizationRequest) IsResourceRequest() bool {
	return r.Spec.ResourceAttributes != nil
}

func (r *AuthorizationRequest) Action() string {
	if !r.IsResourceRequest() {
		return r.Spec.NonResourceAttributes.Verb
	}
	return r.Spec.ResourceAttributes.Verb
}

// Path is in NonResourceAttributes only
func (r *AuthorizationRequest) Path() string {
	if r.IsResourceRequest() {
		return ""
	}
	return r.Spec.NonResourceAttributes.Path
}

func (r *AuthorizationRequest) Group() string {
	if !r.IsResourceRequest() {
		return ""
	}
	return r.Spec.ResourceAttributes.Group
}

func (r *AuthorizationRequest) Resource() string {
	if !r.IsResourceRequest() {
		return ""
	}
	return r.Spec.ResourceAttributes.Resource
}
