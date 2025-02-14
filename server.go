package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aleksannder/url-shortener/common"
	"github.com/aleksannder/url-shortener/handlers"
	"github.com/aleksannder/url-shortener/services"
	"github.com/aleksannder/url-shortener/store"
	"github.com/aleksannder/url-shortener/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	cfg *common.Config
}

func (s *Server) Run() {
	// Init redis
	cacheRepository, err := s.initCacheRepository()
	if err != nil {
		log.Fatal(err)
	}

	// Init consul
	persistentRepository, err := s.initPersistentRepository()
	if err != nil {
		return
	}

	// Init service
	service, err := s.initService(cacheRepository, persistentRepository)
	if err != nil {
		log.Fatal(err)
	}

	// Init handler
	handler, err := s.initHandler(service)
	if err != nil {
		log.Fatal(err)
	}

	// Init syncer and start cron job
	syncer := util.Sync{
		Cache:      cacheRepository,
		Persistent: persistentRepository,
	}
	go syncer.Sync()

	// Start HTTP Server
	s.startHTTPServer(handler)
}

func (s *Server) startHTTPServer(handler *handlers.UrlHandler) {
	// Initialize new router
	router := mux.NewRouter()
	// Set up routing

	router.HandleFunc("/urls/", handler.Insert).Methods("POST")
	router.HandleFunc("/{shortCode}", handler.Redirect).Methods("GET")

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.cfg.ServerAddress, s.cfg.ServerPort),
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go s.listenAndServe(&srv)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
func (s *Server) initCacheRepository() (*store.UrlCacheRepository, error) {
	urlStore, err := store.NewUrlCacheRepository()
	if err != nil {
		return nil, err
	}

	return urlStore, nil
}

func (s *Server) initPersistentRepository() (*store.UrlRepository, error) {
	urlStore, err := store.NewUrlRepository()
	if err != nil {
		return nil, err
	}

	return urlStore, nil
}

func (s *Server) initService(cacheStore *store.UrlCacheRepository, permStore *store.UrlRepository) (*services.UrlService, error) {
	service, err := services.NewUrlService(cacheStore, permStore)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *Server) initHandler(service *services.UrlService) (*handlers.UrlHandler, error) {
	handler, err := handlers.NewUrlHandler(service)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func (s *Server) listenAndServe(srv *http.Server) {
	log.Printf("Starting server on :%s", s.cfg.ServerPort)

	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}
}
