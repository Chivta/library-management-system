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
	"gorm.io/gorm/logger"
)

type Container struct {
	DB               *gorm.DB
	Config           *config.Config
	Cache            *cache.Cache
	BookRepository   repository.BookRepository
	ReaderRepository repository.ReaderRepository
	UserRepository   repository.UserRepository
	Validator        *validation.Validator
}

func NewContainer(dbPath string, configPath string) (*Container, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("Failed to load config from %s, using default config", configPath)
		cfg = config.DefaultConfig()
	}

	cacheInstance := cache.NewCache(cfg.CacheTTLSeconds)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Suppress "record not found" logs
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Book{}, &models.Reader{}, &models.User{})
	if err != nil {
		return nil, err
	}

	// Seed admin user if not exists
	seedAdminUser(db)

	bookRepo := repository.NewBookRepository(db, cacheInstance)
	readerRepo := repository.NewReaderRepository(db, cacheInstance)
	userRepo := repository.NewUserRepository(db)

	validator := validation.NewValidator()

	return &Container{
		DB:               db,
		Config:           cfg,
		Cache:            cacheInstance,
		BookRepository:   bookRepo,
		ReaderRepository: readerRepo,
		UserRepository:   userRepo,
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

// seedAdminUser creates default admin user if it doesn't exist
func seedAdminUser(db *gorm.DB) {
	var adminUser models.User
	result := db.Where("username = ?", "admin").First(&adminUser)

	if result.Error == gorm.ErrRecordNotFound {
		// Create admin user
		adminUser = models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Role:     "admin",
		}
		if err := adminUser.HashPassword("password"); err != nil {
			log.Printf("Failed to hash admin password: %v", err)
			return
		}
		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
			return
		}
		log.Println("✓ Default admin user created (username: admin, password: password)")
	}
}
