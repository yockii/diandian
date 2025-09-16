package model

import (
	"diandian/background/util"

	"gorm.io/gorm"
)

type Base struct {
	ID        uint64 `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt int64  `json:"created_at,omitempty" gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `json:"updated_at,omitempty" gorm:"autoUpdateTime:milli"`
}

// 创建时的钩子
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) { b.ID = util.NextID(); return }
