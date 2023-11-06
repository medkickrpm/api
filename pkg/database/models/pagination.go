package models

import "gorm.io/gorm"

type PageReq struct {
	// Page represents the requested page
	Page int `json:"page" query:"page"`
	// Size represents the number of items in a page
	Size int `json:"size" query:"size"`
}

// Paginate returns scope for pagination
func (p *PageReq) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((p.Page - 1) * p.Size).Limit(p.Size)
	}
}

// NewPageReq returns a new PageReq with default values
func NewPageReq() PageReq {
	return PageReq{1, 10}
}
