package main

type AuthorizationStatus struct {
  Allowed bool `json:"allowed"`
  Reason string `json:"reason,omitempty"`
}

type AuthorizationResponse struct {
  ApiVersion string `json:"apiVersion"`
  Kind string `json:"kind"`
  Status AuthorizationStatus `json:"status"`
}

func NewAuthorizationResponse(status bool, reason ...string) *AuthorizationResponse {
  resp := &AuthorizationResponse{
    ApiVersion: "authorization.k8s.io/v1beta1",
    Kind: "SubjectAccessReview",
    Status: AuthorizationStatus {
      Allowed: status,
    },
  }

  if len(reason) > 0 {
    resp.Status.Reason = reason[0]
  }
  return resp
}
