package models

import "github.com/goravel/framework/database/orm"

// Page is a CMS page — mirrors smbgen's CmsPage model
type Page struct {
	orm.Model
	Title     string `gorm:"column:title;not null"`
	Slug      string `gorm:"column:slug;uniqueIndex;not null"`
	Content   string `gorm:"column:content;type:text"`
	Published bool   `gorm:"column:published;default:false"`
}
