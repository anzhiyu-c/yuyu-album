/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-04-11 06:02:54
 * @LastEditTime: 2025-06-15 12:30:16
 * @LastEditors: 安知鱼
 */
package database

import (
	"context" // 用于传递上下文
	"log"     // 用于日志记录
	"strconv" // 用于将字符串转换为整数

	"album-admin/config" // 你的配置包，确保路径正确

	"github.com/redis/go-redis/v9" // Redis 客户端库
)

// RDB 是全局的 Redis 客户端实例
var RDB *redis.Client

// Ctx 是一个全局的后台上下文，通常用于 Redis 操作
var Ctx = context.Background()

// InitRedis 初始化 Redis 连接
func InitRedis() {
	if config.Conf == nil {
		log.Fatal("Configuration not loaded. Call config.LoadConfig() first.")
	}

	// --- 1. 从配置 (viper) 中获取 Redis 连接信息 ---
	redisAddr := config.Conf.GetString("REDIS_ADDR")
	redisPassword := config.Conf.GetString("REDIS_PASSWORD") // 获取密码
	redisDBStr := config.Conf.GetString("REDIS_DB")          // 获取 DB 编号 (字符串形式)

	// 检查关键配置是否存在
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is not set in the configuration.")
	}

	// --- 2. 处理数据库编号 (DB) ---
	var redisDB int // 用于存储最终的 DB 编号
	var err error
	if redisDBStr == "" {
		// 如果 REDIS_DB 未设置或为空字符串，使用默认值 0
		log.Println("REDIS_DB not set or empty in config, using default DB 0")
		redisDB = 0
	} else {
		// 尝试将从配置中读取的字符串转换为整数
		redisDB, err = strconv.Atoi(redisDBStr)
		if err != nil {
			// 如果转换失败 (例如，配置值不是有效的数字)，记录警告并使用默认值 0
			log.Printf("Warning: Invalid REDIS_DB value '%s' in config. Using default DB 0. Error: %v", redisDBStr, err)
			redisDB = 0 // 安全起见，使用默认值
		}
	}

	// --- 3. 创建 Redis 客户端实例 ---
	// 使用从配置中获取的值来设置选项
	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,     // Redis 服务器地址
		Password: redisPassword, // Redis 密码 (如果为空字符串，则表示没有密码)
		DB:       redisDB,       // Redis 数据库编号
		// 你可以在这里添加其他 redis.Options 配置，例如连接池大小等
		// PoolSize: 10,
	})

	// --- 4. 检查连接和认证 ---
	// 使用 Ping 命令测试连接是否成功以及认证是否通过
	statusCmd := RDB.Ping(Ctx)
	if statusCmd.Err() != nil {
		// 如果 Ping 失败，记录致命错误，程序可能会因此退出
		log.Fatalf("连接 Redis (%s, DB %d) 失败: %v", redisAddr, redisDB, statusCmd.Err())
	} else {
		// 如果 Ping 成功，记录成功的日志信息
		log.Printf("成功连接到 Redis (%s, DB %d)", redisAddr, redisDB)
	}
}

// CloseRedis 可选: 添加一个函数来在程序退出时优雅地关闭 Redis 连接
func CloseRedis() error {
	if RDB != nil {
		log.Println("正在关闭 Redis 连接...")
		return RDB.Close()
	}
	return nil
}
