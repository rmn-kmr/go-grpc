package grpc

import (
	"context"
	"fmt"
	nr "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rmnkmr/lsp/internal/app"
	lspPb "github.com/rmnkmr/lsp/proto"
	"github.com/rmnkmr/protonium"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	server *grpc.Server
	port   int
}

func (s *Server) Initialize(nrApp *nr.Application) {
}

func (s *Server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	port := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening grpc on %s", port)

	// add server stop hook
	go func() {
		<-ctx.Done()
		s.server.GracefulStop()
	}()
	if err := s.server.Serve(lis); err != nil {
		return nil
	}
	return nil
}

func New(s *app.App, port int) protonium.Option {
	server := grpc.NewServer()
	lspPb.RegisterAPIServer(server, s)
	grpcServer := &Server{
		server: server,
		port:   port,
	}
	return protonium.Component("grpc-server", grpcServer)
}
