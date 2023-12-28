package must

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(driverName, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// if err := db.Ping(); err != nil {
	// 	return nil, err
	// }
	return db, nil
}
