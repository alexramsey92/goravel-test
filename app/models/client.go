package models

import "github.com/goravel/framework/database/orm"

// ClientStatus mirrors smbgen's client provisioning workflow
const (
	ClientStatusLead      = "lead"
	ClientStatusActive    = "active"
	ClientStatusInactive  = "inactive"
)

type Client struct {
	orm.Model
	Name    string `gorm:"column:name;not null"`
	Email   string `gorm:"column:email;not null"`
	Phone   string `gorm:"column:phone"`
	Company string `gorm:"column:company"`
	Status  string `gorm:"column:status;default:lead"`
	Notes   string `gorm:"column:notes;type:text"`
	UserID  uint   `gorm:"column:user_id"` // assigned staff member

	// Relationships
	User User `gorm:"foreignKey:UserID"`
}
