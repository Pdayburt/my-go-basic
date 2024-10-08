package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"net/http"
)

func main() {

	InitViperV1()
	InitLogger()
	InitPrometheus()
	app := InitWebServerByWire()

	for _, c := range app.consumer {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}

	cron := app.cron
	cron.Start()
	defer func() {
		stop := cron.Stop()
		<-stop.Done()

	}()

	ginServer := app.server
	ginServer.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello webook")
	})
	ginServer.Run(":8080")

}

func InitPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func InitViperV3ByRemoteETCD() {
	err := viper.AddRemoteProvider("etcd3",
		"127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	keys := viper.AllKeys()
	for _, key := range keys {
		v := viper.GetString(key)
		fmt.Println(key, v)
	}

	if err != nil {
		panic(err)
	}
}

// InitViperV2 不用环境下的配置 测试 生产
func InitViperV2() {
	cfile := pflag.String("config",
		"./config/dev.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(e.Name, e.Op)
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

}

// InitViperV1 使用结构体读取额配置文件
func InitViperV1() {
	viper.SetConfigFile("/Users/anatkh/Downloads/blockChain/golang/my-go-basic/webook/config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func InitViper() {

	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

}

func initWebServer() *gin.Engine {

	/*ginServer := gin.Default()
	ginServer.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "http://localhost") ||
				strings.Contains(origin, "webook.com") {
				return true
			}
			return strings.Contains(origin, "company.com")
		},
		MaxAge: 12 * time.Hour,
	}))*/
	// 生成随机密钥
	/*	authKey, _ := generateRandomKey(32)
		encryptionKey, _ := generateRandomKey(32)

		//store := cookie.NewStore([]byte("secret"))
		store, _ := redis.NewStore(16, "tcp",
			"localhost:6379", "", authKey, encryptionKey)
		ginServer.Use(sessions.Sessions("my_session", store))
		ginServer.Use(middleware.NewLoginMiddlewareBuilder().
			IgnorePath("/users/login").
			IgnorePath("/users/signup").
			Build())*/
	/*	ginServer.Use(middleware.NewLoginJwtMiddlewareBuilder().
		IgnorePath("/users/login").
		IgnorePath("/users/signup").
		IgnorePath("/users/login_sms/code/send").
		IgnorePath("/users/login_sms").
		Build())*/

	/*redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.RedisConf.Addr,
	})*/
	//ginServer.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	//return ginServer
	return nil

}

/*func initUser(db *gorm.DB, redis redis.Cmdable) *web.UserHandler {
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(redis)
	userRepository := repository.NewUserRepository(userDao, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(redis)
	codeRepository := repository.NewCodeRepository(codeCache)
	memService := memory.NewService()
	codeService := service.NewCodeService(codeRepository, memService, "9527")

	userHandler := web.NewUserHandler(userService, codeService)
	return userHandler
}*/

/*func generateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
*/
