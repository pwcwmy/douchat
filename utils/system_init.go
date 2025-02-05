package utils

import (
	// "douchat/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"context"
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config") // 路径一定要写对
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config mysql:", viper.Get("mysql"))
}

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func InitMySQL() {
	// 自定义日志模版 打印SQL语句
	// logger.New(writer, config)
	newLogger := logger.New(
		// 创建一个将会把日志输出到标准输出，
		// 每条日志消息前面会有一个回车换行符以及标准的时间戳格式（包括日期和时间）的 Logger 对象
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQl阈值
			LogLevel:      logger.Info, // 日志对应级别
			Colorful:      true,        // 日志为彩色
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	// user := models.UserBasic{}
	// DB.Find(&user)
	// fmt.Println(user)
}

var ctx = context.Background()
func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.MinIdleConns"),
	})
	pong, err := Redis.Ping(ctx).Result()
	if err != nil {
		fmt.Println("init Redis err", err)
	} else {
		fmt.Println("init Redis finished", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish 将消息发布到redis
func Publish(ctx context.Context, channel string, msg string) error {
	err := Redis.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Redis.Subscribe(ctx, channel)
	fmt.Println("Sub消息ctx", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Sub消息msg", msg.Payload)
	return msg.Payload, err // Payload 转为string
}