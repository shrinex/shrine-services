package dtmx

import (
	"fmt"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/workflow"
	"github.com/dtm-labs/dtmdriver"
	_ "github.com/dtm-labs/dtmdriver-gozero" // use go-zero driver
	driver "github.com/dtm-labs/dtmdriver-gozero"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"os"
	"shrine/std/conf/dtm"
	"shrine/std/conf/rdb"
	"shrine/std/globals"
	"strconv"
	"strings"
	"time"
)

func init() {
	_ = dtmdriver.Use(driver.DriverName)
}

func InitHttp(mysqlConf rdb.MySQLConf, redisConf redis.RedisConf, server *rest.Server) {
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/queryPrepared",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			wrapHandler(w, r, func() error {
				bb := MustBarrierFromQuery(r)
				conn, err := getMySql(mysqlConf)
				if err != nil {
					return err
				}

				return bb.QueryPrepared(conn)
			})
		},
	}, {
		Method: http.MethodGet,
		Path:   "/redisQueryPrepared",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			wrapHandler(w, r, func() error {
				bb := MustBarrierFromQuery(r)
				conn, err := getRedis(redisConf)
				if err != nil {
					return err
				}

				return bb.RedisQueryPrepared(conn, 86400)
			})
		},
	}}, rest.WithPrefix("/api/v1/dtm"))
}

func InitGrpcWorkflow(dtmConf dtm.DtmConf, rpcConf zrpc.RpcServerConf, server *grpc.Server) {
	workflow.InitGrpc(dtmConf.GRPCServer, rpcConf.ListenOn, server)
}

func InitHttpWorkflow(restConf rest.RestConf, dtmConf dtm.DtmConf, mysqlConf rdb.MySQLConf, redisConf redis.RedisConf, server *rest.Server) {
	InitHttp(mysqlConf, redisConf, server)

	server.AddRoutes([]rest.Route{{
		Method: http.MethodPost,
		Path:   "/workflow/resume",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			wrapHandler(w, r, func() error {
				data, err := io.ReadAll(r.Body)
				if err != nil {
					return err
				}

				return workflow.ExecuteByQS(r.URL.Query(), data)
			})
		},
	}}, rest.WithPrefix("/api/v1/dtm"))

	workflow.InitHTTP(dtmConf.HTTPServer, fmt.Sprintf("%s/api/v1/workflow/resume", getRestServer(restConf)))
}

func wrapHandler(w http.ResponseWriter, r *http.Request, fn func() error) {
	began := time.Now()
	err := fn()
	code, body := asResponse(err)
	httpx.WriteJsonCtx(r.Context(), w, code, body)

	if code == http.StatusOK || code == http.StatusTooEarly {
		logx.Infof("%2dms %d %s %s", time.Since(began).Milliseconds(), code, r.Method, r.RequestURI)
	} else {
		logx.Errorf("%2dms %d %s %s", time.Since(began).Milliseconds(), code, r.Method, r.RequestURI)
	}
}

func asResponse(err error) (int, *globals.Response) {
	code := http.StatusOK
	if err != nil {
		if errors.Is(err, dtmcli.ErrFailure) {
			code = http.StatusConflict
		} else if errors.Is(err, dtmcli.ErrOngoing) {
			code = http.StatusTooEarly
		} else {
			code = http.StatusInternalServerError
		}
	}

	return code, &globals.Response{
		Code:    int32(code),
		Message: http.StatusText(code),
	}
}

func getRestServer(restConf rest.RestConf) string {
	var sb strings.Builder
	sb.WriteString("http")
	sb.WriteString("://")
	if restConf.Host != "0.0.0.0" {
		sb.WriteString(restConf.Host)
	} else {
		host := os.Getenv("POD_IP")
		if len(host) != 0 {
			sb.WriteString(host)
		} else {
			sb.WriteString(netx.InternalIp())
		}
	}
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(restConf.Port))
	return sb.String()
}
