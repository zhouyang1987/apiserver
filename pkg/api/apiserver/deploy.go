package apiserver

func InsertDeploy(deploy *Deploy) error {
	return db.Create(deploy).Error
}

func InsertProjectConfig(config *ProjectConfig) error {
	return db.Create(config).Error
}
