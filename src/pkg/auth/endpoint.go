package auth

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type authRequest struct {
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	LoginType string `json:"login_type,omitempty"`
}

type authResponse struct {
	//Success bool   `json:"success"`
	Code  int    `json:"code"`
	Token string `json:"token"`
	Err   error  `json:"error"`
}

func (r authResponse) error() error { return r.Err }

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(authRequest)
		rs, err := s.Login(ctx, req.Username, req.Password)
		var code int
		if err != nil {
			code = 0
		}
		return authResponse{Code: code, Err: err, Token: rs}, nil
	}
}
