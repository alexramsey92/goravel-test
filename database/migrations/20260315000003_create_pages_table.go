package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20260315000003CreatePagesTable struct{}

func (r *M20260315000003CreatePagesTable) Signature() string {
	return "20260315000003_create_pages_table"
}

func (r *M20260315000003CreatePagesTable) Up() error {
	if !facades.Schema().HasTable("pages") {
		return facades.Schema().Create("pages", func(table schema.Blueprint) {
			table.ID()
			table.String("title")
			table.String("slug")
			table.Unique("slug")
			table.LongText("content").Nullable()
			table.Boolean("published").Default(false)
			table.Timestamps()
		})
	}
	return nil
}

func (r *M20260315000003CreatePagesTable) Down() error {
	return facades.Schema().DropIfExists("pages")
}
