package middleware

import (
	"context"
	"errors"
	"fmt"
	kitjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	kpljwt "github.com/nsini/kplcloud/src/jwt"
	"strings"
)

var ErrorASD = errors.New("权限验证失败！")

func CheckAuthMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			//claims := ctx.Value(jwt.JWTClaimsContextKey)

			//fmt.Println(ctx.Value(kithttp.ContextKeyRequestPath))
			//fmt.Println(ctx.Value(kithttp.ContextKeyRequestProto))

			//fmt.Println("-------------> start... claims")
			//fmt.Println(claims)
			//fmt.Println("------------> claims end")

			token := ctx.Value(kithttp.ContextKeyRequestAuthorization).(string)
			if token == "" {
				return nil, ErrorASD
			}
			token = strings.Split(token, "Bearer ")[1]

			var clustom kpljwt.ArithmeticCustomClaims
			tk, err := kitjwt.ParseWithClaims(token, &clustom, kpljwt.JwtKeyFunc)
			if err != nil && tk == nil {
				_ = logger.Log("jwt", "ParseWithClaims", "err", err.Error())
			}

			claim, ok := tk.Claims.(*kpljwt.ArithmeticCustomClaims)
			if !ok {
				_ = logger.Log("tk", "Claims", "err", ok)
				err = ErrorASD
				return
			}

			// todo 数据库或cache验证权限？
			fmt.Println("userId: ", claim.UserId)
			fmt.Println("username: ", claim.Name)
			fmt.Println("request: ", request)

			return next(ctx, request)
		}
	}
}
