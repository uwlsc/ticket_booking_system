package pkg

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"ticket-api/config"
)

// Database modal
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(logger Logger, env config.Config) Database {
	url := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", env.DBUser, env.DBPassword, env.DBHost,
		env.DBPort, env.DBName,
	)
	if env.DBType != "mysql" {
		url = fmt.Sprintf(
			"%s:%s@unix(/cloudsql/%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			env.DBUser,
			env.DBPassword,
			env.DBHost,
			env.DBName,
		)
	}

	logger.Info("opening db connection")
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{Logger: logger.GetGormLogger()})
	if err != nil {
		logger.Info("Url: ", url)
		logger.Panic(err)
	}

	return Database{DB: db}
}
