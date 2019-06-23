package config

type Config interface {
	Get(key string) string
}

type configConst string

const (
	ImageFilePath              = "image_file_path"
	ImageDomain                = "image_domain"
	SessionTimeout             = "session.timeout"
	DeploymentImagePullSecrets = "image_pull_secrets"
)

type config struct {
}

func NewConfig(path string) Config {
	// 处理配置文件
	return &config{}
}

func (c *config) Get(key string) string {

	switch configConst(key) {
	case SessionTimeout:
		return "3600"
	case DeploymentImagePullSecrets:
		return "regcred"
	}

	return ""
}
