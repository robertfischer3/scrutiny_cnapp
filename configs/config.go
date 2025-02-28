package configs

import (
    "github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    // Add other configurations as needed
}

// ServerConfig holds all server-related configuration
type ServerConfig struct {
    Port    int
    Timeout int
}

// DatabaseConfig holds all database-related configuration
type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
}

// LoadConfig reads configuration from files or environment variables
func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    
    // Also read from environment variables
    viper.AutomaticEnv()
    
    err = viper.ReadInConfig()
    if err != nil {
        return
    }
    
    err = viper.Unmarshal(&config)
    return
}