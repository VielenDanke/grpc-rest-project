package userservice

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
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	cp "github.com/vielendanke/grpc-rest-project/user-service/company"
	"github.com/vielendanke/grpc-rest-project/user-service/configs"
	"github.com/vielendanke/grpc-rest-project/user-service/internal/userservice/handler"
	"github.com/vielendanke/grpc-rest-project/user-service/internal/userservice/repository"
	"github.com/vielendanke/grpc-rest-project/user-service/internal/userservice/service"
	u "github.com/vielendanke/grpc-rest-project/user-service/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func initDB(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pgxCfg, cfgErr := pgxpool.ParseConfig(url)
	if cfgErr != nil {
		return nil, cfgErr
	}
	pool, pErr := pgxpool.ConnectConfig(ctx, pgxCfg)

	if pErr != nil {
		return nil, pErr
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
	db, dbErr := initDB(ctx, cfg.DB.URL)

	if dbErr != nil {
		return dbErr
	}
	ur := repository.NewUserRepository(db)

	var connUrl string

	for _, v := range cfg.Services {
		if v.Name == "company" {
			connUrl = v.ConnUrl
		}
	}
	dial, dErr := grpc.Dial(
		connUrl,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcmiddleware.ChainUnaryClient(
			grpcprometheus.UnaryClientInterceptor,
		)),
		grpc.WithStreamInterceptor(grpcmiddleware.ChainStreamClient(
			grpcprometheus.StreamClientInterceptor,
		)),
	)
	if dErr != nil {
		return dErr
	}
	cs := cp.NewCompanyServiceClient(dial)

	ts := service.NewUserService(ur, cs)

	u.RegisterUserServer(srv, handler.NewUserHandler(ts))

	grpcprometheus.Register(srv)

	reflection.Register(srv)

	log.Printf("Starting GRPC server on: %s\n", cfg.GRPC.Addr)
	return srv.Serve(l)
}

func StartMetricsServer(cfg *configs.Config) error {
	sv := http.NewServeMux()
	sv.Handle("/metrics", promhttp.Handler())
	sv.HandleFunc("/live", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	sv.HandleFunc("/ready", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
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
	if pErr := mux.HandlePath(http.MethodPost, "/v1/files", fHandler); pErr != nil {
		return pErr
	}
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

func fHandler(rw http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if err := r.ParseMultipartForm(1e6); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	f, h, fErr := r.FormFile("file")
	defer closeIO(f)
	if fErr != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	data := r.FormValue("body")
	log.Printf("Body: %s\n", data)
	saveF, sfErr := os.OpenFile(h.Filename, os.O_CREATE, os.ModeDevice)
	if sfErr != nil {
		rw.WriteHeader(http.StatusSeeOther)
		_, wErr := rw.Write([]byte("File already exists"))
		if wErr != nil {
			log.Printf("Error write response %s\n", wErr.Error())
		}
		return
	}
	defer closeIO(saveF)
	if sfErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, cErr := io.Copy(saveF, f); cErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func closeIO(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Printf("Error closing resources: %s\n", err.Error())
	}
}
