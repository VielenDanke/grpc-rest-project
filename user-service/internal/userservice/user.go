package userservice

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	cp "github.com/vielendanke/grpc-rest-project/user-service/company"
	"github.com/vielendanke/grpc-rest-project/user-service/configs"
	"github.com/vielendanke/grpc-rest-project/user-service/internal/userservice/handler"
	"github.com/vielendanke/grpc-rest-project/user-service/internal/userservice/repository"
	"github.com/vielendanke/grpc-rest-project/user-service/internal/userservice/service"
	u "github.com/vielendanke/grpc-rest-project/user-service/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

func StartServerGRPS(ctx context.Context, cfg *configs.Config) error {
	l, lErr := net.Listen("tcp", cfg.GRPC.Addr)
	if lErr != nil {
		return lErr
	}
	srv := grpc.NewServer()
	ur := repository.NewUserRepository()
	dial, dErr := grpc.Dial("localhost:9091", grpc.WithInsecure())
	if dErr != nil {
		return dErr
	}
	cs := cp.NewCompanyServiceClient(dial)
	ts := service.NewUserService(ur, cs)
	u.RegisterUserServer(srv, handler.NewUserHandler(ts))
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
	if regErr := u.RegisterUserHandlerFromEndpoint(ctx, mux, cfg.GRPC.Addr, opts); regErr != nil {
		return regErr
	}
	log.Printf("Starting HTTP server on: %s\n", cfg.HTTP.Addr)
	return http.ListenAndServe(cfg.HTTP.Addr, wsproxy.WebsocketProxy(mux))
}
