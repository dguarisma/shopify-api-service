package model

import (
	"time"

	"gorm.io/gorm"
)

type ICustom interface {
	GetId() uint
	SetId(uint)
}

func (cm *CustomModel) GetId() uint   { return cm.ID }
func (cm *CustomModel) SetId(id uint) { cm.ID = id }

type CustomModel struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt *time.Time     `json:"-" gorm:"->;<-:create"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
