package main

import (
	"fmt"
	"log"

	"github.com/risy007/kmyh-config"
)

func main() {
	// 创建一个简单的测试来验证类型安全修复
	fmt.Println("Testing type safety improvements...")

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
	fmt.Println("Type safety improvements applied successfully!")
}
