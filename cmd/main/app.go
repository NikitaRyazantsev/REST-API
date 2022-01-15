package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"project/internal/config"
	"project/internal/user"
	"project/internal/user/db"
	"project/pkg/client/mongodb"
	"project/pkg/logging"
	"project/pkg/shutdown"
	"syscall"
	"time"
)

func main() {
	// Create logger
	logger := logging.GetLogger()

	// Create router
	logger.Info("create router")
	router := httprouter.New()

	// Get data from config
	cfg := config.GetConfig()

	// Connect to MongoDB
	cfgMongo := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username, cfgMongo.Password, cfgMongo.Database, cfgMongo.AuthDB)
	if err != nil {
		logger.Fatal(err)
	}
	userStorage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	// Initialize user service
	userService, err := user.NewService(userStorage, *logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Create handler
	usersHandler := user.Handler{
		Logger:      logger,
		UserService: userService,
	}

	// Routing init
	usersHandler.Register(router)

	// Start application
	logger.Println("start application")
	start(router, cfg, logger)
}

func start(router *httprouter.Router, cfg *config.Config, logger *logging.Logger) {
	logger.Info("start application")

	// Create listener
	var listener net.Listener
	var listenErr error

	// Choose type of connection
	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path: %s", socketPath)

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Info("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	// Star server
	server := &http.Server{
		Handler:      router,
		WriteTimeout: cfg.Timeout.Write * time.Second,
		ReadTimeout:  cfg.Timeout.Read * time.Second,
	}

	// Graceful shutdown
	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	// Start server
	logger.Info("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	logger.Fatalln(server.Serve(listener))
}
