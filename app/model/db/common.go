package db

import (
	"github.com/jinzhu/gorm"
)

// 分页条件
type PageWhereOrder struct {
	Order string
	Where string
	Value []interface{}
}

// GetPage
func GetPage(model, where interface{}, out interface{}, pageIndex, pageSize uint64, totalCount *uint64, whereOrder ...PageWhereOrder) error {
	db := DB.Model(model).Where(where)
	if len(whereOrder) > 0 {
		for _, wo := range whereOrder {
			if wo.Order != "" {
				db = db.Order(wo.Order)
			}
			if wo.Where != "" {
				db = db.Where(wo.Where, wo.Value...)
			}
		}
	}
	err := db.Count(totalCount).Error
	if err != nil {
		return err
	}
	if *totalCount == 0 {
		return nil
	}
	return db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error
}

// First
func First(where interface{}, out interface{}) (notFound bool, err error) {
	err = DB.Where(where).First(out).Error
	if err != nil {
		notFound = gorm.IsRecordNotFoundError(err)
	}
	return
}

// Find
func Find(where interface{}, out interface{}, orders ...string) error {
	db := DB.Where(where)
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}

// Create
func Create(value interface{}) error {
	return DB.Create(value).Error
}

// Save
func Save(value interface{}) error {
	return DB.Save(value).Error
}

// PluckList
func PluckList(model, where interface{}, out interface{}, fieldName string) error {
	return DB.Model(model).Where(where).Pluck(fieldName, out).Error
}
