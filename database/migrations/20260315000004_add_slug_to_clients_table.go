package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260315000004AddSlugToClientsTable struct{}

func (r *M20260315000004AddSlugToClientsTable) Signature() string {
	return "20260315000004_add_slug_to_clients_table"
}

func (r *M20260315000004AddSlugToClientsTable) Up() error {
	if facades.Schema().HasTable("clients") && !facades.Schema().HasColumn("clients", "slug") {
		return facades.Schema().Table("clients", func(table schema.Blueprint) {
			table.String("slug").Nullable()
		})
	}
	return nil
}

func (r *M20260315000004AddSlugToClientsTable) Down() error {
	return nil
}
