package fnos

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Enabled        bool
	Socket, Prefix string
}
type Listener struct {
	net.Listener
	TrustedGateway bool
}
type gatewayContextKey struct{}

func IsGatewayRequest(ctx context.Context) bool { return ctx.Value(gatewayContextKey{}) == true }

// GatewayContext 给请求上下文打上网关标记。生产路径由 NewHTTPServer 的
// ConnContext 对整条网关 socket 连接统一设置；导出是为了让鉴权中间件的
// 测试能在单个请求上模拟网关连接，测试之外不要使用。
func GatewayContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, gatewayContextKey{}, true)
}

func Listen(address string, cfg Config) ([]Listener, error) {
	tcpListener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	listeners := []Listener{{Listener: tcpListener}}
	if !cfg.Enabled {
		return listeners, nil
	}
	if cfg.Socket == "" {
		tcpListener.Close()
		return nil, fmt.Errorf("-fnos-app requires -gateway-socket")
	}
	if !strings.HasPrefix(cfg.Prefix, "/app/") {
		tcpListener.Close()
		return nil, fmt.Errorf("-gateway-prefix must begin with /app/")
	}
	if err := os.MkdirAll(filepath.Dir(cfg.Socket), 0755); err != nil {
		tcpListener.Close()
		return nil, err
	}
	if err := os.Remove(cfg.Socket); err != nil && !os.IsNotExist(err) {
		tcpListener.Close()
		return nil, err
	}
	socketListener, err := net.Listen("unix", cfg.Socket)
	if err != nil {
		tcpListener.Close()
		return nil, err
	}
	return append(listeners, Listener{Listener: socketListener, TrustedGateway: true}), nil
}

func NewHTTPServer(handler http.Handler, trusted bool) *http.Server {
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	if trusted {
		server.ConnContext = func(ctx context.Context, _ net.Conn) context.Context {
			return context.WithValue(ctx, gatewayContextKey{}, true)
		}
	}
	return server
}
