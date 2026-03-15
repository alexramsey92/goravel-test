package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260315000002CreateClientsTable struct{}

func (r *M20260315000002CreateClientsTable) Signature() string {
	return "20260315000002_create_clients_table"
}

func (r *M20260315000002CreateClientsTable) Up() error {
	if !facades.Schema().HasTable("clients") {
		return facades.Schema().Create("clients", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("email")
			table.String("phone").Nullable()
			table.String("company").Nullable()
			table.String("status").Default("lead") // lead | active | inactive
			table.Text("notes").Nullable()
			table.UnsignedBigInteger("user_id").Nullable()
			table.Foreign("user_id").References("id").On("users").NullOnDelete()
			table.Timestamps()
		})
	}
	return nil
}

func (r *M20260315000002CreateClientsTable) Down() error {
	return facades.Schema().DropIfExists("clients")
}
