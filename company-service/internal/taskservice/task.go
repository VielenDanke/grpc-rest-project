package task

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vielendanke/grpc-rest-project/company-service/internal/taskservice/handler"
	"github.com/vielendanke/grpc-rest-project/company-service/internal/taskservice/service"
	"log"
	"net"
	"net/http"

	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"github.com/vielendanke/grpc-rest-project/company-service/configs"
	cp "github.com/vielendanke/grpc-rest-project/company-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartServerGRPS(ctx context.Context, cfg *configs.Config) error {
	l, lErr := net.Listen("tcp", cfg.GRPC.Addr)
	if lErr != nil {
		return lErr
	}
	srv := grpc.NewServer()
	cli := &http.Client{}
	ts := service.NewTaskService(cli)
	cp.RegisterCompanyServiceServer(srv, handler.NewTaskHandler(ts))
	reflection.Register(srv)
	log.Printf("Starting GRPC server on: %s\n", cfg.GRPC.Addr)
	return srv.Serve(l)
}

func StartServerHTTP(ctx context.Context, cfg *configs.Config) error {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(header string) (string, bool) {
			return header, true
		}),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
	)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50000000)),
	}
	if regErr := cp.RegisterCompanyServiceHandlerFromEndpoint(ctx, mux, cfg.GRPC.Addr, opts); regErr != nil {
		return regErr
	}
	log.Printf("Starting HTTP server on: %s\n", cfg.HTTP.Addr)
	return http.ListenAndServe(cfg.HTTP.Addr, wsproxy.WebsocketProxy(mux))
}
