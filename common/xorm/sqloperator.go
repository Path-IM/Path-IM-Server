package xorm

import (
	xormerr "github.com/showurl/Zero-IM-Server/common/xorm/err"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Transaction(tx gorm.DB, fs ...func(tx *gorm.DB) error) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		var err error
		for _, f := range fs {
			if f != nil {
				err = f(tx)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func DetailByWhere(tx *gorm.DB, model interface{}, wheres ...GormWhere) error {
	tableName := model.(schema.Tabler).TableName()
	tx = tx.Table(tableName)
	for _, where := range wheres {
		tx = tx.Where(where.Where, where.args...)
	}
	err := tx.First(model).Error
	if err != nil {
		// 表不存在
		if xormerr.TableNotFound(err) {
			// 创建表
			_ = tx.Table(tableName).AutoMigrate(model)
		}
		return err
	}
	return nil
}

// Count 获取数量
func Count(tx *gorm.DB, model interface{}, where string, args ...interface{}) (int64, error) {
	var total int64
	err := tx.Model(model).Where(where, args...).Count(&total).Error
	if xormerr.TableNotFound(err) {
		_ = tx.AutoMigrate(model)
		err = tx.Model(model).Where(where, args...).Count(&total).Error
	}
	return total, err
}

func Insert(tx *gorm.DB, model interface{}) error {
	tableName := model.(schema.Tabler).TableName()
	err := tx.Table(tableName).Create(model).Error
	if err != nil {
		// 表不存在
		if xormerr.TableNotFound(err) {
			// 创建表
			err = tx.Table(tableName).AutoMigrate(model)
			if err != nil {
				return err
			} else {
				// 创建记录
				return tx.Table(tableName).Create(model).Error
			}
		} else {
			return err
		}
	}
	return nil
}

func ListWithPaging(
	tx *gorm.DB,
	models interface{},
	model interface{},
	no int, size int,
	where string, args ...interface{}) (int64, error) {
	tableName := model.(schema.Tabler).TableName()
	var count int64
	db := tx.Table(tableName).Where(where, args...)
	db.Count(&count)
	return count, Paging(db, no, size).Find(models).Error
}

func Paging(tx *gorm.DB, no int, size int) *gorm.DB {
	return tx.Offset((no - 1) * size).Limit(size)
}
