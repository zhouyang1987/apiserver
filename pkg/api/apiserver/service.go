package apiserver

func QueryServices(serviceName string, pageCnt, pageNum int, appId uint) (list []*Service, total int) {

	if appId == 0 {
		if serviceName != "" {
			db.Where("name like ? ", `%`+serviceName+`%`).Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list)
			db.Model(new(Service)).Where("name like ?", serviceName).Count(&total)
		} else {
			db.Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list)
			db.Model(new(Service)).Count(&total)
		}
	} else {
		if serviceName != "" {
			db.Where("name like ? and app_id=?", `%`+serviceName+`%`, appId).Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list)
			db.Model(new(Service)).Where("name like ? and app_id=?", `%`+serviceName+`%`, appId).Count(&total)
		} else {
			db.Where("app_id=?", appId).Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list)
			db.Model(new(Service)).Where("app_id=?", appId).Count(&total)
		}
	}

	for _, svc := range list {
		if svc.ID != 0 {
			var (
				containerList []*Container
				config        = &ContainerConfig{}
				serviceConfig = &ServiceConfig{}
				base          = &BaseConfig{}
				configmap     = &ConfigMap{}
				superConfig   = &SuperConfig{}
				volumes       []*Volume
				envs          []*Env
				ports         []*Port
			)
			db.Find(&containerList, Container{ServiceId: svc.ID})
			svc.Items = containerList
			for _, container := range containerList {
				db.Find(config, ContainerConfig{ContainerId: container.ID})
				container.Config = config

				db.First(base, BaseConfig{ServiceConfigId: config.ID})
				db.Find(&volumes, Volume{BaseConfigId: base.ID})
				base.Volumes = volumes
				config.BaseConfig = base

				db.First(configmap, ConfigMap{ServiceConfigId: config.ID})
				config.ConfigMap = configmap

				db.First(superConfig, SuperConfig{ServiceConfigId: config.ID})
				db.Find(&envs, Env{SuperConfigId: superConfig.ID})
				db.Find(&ports, Port{SuperConfigId: superConfig.ID})
				superConfig.Envs = envs
				superConfig.Ports = ports
				config.SuperConfig = superConfig
			}
			db.First(serviceConfig, ServiceConfig{ServiceId: svc.ID})
			serviceConfig.BaseConfig = base
			serviceConfig.ConfigMap = configmap
			serviceConfig.SuperConfig = superConfig
			svc.Config = serviceConfig
		}
	}
	return
}

func InsertService(svc *Service) {
	if db.Model(svc).Where("name=?", svc.Name).First(svc).RecordNotFound() {
		db.Model(svc).Save(svc)
	}
}

func DeleteService(svc *Service) {
	db.Delete(svc, svc.ID)
	svcCfg := svc.Config
	db.Delete(svcCfg, "service_id=?", svc.ID)

	svcCfgBase := svcCfg.BaseConfig
	db.Delete(svcCfgBase, "service_config_id=?", svcCfg.ID)

	for _, volume := range svcCfgBase.Volumes {
		db.Delete(volume, "base_config_id=?", svcCfgBase.ID)
	}

	svcCfgMap := svcCfg.ConfigMap
	db.Delete(svcCfgMap, "service_config_id=?", svcCfg.ID)

	svcSuper := svcCfg.SuperConfig
	db.Delete(svcSuper, "service_config_id=?", svcCfg.ID)

	for _, env := range svcSuper.Envs {
		db.Delete(env, "super_config_id=?", svcSuper.ID)
	}
	for _, port := range svcSuper.Ports {
		db.Delete(port, "super_config_id=?", svcSuper.ID)
	}

	for _, c := range svc.Items {
		db.Delete(c, "container_id=?", svc.ID)
	}
}

func UpdateService(svc *Service) {
	db.Model(svc).Update(svc)
}

func QueryServiceById(id uint) *Service {
	svc := &Service{}
	var (
		containerList []*Container
		config        = &ContainerConfig{}
		base          = &BaseConfig{}
		configmap     = &ConfigMap{}
		superConfig   = &SuperConfig{}
		volumes       []*Volume
		envs          []*Env
		ports         []*Port
	)
	if db.First(svc, id).Error == nil && svc.ID != 0 {
		db.Find(&containerList, Container{ServiceId: svc.ID})
		svc.Items = containerList
		for _, container := range containerList {
			db.Find(config, ContainerConfig{ContainerId: container.ID})
			container.Config = config

			db.First(base, BaseConfig{ServiceConfigId: config.ID})
			db.Find(&volumes, Volume{BaseConfigId: base.ID})
			base.Volumes = volumes
			config.BaseConfig = base

			db.First(configmap, ConfigMap{ServiceConfigId: config.ID})
			config.ConfigMap = configmap

			db.First(superConfig, SuperConfig{ServiceConfigId: config.ID})
			db.Find(&envs, Env{SuperConfigId: superConfig.ID})
			db.Find(&ports, Port{SuperConfigId: superConfig.ID})
			superConfig.Envs = envs
			superConfig.Ports = ports
			config.SuperConfig = superConfig

			svc.Config = &ServiceConfig{
				BaseConfig:  config.BaseConfig,
				ConfigMap:   config.ConfigMap,
				SuperConfig: config.SuperConfig,
			}
		}
	}

	return svc
}
