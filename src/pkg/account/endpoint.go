package account

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type accountRequest struct {
	Id int64 `json:"id,omitempty"`
}

type accountResponse struct {
	Code int                    `json:"code,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
	Err  error                  `json:"error,omitempty"`
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(accountRequest)
		var err error
		if rs, err := s.Detail(ctx, req.Id); err == nil {
			return accountResponse{0, rs, nil}, nil
		}
		return accountResponse{Code: -1, Err: err}, err
	}
}
