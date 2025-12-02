package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	CacheTTLSeconds    int64  `json:"cache_ttl_seconds"`
	EnableGetBooks     bool   `json:"enable_get_books"`
	EnablePostBooks    bool   `json:"enable_post_books"`
	EnablePutBooks     bool   `json:"enable_put_books"`
	EnableDeleteBooks  bool   `json:"enable_delete_books"`
	EnableGetReaders   bool   `json:"enable_get_readers"`
	EnablePostReaders  bool   `json:"enable_post_readers"`
	EnablePutReaders   bool   `json:"enable_put_readers"`
	EnableDeleteReaders bool  `json:"enable_delete_readers"`
}

func LoadConfig(filePath string) (*Config, error) {
	log.Printf("Loading configuration from: %s", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("Error parsing config file: %v", err)
		return nil, err
	}

	log.Printf("Configuration loaded successfully: cache_ttl_seconds=%d", cfg.CacheTTLSeconds)
	return &cfg, nil
}

func DefaultConfig() *Config {
	return &Config{
		CacheTTLSeconds:    300, // 5 minutes default
		EnableGetBooks:     true,
		EnablePostBooks:    true,
		EnablePutBooks:     true,
		EnableDeleteBooks:  true,
		EnableGetReaders:   true,
		EnablePostReaders:  true,
		EnablePutReaders:   true,
		EnableDeleteReaders: true,
	}
}
