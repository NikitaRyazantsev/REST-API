package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
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

	// Star server
	server := &http.Server{
		Addr:         ":" + cfg.Listen.Port,
		Handler:      router,
		WriteTimeout: cfg.Timeout.Write * time.Second,
		ReadTimeout:  cfg.Timeout.Read * time.Second,
	}

	// Graceful shutdown
	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	// Start server
	logger.Info("server is listening port: %s", cfg.Listen.Port)
	logger.Fatalln(server.ListenAndServe())
}
