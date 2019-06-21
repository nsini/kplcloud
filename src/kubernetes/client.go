package kubernetes

import (
	"github.com/nsini/kplcloud/src/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

const (
	// High enough QPS to fit all expected use cases.
	defaultQPS = 1e2
	// High enough Burst to fit all expected use cases.
	defaultBurst = 1e2
	// full resyc cache resource time
	defaultResyncPeriod = 30 * time.Second
)

type K8sClient interface {
	Do() *kubernetes.Clientset
	Config() *rest.Config
}

type client struct {
	clientSet *kubernetes.Clientset
	config    *rest.Config
}

func NewClient(cf *config.Config) (cli K8sClient, err error) {
	// todo 怎么连呢？ 这里需要调整
	cliConfig, err := clientcmd.BuildConfigFromFlags("", "./config.yaml")
	if err != nil {
		return
	}

	cliConfig.QPS = defaultQPS
	cliConfig.Burst = defaultBurst
	cliConfig.Timeout = time.Second * 10

	// create the clientset
	clientset, err := kubernetes.NewForConfig(cliConfig)
	if err != nil {
		return
	}

	return &client{clientSet: clientset, config: cliConfig}, nil
}

func (c *client) Do() *kubernetes.Clientset {
	return c.clientSet
}

func (c *client) Config() *rest.Config {
	return c.config
}
