package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/nsini/kplcloud/src/config"
	kpljwt "github.com/nsini/kplcloud/src/jwt"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

var (
	ErrInvalidArgument  = errors.New("invalid argument")
	UserOrPasswordError = errors.New("username or password error.")
)

type Service interface {
	Login(ctx context.Context, username, password string) (rs string, err error)
}

type service struct {
	logger log.Logger
	config config.Config
}

/**
 * @Title 详情页
 */
func (c *service) Login(ctx context.Context, username, password string) (rs string, err error) {

	_ = c.logger.Log("name", username, "password", password)

	if username != "hello" && password != "world" {
		return "", UserOrPasswordError
	}

	rs, err = c.sign(username, "1")
	rs = "Bearer " + rs
	return
}

func (c *service) sign(username string, uid string) (string, error) {
	//为了演示方便，设置两分钟后过期
	sessionTimeout, err := strconv.Atoi(c.config.Get(config.SessionTimeout))
	if err != nil {
		sessionTimeout = 3600
	}
	expAt := time.Now().Add(time.Duration(sessionTimeout) * time.Second).Unix()

	// 创建声明
	claims := kpljwt.ArithmeticCustomClaims{
		UserId: uid,
		Name:   username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	//创建token，指定加密算法为HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//生成token
	return token.SignedString([]byte(kpljwt.GetJwtKey()))
}

func NewService(logger log.Logger, cf config.Config) Service {
	return &service{
		logger: logger,
		config: cf,
	}
}
