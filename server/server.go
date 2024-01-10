package server

import (
	"BBoltViewer/cmd"
	"BBoltViewer/version"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"BBoltViewer/router"

	"BBoltViewer/g"

	"github.com/gin-gonic/gin"
)

func listenAndServe(srv *http.Server) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

func listenAndServeTLS(srv *http.Server, certFile, keyFile string) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	return srv.ServeTLS(ln, certFile, keyFile)
}

func initHttpServer(port int, handler http.Handler, bIpv4 bool) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	if bIpv4 {
		go func() {
			if err := listenAndServe(srv); err != nil {
				g.Log.Fatal(err)
			}
		}()
		return srv
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			g.Log.Fatal(err)
		}
	}()
	return srv
}

func initHttpsServer(port int, handler http.Handler, bIpv4 bool) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	certFile := "./cert/server.pem"
	keyFile := "./cert/server.key"
	if bIpv4 {
		go func() {
			if err := listenAndServeTLS(srv, certFile, keyFile); err != nil {
				g.Log.Fatal(err)
			}
		}()
		return srv
	}
	go func() {
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			g.Log.Fatal(err)
		}
	}()
	return srv
}

func Run(c *cmd.Command) {
	if c.Loglevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	apiSrv, sProtocol := func() (*http.Server, string) {
		if c.Http {
			return initHttpServer(c.ServerPort, router.Init(c), c.ForceIpv4), "http"
		} else {
			return initHttpsServer(c.ServerPort, router.Init(c), c.ForceIpv4), "https"
		}
	}()
	g.Log.Info(version.Get())
	g.Log.Infof("api server run %s, protocol: %s", apiSrv.Addr, sProtocol)
	signalExit := make(chan os.Signal, 1)
	signal.Notify(signalExit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-signalExit
	close(signalExit)
	close(g.SignalExit)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := apiSrv.Shutdown(ctx); err != nil {
		g.Log.Fatal(err)
	}
	g.Log.Infof("api server exit %s, protocol: %s", apiSrv.Addr, sProtocol)
}
