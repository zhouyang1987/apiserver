package apiserver

import (
	"apiserver/pkg/storage/mysqld"
)

var (
	db = mysqld.GetDB()
)

func init() {
	db.SingularTable(true)
	db.CreateTable(
		new(App),
		new(Service),
		new(Container),
		new(Port),
		new(Env),
		new(SuperConfig),
		new(ConfigMap),
		new(Volume),
		new(BaseConfig),
		new(ServiceConfig),
		new(ContainerConfig),
		new(ConfigGroup),
		new(Deploy),
		new(DeployItem),
		new(ProjectConfig),
		new(Result),
		new(ResultItem),
	)
}
