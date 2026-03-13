package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var port string

func main() {
	var logger slogLimittedLogger
	logger = configureSlogger(slog.LevelDebug)
	port = "8000"
	if s := os.Getenv("HEALTH_PORT"); s == "" {
		logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("SET HEALTH_PORT to default value: %s", "8000"))
	} else {
		port = s
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthChkHandler)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGQUIT)
	defer stop()
	svr := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 180 * time.Second,
		IdleTimeout:  -1,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}

	svr.ListenAndServe()
	defer svr.Shutdown(context.Background())
	go func() {
		<-ctx.Done()
		logger.LogAttrs(context.Background(), slog.LevelInfo, "Shutting down server")
		stop()
		svr.Shutdown(context.Background())
	}()
}

func createCmd(c context.Context) *exec.Cmd {
	cmd := exec.CommandContext(c, "php-fpm-healthcheck")
	cmd.Env = append(os.Environ(),
		"FCGI_CONNECT_DEFAULT=\"localhost:9001\"",
	)
	return cmd
}

func healthChkHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	cmd := createCmd(c)

	err := cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
