package namespace

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	kpljwt "github.com/yijizhichang/kplcloud/src/jwt"
	"github.com/yijizhichang/kplcloud/src/middleware"
	"net/http"
)

var errBadRoute = errors.New("bad route")

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		//kithttp.ServerBefore(kpljwt.JwtRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	ep := makeDetailEndpoint(svc)
	ep = kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory)(ep)
	ep = middleware.CheckAuthMiddleware(logger)(ep)

	get := kithttp.NewServer(
		ep,
		decodeGetRequest,
		encodeGetResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/namespace/{name}", get).Methods("GET")

	//r.Handle("/namespace/{name}", ns).Methods("DELETE")
	//r.Handle("/namespace/create", create).Methods("POST")
	//r.Handle("/namespace/{name}", put).Methods("PUT")

	return r
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, errBadRoute
	}
	//fmt.Println(r.Header.Get("Authorization"))

	return nsRequest{
		Name: name,
	}, nil
}

func encodeGetResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case kitjwt.ErrTokenContextMissing, kitjwt.ErrTokenExpired:
		w.WriteHeader(http.StatusForbidden)
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  -1,
		"error": err.Error(),
	})
}
