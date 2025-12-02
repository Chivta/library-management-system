package container

import (
	"lab1/cache"
	"lab1/config"
	"lab1/models"
	"lab1/repository"
	"lab1/validation"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Container struct {
	DB               *gorm.DB
	Config           *config.Config
	Cache            *cache.Cache
	BookRepository   repository.BookRepository
	ReaderRepository repository.ReaderRepository
	Validator        *validation.Validator
}

func NewContainer(dbPath string, configPath string) (*Container, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("Failed to load config from %s, using default config", configPath)
		cfg = config.DefaultConfig()
	}

	cacheInstance := cache.NewCache(cfg.CacheTTLSeconds)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Book{}, &models.Reader{})
	if err != nil {
		return nil, err
	}

	bookRepo := repository.NewBookRepository(db, cacheInstance)
	readerRepo := repository.NewReaderRepository(db, cacheInstance)

	validator := validation.NewValidator()

	return &Container{
		DB:               db,
		Config:           cfg,
		Cache:            cacheInstance,
		BookRepository:   bookRepo,
		ReaderRepository: readerRepo,
		Validator:        validator,
	}, nil
}

func (c *Container) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
