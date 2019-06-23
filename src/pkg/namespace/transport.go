package namespace

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	kpljwt "github.com/nsini/kplcloud/src/jwt"
	"github.com/nsini/kplcloud/src/middleware"
	"io/ioutil"
	"net/http"
)

var errBadRoute = errors.New("bad route")

//type NamespaceEndpoints struct {
//	GetEndpoint  endpoint.Endpoint
//	PostEndpoint endpoint.Endpoint
//}

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	epsMap := map[string]endpoint.Endpoint{
		"get":    makeGetEndpoint(svc),
		"post":   makePostEndpoint(svc),
		"sync":   makeSyncEndpoint(svc),
		"update": makeUpdateEndpoint(svc),
	}

	for key, val := range epsMap {
		epsMap[key] = kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory)(val)
		epsMap[key] = middleware.CheckAuthMiddleware(logger)(val)
	}

	get := kithttp.NewServer(
		epsMap["get"],
		decodeGetRequest,
		encodeResponse,
		opts...,
	)

	create := kithttp.NewServer(
		epsMap["post"],
		decodePostRequest,
		encodeResponse,
		opts...,
	)

	sync := kithttp.NewServer(
		epsMap["sync"],
		decodeGetRequest,
		encodeResponse,
		opts...,
	)

	update := kithttp.NewServer(
		epsMap["update"],
		decodeGetRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/namespace/{name}", get).Methods("GET")
	r.Handle("/namespace/{name}", update).Methods("PUT")
	r.Handle("/namespace/", create).Methods("POST")
	r.Handle("/namespace/sync/all", sync).Methods("GET")

	//r.Handle("/namespace/{name}", ns).Methods("DELETE")

	return r
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, errBadRoute
	}
	return nsRequest{
		Name: name,
	}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req nsRequest

	if err = json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
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
	case kitjwt.ErrTokenContextMissing, kitjwt.ErrTokenExpired, middleware.ErrorASD:
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
