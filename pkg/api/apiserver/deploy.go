package apiserver

func InsertDeploy(deploy *Deploy) error {
	return db.Create(deploy).Error
}
