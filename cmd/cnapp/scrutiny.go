package main

import (
    "fmt"
    "net/http"
    "path/filepath"
    "runtime"

    "github.com/robertfischer3/scrutiny_cnapp/configs"
    "github.com/robertfischer3/scrutiny_cnapp/internal/app/handler"
    "github.com/robertfischer3/scrutiny_cnapp/internal/pkg/logger"
    
    "github.com/gorilla/mux"
)

func main() {
    // Set up logger
    log := logger.GetLogger()
    
    // Load configuration
    _, b, _, _ := runtime.Caller(0)
    basepath := filepath.Dir(b)
    configPath := filepath.Join(filepath.Dir(filepath.Dir(basepath)), "configs")
    
    config, err := configs.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Set up router
    r := mux.NewRouter()
    
    // Register handlers
    handler.RegisterHandlers(r)
    
    // Set up middleware
    r.Use(handler.LoggingMiddleware)
    
    // Start server
    addr := fmt.Sprintf(":%d", config.Server.Port)
    log.Infof("Starting server on %s", addr)
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}