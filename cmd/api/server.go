package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	addr := fmt.Sprintf("0.0.0.0:%d", app.conf.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// shutdown channel
	// 接收 Shutdown() 返回的错误，实现优雅关机
	shutdownError := make(chan error)

	// 捕获关机信号
	go func() {
		// 接收退出信号channel
		quit := make(chan os.Signal, 1)

		// 监听系统的SIGINT和SIGTERM信号，如果有写入quit channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// 阻塞，等待系统信号
		s := <-quit

		slog.Info("接收到关闭服务器信号", slog.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		slog.Info("等到服务后台完成任务")
		app.wg.Wait()

		// 写入关机信号，这里并没有错误，用户主动停止服务
		shutdownError <- nil
	}()

	slog.Info("Starting server on "+addr, slog.String("mode", app.conf.Mode))

	// 调用 Shutdown() 会返回 ErrServerClosed，因此需要排除 http.ErrServerClosed
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// 阻塞，等待关机信号
	err = <-shutdownError
	if err != nil {
		return err
	}

	slog.Info("stopped server", slog.Int("port", app.conf.Port))

	return nil
}
