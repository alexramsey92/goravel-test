package models

import "github.com/goravel/framework/database/orm"

// Role constants mirroring smbgen's role system
const (
	RoleAdmin  = "admin"
	RoleUser   = "user"
	RoleClient = "client"
)

type User struct {
	orm.Model
	Name     string `gorm:"column:name;not null"`
	Email    string `gorm:"column:email;uniqueIndex;not null"`
	Password string `gorm:"column:password;not null"`
	Role     string `gorm:"column:role;default:user"`
}

// GetID satisfies the auth.Authenticatable interface
func (u *User) GetID() any {
	return u.ID
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsClient() bool {
	return u.Role == RoleClient
}
