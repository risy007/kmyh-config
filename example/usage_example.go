package main

import (
	"fmt"
	"log"

	"github.com/risy007/kmyh-config"
)

func main() {
	// 加载应用配置
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

	// 创建日志记录器
	logger := config.NewZapLogger(appConfig.Logger)

	// 创建etcd客户端
	etcdClient, err := config.NewEtcdClient(appConfig.Etcd, logger)
	if err != nil {
		log.Fatalf("Failed to create etcd client: %v", err)
	}
	defer etcdClient.Close()

	// 创建 ConfigManager
	configManager := config.NewConfigManagerDirect(etcdClient, logger, appConfig)

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

	// 实际使用示例（需要etcd中有对应配置才能成功）
	dbConfig, err := config.GetConfig[config.DatabaseConfig](configManager, appConfig.AppName, appConfig.Env)
	if err != nil {
		fmt.Printf("Could not load database config (expected if etcd not running): %v\n", err)
	} else {
		fmt.Printf("Database config loaded: %+v\n", dbConfig)
	}
}
