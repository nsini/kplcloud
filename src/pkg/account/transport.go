package account

import (
	"context"
	"encoding/json"
	"errors"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var errBadRoute = errors.New("bad route")

func MakeHandler(ps Service, logger kitlog.Logger) http.Handler {
	//ctx := context.Background()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	detail := kithttp.NewServer(
		makeDetailEndpoint(ps),
		decodeDetailRequest,
		encodeDetailResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/account/", detail).Methods("GET")
	return r
}

func decodeDetailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	postId, err := strconv.Atoi(id)
	if err != nil {
		return nil, errBadRoute
	}
	return accountRequest{
		Id: int64(postId),
	}, nil
}

func encodeDetailResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "blog-single")

	resp := response.(accountResponse)

	return json.NewEncoder(w).Encode(resp)
}

type errorer interface {
	error() error
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	//case repository.PostNotFound:
	//	w.WriteHeader(http.StatusNotFound)
	//	//ctx = context.WithValue(ctx, "method", "404")
	//	// _ = templates.RenderHtml(ctx, w, map[string]interface{}{})
	//	return
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, _ = w.Write([]byte(err.Error()))
}
