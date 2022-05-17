package global

import (
	envUtils "github.com/showurl/Zero-IM-Server/common/utils/env"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string         `json:"id" gorm:"primarykey;type:bigint(20);comment:'主键'"`
	CreatedAt ModelTime      `json:"created_at" gorm:"column:created_at;index;comment:'创建时间'"`
	UpdatedAt ModelTime      `json:"-" gorm:"column:updated_at;index;comment:'更新时间'"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:'删除时间'"`
}

// DetailExpiredSecond 缓存默认5分钟过期
func (b *BaseModel) DetailExpiredSecond() int {
	if envUtils.IsDev() {
		return 60 * 60 * 24 * 30
	}
	return 60 * 5
}

// RelationExpiredSecond 缓存默认5分钟过期
func (b *BaseModel) RelationExpiredSecond() int {
	if envUtils.IsDev() {
		return 60 * 60 * 24 * 30
	}
	return 60 * 5
}

func (b *BaseModel) GetIdString() string {
	return b.ID
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "0" || b.ID == "" {
		b.ID = GetID()
	}
	return nil
}

func (b *BaseModel) AfterFind(tx *gorm.DB) error {
	return nil
}
