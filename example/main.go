package main

import (
	"fmt"
	"log"

	"github.com/risy007/kmyh-config"
)

func main() {
	// 示例：模拟配置加载
	appConfig, err := config.NewAppConfig()
	if err != nil {
		log.Printf("Warning: Could not load config: %v", err)
		// 创建默认配置用于测试
		appConfig = &config.AppConfig{
			AppName: "test-app",
			Env:     "dev",
			Etcd: config.EtcdConfig{
				Endpoints:   []string{"localhost:2379"},
				DialTimeout: 5,
				Prefix:      "/config",
			},
			Logger: config.LogConfig{
				Level:     "info",
				Format:    "json",
				Directory: "./logs",
			},
		}
	}

	fmt.Printf("App Name: %s\n", appConfig.AppName)
	fmt.Printf("Environment: %s\n", appConfig.Env)

	// 演示新的泛型配置获取功能
	fmt.Println("\nTesting generic config retrieval...")

	// 注意：实际使用时需要有运行中的etcd服务器和相应的配置数据
	// 这里仅展示API使用方式
	fmt.Println("Generic GetConfig function added to the package")
	fmt.Println("Usage: config.GetConfig[config.DatabaseConfig](configManager, \"myapp\", \"prod\")")
	fmt.Println("This will automatically map to the 'database' config group")

	// 示例映射规则：
	fmt.Println("\nMapping examples:")
	fmt.Println("- DatabaseConfig -> database")
	fmt.Println("- HttpConfig -> http")
	fmt.Println("- DifyConfig -> dify")
	fmt.Println("- AppConfig -> app")
}
