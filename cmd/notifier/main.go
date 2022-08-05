package main

import (
	"context"
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/repositories/factory"
	"github.com/antonioo83/notifier-server/internal/server"
	grpc3 "github.com/antonioo83/notifier-server/internal/server/grpc"
	mpb "github.com/antonioo83/notifier-server/internal/server/grpc/message_proto"
	upb "github.com/antonioo83/notifier-server/internal/server/grpc/user_proto"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/antonioo83/notifier-server/internal/services/auth"
	factory2 "github.com/antonioo83/notifier-server/internal/services/auth/factory"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	fmt.Printf("Build version:%s\n", buildVersion)
	fmt.Printf("Build date:%s\n", buildDate)
	fmt.Printf("Build commit:%s\n", buildCommit)

	configFromFile, err := config.LoadConfigFile("config.json")
	if err != nil {
		log.Fatalf("i can't load configuration file:" + err.Error())
	}
	cfg, err := config.GetConfigSettings(configFromFile)
	if err != nil {
		log.Fatalf("Can't read config: %s", err.Error())
	}

	var pool *pgxpool.Pool
	context := context.Background()
	pool, err = pgxpool.Connect(context, cfg.DatabaseDsn)
	if err != nil {
		log.Fatalf("Can't connect to the database server: %s", err.Error())
	}
	defer pool.Close()

	userRepository := factory.NewUserRepository(context, pool)
	resourceRepository := factory.NewResourceRepository(context, pool)
	messageRepository := factory.NewMessageRepository(context, pool)
	journalRepository := factory.NewJournalRepository(context, pool)
	userAuthHandler := factory2.NewUserAuthHandler(userRepository, cfg)
	routeParameters :=
		server.RouteParameters{
			Config:             cfg,
			UserRepository:     userRepository,
			ResourceRepository: resourceRepository,
			MessageRepository:  messageRepository,
			JournalRepository:  journalRepository,
		}

	senderService := services.NewMessageSenderService(cfg, messageRepository, journalRepository)
	senderService.Run()

	handler := server.GetRouters(userAuthHandler, routeParameters)
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	if cfg.ServerType == config.HTTPServer {
		var srv = http.Server{Addr: cfg.ServerAddress, Handler: handler}
		runHTTPServer(cfg, &srv)
		shutdownGracefullyHTTPServer(context, &srv, idleConnsClosed, sigint)
		<-idleConnsClosed
		fmt.Println("Server HTTP Shutdown gracefully")
		srv.Shutdown(context)
	} else if cfg.ServerType == config.GRPCServer {
		srv := grpc.NewServer()
		runGRPCServer(cfg, userAuthHandler, srv, routeParameters)
		shutdownGracefullyGRPCServer(srv, idleConnsClosed, sigint)
		<-idleConnsClosed
		fmt.Println("Server GRPC Shutdown gracefully")
		srv.GracefulStop()
	} else {
		log.Fatalf("Unknowned server type")
	}
}

func runHTTPServer(config config.Config, srv *http.Server) {
	if config.EnableHTTPS {
		c := services.NewServerCertificate509Service(1658, "Yandex.Praktikum", "RU")
		if err := c.SaveCertificateAndPrivateKeyToFiles("cert.pem", "private.key"); err != nil {
			log.Fatalf("I can't save certificate and private key to files: %v", err)
		}
		if err := srv.ListenAndServeTLS("cert.pem", "private.key"); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	} else {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}
}

func runGRPCServer(cfg config.Config, uh *auth.UserAuthService, srv *grpc.Server, p server.RouteParameters) {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	var s grpc3.UserServer
	s.Config = p.Config
	s.UserRepository = p.UserRepository
	upb.RegisterUserServer(srv, &s)

	ms := grpc3.MessageServer{
		Config:             cfg,
		UserAuth:           *uh,
		UserRepository:     p.UserRepository,
		ResourceRepository: p.ResourceRepository,
		MessageRepository:  p.MessageRepository,
		JournalRepository:  p.JournalRepository,
	}
	mpb.RegisterMessageServer(srv, &ms)

	fmt.Println("Сервер gRPC начал работу")
	if err := srv.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func shutdownGracefullyHTTPServer(ctx context.Context, srv *http.Server, idleConnsClosed chan struct{}, sigint chan os.Signal) {
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			// ошибки закрытия Listener
			log.Printf("HTTP server Shutdown: %v", err)
		}
		// сообщаем основному потоку, что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()
}

func shutdownGracefullyGRPCServer(srv *grpc.Server, idleConnsClosed chan struct{}, sigint chan os.Signal) {
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		srv.GracefulStop()
		// сообщаем основному потоку, что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()
}
