package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260315000001CreateUsersTable struct{}

func (r *M20260315000001CreateUsersTable) Signature() string {
	return "20260315000001_create_users_table"
}

func (r *M20260315000001CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		return facades.Schema().Create("users", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("email")
			table.Unique("email")
			table.String("password")
			table.String("role").Default("user") // admin | user | client
			table.Timestamps()
		})
	}
	return nil
}

func (r *M20260315000001CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
