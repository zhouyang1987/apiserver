package autorization

import (
	"apiserver/pkg/storage/mysqld"
)

var (
	db = mysqld.GetDB()
)

func init() {
	db.SingularTable(true)
	db.CreateTable(
		new(Team),
		new(User),
		new(Permission),
		new(Role),
	)
}
