package main

import (
	"be-capstone-project/src/cmd/public/config"
	"be-capstone-project/src/cmd/public/controller"
	"be-capstone-project/src/cmd/public/middleware"
	"be-capstone-project/src/cmd/public/router"
	"be-capstone-project/src/internal/adapter/repository/postgres"
	"be-capstone-project/src/internal/adapter/services"
	"be-capstone-project/src/internal/core/logger"
	"be-capstone-project/src/internal/core/storage"
	"be-capstone-project/src/internal/core/validator"
	webserver_http "be-capstone-project/src/internal/core/web/http"

	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Provide -c flag path to yaml file to load config file ( optional)
var confFile = flag.String("c", "", "Path to the server configuration file.")

func BootstrapAndRun() {
	flag.Parse()
	cfg, err := config.NewAppConfigs(*confFile)
	if err != nil {
		logger.Fatalf("Parse config fail with error: %v", err)
	}
	ctx := context.Background()

	//Helper
	validator.InitCustomValidator(cfg)
	// Repository layer
	//kafkaUserAuthV2Config, err := common_configs.BuildSaramaConfig(cfg.Kafka)
	//
	//kafkaProducerService, err := kafka.NewSyncProducerKafkaService(cfg.Kafka.Brokers, kafkaUserAuthV2Config)
	//if err != nil {
	//	logger.Fatalf("Build kafka service error", err)
	//}
	//traceableHttpClient := client.NewTraceableHttpClient(cfg.App)

	//redisOptions, err := redis.ParseURL(cfg.Redis.Host)
	//if err != nil {
	//	logger.Fatalf("Init redis config fail: %v", err)
	//}
	//redisClient := redis.NewClient(redisOptions)
	//redisRepo, err := redis_repo.NewRedisRepository(redisClient)
	if err != nil {
		logger.Fatalf("Init redis config fail: %v", err)
	}

	postgresClient, err := storage.NewPostgresClient(&cfg.Store)
	if err != nil {
		logger.Fatalf("Unable to init postgres client with err: %v", err)
	}

	// Adapter
	sampleRepository := postgres.NewSampleRepository(postgresClient)

	// Service layer
	sampleService := services.NewSampleService(sampleRepository)

	// Controller layer
	sampleController := controller.NewSampleController(sampleService)

	engine := gin.New()
	//Register middleware and router
	middleware.EnableCoreMiddlewareRequestTracing(engine, *cfg)
	engine.Use(
		middleware.InitContext(),
		gin.CustomRecoveryWithWriter(logger.GetGlobal(), func(c *gin.Context, err any) {
			c.AbortWithStatus(http.StatusInternalServerError)
		}), // replace default panic handler writer by global logger to make a gentle json output of webserver
	)
	router.RegisterGinRouters(engine, &sampleController)

	srv := webserver_http.NewHttpServer(engine, cfg)

	go func() {
		logger.Infof("HTTP Server start at port %v", cfg.App.Port)
		if errStartHttpServer := srv.ListenAndServe(); errStartHttpServer != nil && errStartHttpServer != http.ErrServerClosed {
			logger.Fatalf("HTTP Server start fail on port %v, error: %v", cfg.App.Port, errStartHttpServer)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of configurable seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Graceful shutdown timeout of 0 seconds...")

	ctx, cancel := context.WithTimeout(ctx, 0*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server Shutdown  ErrorCtx:", err)
	}
	// catching ctx.Done(). timeout of x seconds.
	select {
	case <-ctx.Done():
		logger.Info("Application  shutdown.")
	}
	logger.Info("Server exiting")
}
