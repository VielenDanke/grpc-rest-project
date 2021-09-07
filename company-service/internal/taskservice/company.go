package task

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vielendanke/grpc-rest-project/company-service/internal/taskservice/handler"
	"github.com/vielendanke/grpc-rest-project/company-service/internal/taskservice/repository"
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

func initDB(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, pErr := pgxpool.ParseConfig(url)

	if pErr != nil {
		return nil, pErr
	}
	pool, cErr := pgxpool.ConnectConfig(ctx, cfg)

	if cErr != nil {
		return nil, cErr
	}
	if pingErr := pool.Ping(ctx); pingErr != nil {
		return nil, pingErr
	}
	return pool, nil
}

func StartServerGRPS(ctx context.Context, cfg *configs.Config) error {
	l, lErr := net.Listen("tcp", cfg.GRPC.Addr)
	if lErr != nil {
		return lErr
	}
	srv := grpc.NewServer()
	pool, connErr := initDB(ctx, cfg.DB.URL)
	if connErr != nil {
		return connErr
	}
	r := repository.NewCompanyRepository(pool)
	ts := service.NewTaskService(r)
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
