package namespace

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/yijizhichang/kplcloud/src/config"
	"github.com/yijizhichang/kplcloud/src/repository"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Detail(ctx context.Context, name string) (rs map[string]interface{}, err error)
}

type service struct {
	namespace repository.NamespaceRepository
	logger    log.Logger
	config    config.Config
}

/**
 * @Title 详情页
 */
func (c *service) Detail(ctx context.Context, name string) (rs map[string]interface{}, err error) {
	fmt.Println(ctx)
	fmt.Println(">>>>>>>>>>>>>>>>")

	//resp, err := c.namespace.Find(name)
	//if err != nil {
	//	return
	//}

	_ = c.logger.Log("name", "")

	return
}

func NewService(logger log.Logger, cf config.Config, namespace repository.NamespaceRepository) Service {
	return &service{
		logger:    logger,
		config:    cf,
		namespace: namespace,
	}
}
