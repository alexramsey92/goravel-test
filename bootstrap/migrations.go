package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20260315000001CreateUsersTable{},
		&migrations.M20260315000002CreateClientsTable{},
		&migrations.M20260315000003CreatePagesTable{},
		&migrations.M20260315000004AddSlugToClientsTable{},
	}
}
