package model

import "time"

// Base a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
// It may be embedded into your model or you may build your own model without it
//    type User struct {
//      gorm.Model
//    }
type Base struct {
	ID        uint `gorm:"primarykey"`
	Isvalid   uint `gorm:"default:1;comment:是否删除 "`
	CreatedAt time.Time
	UpdatedAt time.Time
}
