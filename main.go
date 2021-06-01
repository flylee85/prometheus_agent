package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/signal"
	"prometheus_agent/config"
	"syscall"
	"time"
)

func initConfig(path string) *config.AgentConfig {
	var config config.AgentConfig

	// 环境变量读取
	viper.AutomaticEnv()
	viper.SetEnvPrefix("PROM_AGENT")

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		logrus.Fatal(err)
	}
	return &config
}

func initLog(verbose bool, config *config.AgentConfig) {
	logger := &lumberjack.Logger{
		Filename:   config.LogConfig.Filename,
		MaxSize:    config.LogConfig.Maxsize,
		MaxBackups: config.LogConfig.Maxbackups,
		Compress:   config.LogConfig.Compress,
	}
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(logger)
	}
}

func main() {
	var (
		verbose, h, help bool
		path             string
	)
	flag.BoolVar(&verbose, "verbose", false, "verbose")
	flag.BoolVar(&h, "h", false, "h")
	flag.BoolVar(&help, "help", false, "help")
	flag.StringVar(&path, "path", "./etc/promagent.yaml", path)

	flag.Usage = func() {
		fmt.Println("Usage: promagent [--verbose] [--config file]")
		flag.PrintDefaults()
	}

	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if h {
		flag.Usage()
		os.Exit(0)
	}

	// 初始化配置文件
	config := initConfig(path)
	fmt.Println(config)
	// 初始化日志
	initLog(verbose, config)
	// 启动
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGKILL, syscall.SIGINT)

	reload := make(chan os.Signal, 1)
	signal.Notify(reload, syscall.SIGHUP)

	go func() {
		for {
			logrus.Debug("doing")
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			<-reload
			fmt.Println("reload")
		}
	}()

	<-stop

}
