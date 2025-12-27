package Dao

import "gorm.io/gorm"

// InitTables 初始化,自动同步字段
func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
	)
}
