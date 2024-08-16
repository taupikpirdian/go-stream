package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/logger"

	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/pkg/errors"

	pkg "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/config"
	pkgRedis "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/service/redis_client"
	"google.golang.org/grpc"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/pkg/utils"
	configServer "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/config"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/delivery/grpc_server"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/psql_user"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/psql_whitelist"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/service"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/usecase"
)

func main() {
	cfg := configServer.Chat()
	rdb := pkg.RedisConnection(cfg.Redis)
	db, err := configServer.InitDatabase()
	if err != nil {
		panic(errors.Wrap(err, "db: "))
	}
	nr, err := pkg.NewRelicApplication()
	if err != nil {
		panic(errors.Wrap(err, "newrelic: "))
	}
	appLogger := logger.NewApiLogger(nr)
	ctx, _ := context.WithCancel(context.Background())

	// amqpConn, err := rabbitmq.NewRabbitMQConn()
	// if err != nil {
	// 	appLogger.Fatal(err)
	// }
	// println(amqpConn, cancel)

	lis, err := net.Listen("tcp", os.Getenv("CHAT_GENAI_PORT"))
	if err != nil {
		log.Fatalln(err)
	}
	redis := pkgRedis.NewRedisClient(cfg.Redis, rdb)

	// repo initiate
	clientGenAi := utils.NewHttpClient(cfg.Timeout*time.Second, cfg.ProxyUrl)
	repoChatGenAi := service.ChatRepositoryFactory{
		Cfg:         cfg,
		Logger:      appLogger,
		Client:      clientGenAi,
		RedisClient: redis,
	}
	chatGenAiRepo, err := repoChatGenAi.Create()
	if err != nil {
		appLogger.Fatal(err)
	}

	/*
		user repo factory
	*/
	psqlUser := psql_user.UserRepositoryFactory{
		Db: db,
	}
	psqlUserRepo, err := psqlUser.Create()
	if err != nil {
		panic(err)
	}

	/*
		whitelist repo factory
	*/
	psqlWhitelist := psql_whitelist.WhiteListRepositoryFactory{
		Db: db,
	}
	psqlWhitelistRepo, err := psqlWhitelist.Create()
	if err != nil {
		panic(err)
	}

	// usecase initiate
	usecaseChatGenAi := usecase.ChatUCFactory{
		L:             appLogger,
		ServiceGenAi:  chatGenAiRepo,
		RedisClient:   redis,
		RepoUser:      psqlUserRepo,
		Cfg:           cfg,
		RepoWhitelist: psqlWhitelistRepo,
	}
	chatGenAiUC, err := usecaseChatGenAi.Create()
	if err != nil {
		appLogger.Fatal(err)
	}

	// handler initiate
	s := grpc.NewServer(
		grpc.UnaryInterceptor(nrgrpc.UnaryServerInterceptor(nr)),
		grpc.StreamInterceptor(nrgrpc.StreamServerInterceptor(nr)),
	)
	chatServer := grpc_server.Server{
		GRPCServer: s,
		L:          appLogger,
		ChatUC:     chatGenAiUC,
	}
	handlerChat, err := chatServer.CreateHandlerChatService()
	if err != nil {
		appLogger.Fatal(err)
	}
	handlerChat.Handle()

	go func() {
		appLogger.Info("Serving gRPC on " + os.Getenv("CHAT_GENAI_PORT"))
		log.Fatal(s.Serve(lis))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		appLogger.Error(fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		appLogger.Error(fmt.Sprintf("ctx.Done: %v", done))
	}

	s.GracefulStop()
	appLogger.Info("Server Exited Properly")
}
