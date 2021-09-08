package task

import (
	"context"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcprometheus.StreamServerInterceptor,
			grpcrecovery.StreamServerInterceptor(),
			grpcctxtags.StreamServerInterceptor(),
			grpcopentracing.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcprometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
			grpcctxtags.UnaryServerInterceptor(),
			grpcopentracing.UnaryServerInterceptor(),
		)),
	)
	pool, connErr := initDB(ctx, cfg.DB.URL)

	if connErr != nil {
		return connErr
	}
	r := repository.NewCompanyRepository(pool)

	ts := service.NewTaskService(r)

	cp.RegisterCompanyServiceServer(srv, handler.NewTaskHandler(ts))

	reflection.Register(srv)

	grpcprometheus.Register(srv)

	log.Printf("Starting GRPC server on: %s\n", cfg.GRPC.Addr)
	return srv.Serve(l)
}

func StartMetricsServer(cfg *configs.Config) error {
	sv := http.NewServeMux()
	sv.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting metrics server on: %s\n", cfg.Metrics.Addr)
	return http.ListenAndServe(cfg.Metrics.Addr, sv)
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
