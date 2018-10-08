package sqle

import "fmt"

// UnsafeMysqlCount counts the rows for a single column in a specified table.
//
// This method IS NOT SAFE AGAINST SQL-INJECTION. Use it only with trusted
// input!
//
// As the method's name already clarifies, a Mysql-specific feature is used.
// Therefore don't use this method against other databases.
func (s Sqle) UnsafeMysqlCount(table, column string) (count int, err error) {
	err = s.Select(fmt.Sprintf("SELECT COUNT(%s) FROM %s", column, table), []interface{}{}, []interface{}{&count})
	return
}

// MysqlExists checks whether the statement defined by the `query` and `args`
// would return a result.
//
// As the method's name already clarifies, a Mysql-specific feature is used.
// Therefore don't use this method against other databases.
func (s Sqle) MysqlExists(query string, args ...interface{}) (exists bool, err error) {
	err = s.Select(fmt.Sprintf("SELECT EXISTS (%s)", query), args, []interface{}{&exists})
	return
}
