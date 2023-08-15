package main

import (
	"context"
	"fmt"
	"go-gin-restful-service/config"
	"go-gin-restful-service/log"
	sysrouter "go-gin-restful-service/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	cfg := config.GetConfig()
	sysRouter := sysrouter.Setup(context.Background(), cfg, router)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      sysRouter,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	go func() {
		// 服务连接
		log.Logger.Infof("Server runing at port: %d ", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Errorf("listen: %s\n", err)
		}
	}()
	// // 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Logger.Error("Server Shutdown:", err)
	}
	log.Logger.Info("Server exiting")
}
