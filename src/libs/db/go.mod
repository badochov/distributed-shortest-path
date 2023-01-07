module github.com/badochov/distributed-shortest-path/src/libs/db

go 1.19

require (
	github.com/badochov/distributed-shortest-path/src/libs/api v0.0.0
	gorm.io/driver/postgres v1.4.5
	gorm.io/gen v0.3.19
	gorm.io/gorm v1.24.3
	gorm.io/plugin/dbresolver v1.3.0
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
	gorm.io/datatypes v1.0.7 // indirect
	gorm.io/driver/mysql v1.4.0 // indirect
	gorm.io/hints v1.1.0 // indirect
)

replace github.com/badochov/distributed-shortest-path/src/libs/api v0.0.0 => ../api
