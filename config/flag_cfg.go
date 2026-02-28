package config

import (
	"dbgod_workOrder_job/loggerzap"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type arrayFlags []string

func (f *arrayFlags) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *arrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

// 初始化全局配置
func init() {
	// 初始化zap配置
	logCfg := &loggerzap.LogConfig{
		Level:      "info",
		FilePath:   "logs/run.log",
		MaxSize:    100, // 100MB
		MaxBackups: 20,  // 控制保留的历史文件数量上限
		//MaxAge:   5, // 日志文件的最大保留天数
		Compress: false,
	}

	// 初始化日志
	if err := loggerzap.Init("production", logCfg); err != nil {
		loggerzap.L().Panic("init production failed", zap.Error(err))
	}
	defer loggerzap.Sync()

	// 从命令行获取flag参数
	var Configs arrayFlags
	flag.Var(&Configs, "conf", "config files path")
	// 通过Parse解析获取到配置文件
	flag.Parse()
	if len(Configs) == 0 {
		loggerzap.L().Panic("config file not specified.")
	}

	//viper.SetConfigName("slow_config") //指定配置文件的文件名称(不需要制定配置文件的扩展名)
	// 配置文件名称来源于命令行参数--conf
	viper.SetConfigName(Configs[0])

	//设置配置文件类型
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config") // 设置配置文件的搜索目录
	viper.AddConfigPath("../../config")

	viper.AutomaticEnv()
	// 如果上面路径找不到，读取环境变量的值，export SLOWCONF=/opt/go-slow-log/config
	viper.AddConfigPath(viper.GetString("SLOWCONF"))
	if err := viper.ReadInConfig(); err != nil {
		loggerzap.L().Panic("read config file wrong", zap.Error(err))
	}

	err := viper.Unmarshal(&_config) // 将配置信息绑定到结构体上
	if err != nil {
		loggerzap.L().Panic("unmarshal config file wrong", zap.Error(err))
	}

	//fmt.Println(_config)
	//viper.WatchConfig()
	////可以通过https://fsnotify.org 监听config文件变化更新配置信息
	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("配置发生变更：", e.Name)
	//})
}

// 获取全局配置
func GetConfig() *Config {
	return _config
}
