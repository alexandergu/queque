package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexandergu/queque/internal/api"
	"github.com/alexandergu/queque/internal/queue"
)

func main() {
	engine := queue.NewEngine()
	for key, handler := range QueueHandlers {
		engine.RegisterHandler(key, handler)
	}
	engine.Start()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	router := api.NewRouter(engine)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Listen server error")
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Error shutdown server")
	}

	if err := engine.Stop(); err != nil {
		fmt.Println("Error stop engine")
	}

	fmt.Println("Stopped")
}
