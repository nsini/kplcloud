package config

type Config interface {
	Get(key string) string
}

const (
	ImageFilePath = "image_file_path"
	ImageDomain   = "image_domain"
	SessionTimeout = "session_timeout"
)

type config struct {
}

func NewConfig(path string) Config {
	// 处理配置文件
	return &config{}
}

func (c *config) Get(key string) string {

	switch key {
	case "image_domain":
		return "http://source.lattecake.com/"
	case "image_file_path":
		return "./image/"
	case "session_timeout":
		return "3600"
	}

	return ""
}
