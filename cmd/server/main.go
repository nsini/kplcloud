package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/yijizhichang/kplcloud/src/cmd"
	"github.com/yijizhichang/kplcloud/src/config"
	"github.com/yijizhichang/kplcloud/src/pkg/auth"
	"github.com/yijizhichang/kplcloud/src/pkg/namespace"
	"github.com/yijizhichang/kplcloud/src/repository"
	"github.com/yijizhichang/kplcloud/src/repository/initialization"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	DefaultHttpPort   = ":8080"
	DefaultConfigPath = "./config/app.yaml"
	DefaultStaticPath = "./static/"
)

var (
	httpAddr   = envString("HTTP_ADDR", DefaultHttpPort)
	configPath = envString("CONFIG_PATH", DefaultConfigPath)
	staticPath = envString("STATIC_PATH", DefaultStaticPath)

	cf config.Config
)

var (
	rootCmd = &cobra.Command{
		Use:               "server",
		Short:             "开普勒平台服务端",
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Long: `# 开普勒平台服务端

您可以通过改命令来启动您的服务

可用的配置类型：

[start]

有关开普勒平台的相关概述，请参阅 https://github.com/yijizhichang/kplcloud
`,
	}

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "启动服务",
		Example: `## 启动命令
server start -p :8080 -c ./config/app.yaml
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			run()
			return nil
		},
	}
)

func init() {

	rootCmd.PersistentFlags().StringVarP(&httpAddr, "http.port", "p", DefaultHttpPort, "服务启动的端口: :8080")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config.path", "c", DefaultConfigPath, "配置文件路径: ./config/app.yaml")
	startCmd.PersistentFlags().StringVarP(&staticPath, "static.path", "s", DefaultStaticPath, "静态文件目录: ./static/")

	cmd.AddFlags(rootCmd)
	rootCmd.AddCommand(startCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func run() {
	cf = config.NewConfig(configPath)

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.StdlibWriter{})
	logger = log.With(logger, "caller", log.DefaultCaller)

	db, err := initialization.NewDb(logger, cf)
	if err != nil {
		_ = logger.Log("db", "connect", "err", err)
		panic(err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	var (
		namespaceRepository = repository.NewNamespaceRepository(db)
	)

	fieldKeys := []string{"method"}

	// namespace service
	{

		var authSvc auth.Service
		authSvc = auth.NewService(logger, cf)
		authSvc = auth.NewLoggingService(logger, authSvc)
		authSvc = auth.NewInstrumentingService(
			prometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "auth",
				Subsystem: "auth_service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "auth",
				Subsystem: "auth_service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys), authSvc)

		var namespaceSvc namespace.Service
		namespaceSvc = namespace.NewService(logger, cf, namespaceRepository)
		namespaceSvc = namespace.NewLoggingService(logger, namespaceSvc) // 日志

		httpLogger := log.With(logger, "component", "http")

		mux := http.NewServeMux()

		http.Handle("/login", auth.MakeHandler(authSvc, httpLogger))
		http.Handle("/namespace/", namespace.MakeHandler(namespaceSvc, httpLogger))

		http.Handle("/metrics", promhttp.Handler())
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
		http.Handle("/", accessControl(mux, logger))
	}

	errs := make(chan error, 2)
	go func() {
		_ = logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	_ = logger.Log("terminated", <-errs)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func accessControl(h http.Handler, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		_ = logger.Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)

		h.ServeHTTP(w, r)
	})
}
