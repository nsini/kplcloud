package namespace

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/kplcloud/src/config"
	"github.com/nsini/kplcloud/src/kubernetes"
	"github.com/nsini/kplcloud/src/repository"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrNamespaceIsExists = errors.New("空间已经存在.")
	ErrNamespaceCreate   = errors.New("空间创建失败.")
)

type Service interface {
	Detail(ctx context.Context, name string) (resp *repository.Namespace, err error)
	Post(ctx context.Context, name, displayName string) error
}

type service struct {
	namespace repository.NamespaceRepository
	logger    log.Logger
	config    config.Config
	k8sClient kubernetes.K8sClient
}

/**
 * @Title 详情页
 */
func (c *service) Detail(ctx context.Context, name string) (resp *repository.Namespace, err error) {
	return c.namespace.Find(name)
}

/**
 * @Title 创建namespace
 */
func (c *service) Post(ctx context.Context, name, displayName string) error {
	res, err := c.namespace.Find(name)
	if err != nil {
		_ = c.logger.Log("displayName", res.Name)
		return ErrNamespaceIsExists
	}

	namespace := new(v1.Namespace)
	namespace.Name = name

	if _, err := c.k8sClient.Do().CoreV1().Namespaces().Create(namespace); err != nil {
		_ = c.logger.Log("k8s", "create", "err", err.Error())
		return ErrNamespaceCreate
	}

	if err = c.namespace.Create(&repository.Namespace{
		Name:   displayName,
		NameEn: name,
	}); err != nil {
		_ = c.logger.Log("ns", "create", "err", err.Error())
		return ErrNamespaceCreate
	}

	// 是否需要创建secret
	if c.config.Get(config.K8sImagePullSecrets) != "" {
		if _, err := c.k8sClient.Do().CoreV1().Secrets(name).Create(&v1.Secret{
			Type: v1.SecretTypeDockerConfigJson,
			ObjectMeta: metav1.ObjectMeta{
				Name:      c.config.Get(config.K8sImagePullSecrets),
				Namespace: name,
			},
		}); err != nil {
			_ = c.logger.Log("k8s", "secrets", "err", err.Error())
		}
	}

	// todo 创建jenkins 视图

	return nil
}

func NewService(logger log.Logger, cf config.Config, namespace repository.NamespaceRepository, client kubernetes.K8sClient) Service {
	return &service{
		logger:    logger,
		config:    cf,
		k8sClient: client,
		namespace: namespace,
	}
}
