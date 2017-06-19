package apiserver

func QueryApps(namespace, appName string, pageCnt, pageNum int) (list []*App, total int) {
	if appName != "" {
		db.Where("user_name=? and name like ? ", namespace, `%`+appName+`%`).Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list)
		db.Model(new(App)).Where("user_name=? and name like ? ", namespace, `%`+appName+`%`).Count(&total)
	} else {
		db.Where("user_name=?", namespace).Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list)
		db.Model(new(App)).Where("user_name=?", namespace).Count(&total)
	}
	var (
		svcList       []*Service
		containerList []*Container
	)
	for _, app := range list {
		db.Find(&svcList, Service{AppId: app.ID})
		app.Items = svcList
		for _, svc := range svcList {
			config := &ServiceConfig{}
			var (
				base        = &BaseConfig{}
				configGroup = &ConfigGroup{}
				configmaps  = []*ConfigMap{}
				superConfig = &SuperConfig{}
				volumes     []*Volume
				envs        []*Env
				ports       []*Port
			)

			db.Find(&containerList, Container{ServiceId: svc.ID})
			svc.Items = containerList

			db.Find(config, ServiceConfig{ServiceId: svc.ID})
			db.First(base, BaseConfig{ServiceConfigId: config.ID})
			db.Find(&volumes, Volume{BaseConfigId: base.ID})
			base.Volumes = volumes
			config.BaseConfig = base

			db.First(configGroup, ConfigGroup{ServiceConfigId: config.ID})
			config.ConfigGroup = configGroup

			db.Find(&configmaps, ConfigMap{ConfigGroupId: configGroup.ID})
			configGroup.ConfigMaps = configmaps

			db.First(superConfig, SuperConfig{ServiceConfigId: config.ID})
			db.Find(&envs, Env{SuperConfigId: superConfig.ID})
			db.Find(&ports, Port{SuperConfigId: superConfig.ID})
			superConfig.Envs = envs
			superConfig.Ports = ports
			config.SuperConfig = superConfig

			svc.Config = config

		}
		app.ServiceCount = len(svcList)
	}
	return
}

func InsertApp(app *App) {
	svcConfig := app.Items[0].Config
	if len(app.Items[0].Items) != 0 {
		app.Items[0].Items[0].Config = &ContainerConfig{
			BaseConfig:  svcConfig.BaseConfig,
			SuperConfig: svcConfig.SuperConfig,
		}
	}

	if svcConfig.ConfigGroup != nil {
		for _, c := range svcConfig.ConfigGroup.ConfigMaps {
			UpdateConfigMap(c)
		}
	}

	if svcConfig.ConfigGroup != nil {
		configGroupId := svcConfig.ConfigGroup.ID
		db.Model(new(ConfigGroup)).Set("gorm:save_associations", false).Update(&ConfigGroup{ServiceConfigId: app.Items[0].Config.ID, ServiceName: app.Items[0].Name, ID: configGroupId})
	}
	svcConfig.ConfigGroup = nil
	if db.Model(app).Where("name=?", app.Name).First(app).RecordNotFound() {
		app.Items[0].AppName = app.Name
		db.Model(app).Save(app)
	}

}

func UpdateApp(app *App) {
	for _, svc := range app.Items {
		svc.Status = app.AppStatus
		db.Model(svc).Update(svc)
	}
	db.Model(app).Update(app)
}

func UpdateAppOnly(app *App) {
	db.Model(app).Set("gorm:save_associations", false).Update(app)
}

func DeleteApp(app *App) {
	db.Delete(app)
	for _, svc := range app.Items {
		db.Delete(svc, "app_id=?", app.ID)

		svcCfg := svc.Config
		db.Delete(svcCfg, "service_id=?", svc.ID)

		svcCfgBase := svcCfg.BaseConfig
		db.Delete(svcCfgBase, "service_config_id=?", svcCfg.ID)

		for _, volume := range svcCfgBase.Volumes {
			db.Delete(volume, "base_config_id=?", svcCfgBase.ID)
		}

		svcSuper := svcCfg.SuperConfig
		db.Delete(svcSuper, "service_config_id=?", svcCfg.ID)

		for _, env := range svcSuper.Envs {
			db.Delete(env, "super_config_id=?", svcSuper.ID)
		}
		for _, port := range svcSuper.Ports {
			db.Delete(port, "super_config_id=?", svcSuper.ID)
		}

		for _, c := range svc.Items {
			db.Delete(c)
			db.Delete(c.Config, "container_id=?", c.ID)
		}

		if svc.Config.ConfigGroup != nil {
			for _, c := range svc.Config.ConfigGroup.ConfigMaps {
				c.ContainerPath = ""
			}
			db.Model(new(ConfigGroup)).Update(svc.Config.ConfigGroup)
			db.Exec("update config_group set service_config_id = ?,service_name=?", 0, "")
		}

	}
}

func QueryAppById(id uint) *App {
	app := &App{}
	db.First(app, id)

	var (
		svcList       []*Service
		containerList []*Container
	)
	db.Find(&svcList, Service{AppId: app.ID})
	app.Items = svcList
	for _, svc := range svcList {
		config := &ServiceConfig{}
		var (
			base        = &BaseConfig{}
			configGroup = &ConfigGroup{}
			configmaps  = []*ConfigMap{}
			superConfig = &SuperConfig{}
			volumes     []*Volume
			envs        []*Env
			ports       []*Port
		)

		db.Find(&containerList, Container{ServiceId: svc.ID})
		for _, c := range containerList {
			contaienrConfig := &ContainerConfig{}
			db.First(contaienrConfig, ContainerConfig{ContainerId: c.ID})
			c.Config = contaienrConfig
		}
		svc.Items = containerList

		db.Find(config, ServiceConfig{ServiceId: svc.ID})
		db.First(base, BaseConfig{ServiceConfigId: config.ID})
		db.Find(&volumes, Volume{BaseConfigId: base.ID})
		base.Volumes = volumes
		config.BaseConfig = base

		db.First(configGroup, ConfigGroup{ServiceConfigId: config.ID})
		config.ConfigGroup = configGroup

		db.Find(&configmaps, ConfigMap{ConfigGroupId: configGroup.ID})
		configGroup.ConfigMaps = configmaps

		db.First(superConfig, SuperConfig{ServiceConfigId: config.ID})
		db.Find(&envs, Env{SuperConfigId: superConfig.ID})
		db.Find(&ports, Port{SuperConfigId: superConfig.ID})
		superConfig.Envs = envs
		superConfig.Ports = ports
		config.SuperConfig = superConfig

		svc.Config = config
	}
	app.ServiceCount = len(svcList)
	return app
}

//only query the app table, doesn't query it's chirld table
func GetAppOnly(id uint) *App {
	app := &App{}
	db.First(app, id)
	return app
}

func ExistApp(app *App) bool {
	return !db.First(app).RecordNotFound()
}

func QueryAppsByNamespace(namespace string) []*App {
	apps := []*App{}
	db.Find(&apps, "user_name=?", namespace)
	return apps
}
