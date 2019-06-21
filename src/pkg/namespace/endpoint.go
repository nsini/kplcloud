package namespace

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/kplcloud/src/repository"
)

type nsRequest struct {
	Name        string `json:"name_en,omitempty"`
	DisplayName string `json:"name,omitempty"`
}

type nsResponse struct {
	Code int                   `json:"code"`
	Data *repository.Namespace `json:"data,omitempty"`
	Err  error                 `json:"error,omitempty"`
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(nsRequest)
		var err error
		if rs, err := s.Detail(ctx, req.Name); err == nil {
			return nsResponse{0, rs, nil}, nil
		}
		return nsResponse{Code: -1, Err: err}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(nsRequest)
		var err error
		if err := s.Post(ctx, req.Name, req.DisplayName); err == nil {
			return nsResponse{0, nil, nil}, nil
		}
		return nsResponse{Code: -1, Err: err}, err
	}
}
