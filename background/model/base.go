package model

import (
	"changeme/background/util"

	"gorm.io/gorm"
)

type Base struct {
	ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt int64  `json:"created_at,omitempty" gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `json:"updated_at,omitempty" gorm:"autoUpdateTime:milli"`
}

// 创建时的钩子
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) { b.ID = util.NextID(); return }
