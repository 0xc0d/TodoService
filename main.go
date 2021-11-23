package main

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"todoService/internal/config"
	"todoService/internal/pb"
	"todoService/internal/repository"
	"todoService/internal/service"
	"todoService/internal/trace"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

func main() {
	cfg := config.Load()

	if cfg.Environment == EnvProduction {
		// MemProfileRate controls the fraction of memory allocations that are recorded
		// and reported in the memory profile. The profiler aims to sample an average of
		// one allocation per MemProfileRate bytes allocated.
		//
		// Setting this variable to zero will turn off memory profiling entirely and
		// eventually saves ~1.406,24 KBs of memory in 64bit machines. due to 179999
		// allocation of memory profile.
		runtime.MemProfileRate = 0
	}

	serverAddress := net.JoinHostPort(cfg.Host, cfg.Port)
	sock, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal().Msgf("failed to listen: %+v", err)
	}
	log.Info().Msgf("server started on %s", serverAddress)

	grpcServer := grpc.NewServer()
	todoServer := service.New(repository.New())
	pb.RegisterTodoServer(grpcServer, todoServer)

	if cfg.Environment == EnvDevelopment {
		reflection.Register(grpcServer)
	}

	if cfg.JaegerEndpoint != "" {
		log.Info().Msgf("starting Jaeger on %s", cfg.JaegerEndpoint)
		trace.SetupJaeger(cfg.JaegerEndpoint, "TodoService", cfg.Environment)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sigChan:
			log.Info().Msg("gracefully stopping server...")
			grpcServer.GracefulStop()
			log.Info().Msg("done")
		}
	}()

	if err = grpcServer.Serve(sock); err != nil {
		log.Fatal().Msgf("failed to serve: %+v", err)
	}
}
