package dao

import "database/sql"

// DBMaster and DBSlave
// Wrapped structs for avoiding wire’s error when there are two same types of input parameters.
type DBMaster struct {
	*sql.DB
}
type DBSlave struct {
	*sql.DB
}
