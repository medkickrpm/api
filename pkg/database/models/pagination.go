package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

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

// SortReq is used to sort data
type SortReq struct {
	// By represents the field to sort by
	By string `json:"sort_by" query:"sort_by"`
	// Direction represents the direction of sorting
	Direction SortDirection `json:"sort_direction" query:"sort_direction"`
}

// SortDirection represents the direction of sorting
type SortDirection string

const (
	// ASC represents ascending order
	ASC SortDirection = "ASC"
	// DESC represents descending order
	DESC SortDirection = "DESC"
)

// Sort returns scope for sorting
func (s *SortReq) Sort() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: s.By}, Desc: strings.ToUpper(string(s.Direction)) == string(DESC)})
	}
}

// NewSortReq returns a new SortReq with default values
func NewSortReq() SortReq {
	return SortReq{"id", DESC}
}
