package apiserver

import (
	"apiserver/pkg/storage/mysqld"

	"apiserver/pkg/util/jsonx"
	"apiserver/pkg/util/log"
)

var (
	db = mysqld.GetDB()
)

func init() {
	db.SingularTable(true)
	db.CreateTable(&App{}, &Service{}, new(Container), new(Port), new(Env), new(SuperConfig), new(ConfigMap), new(Volume), new(BaseConfig), new(ServiceConfig), new(ContainerConfig))
}

func QueryAll(namespace string) (list []*App) {
	db.Find(&list, App{UserName: namespace})

	var (
		svcList       []*Service
		containerList []*Container
		config        = &ContainerConfig{}
		base          = &BaseConfig{}
		configmap     = &ConfigMap{}
		superConfig   = &SuperConfig{}
		volumes       []*Volume
		envs          []*Env
		ports         []*Port
	)
	for _, app := range list {
		db.Find(&svcList, Service{AppId: app.ID})
		app.Items = svcList
		for _, svc := range svcList {
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
		}
	}
	return
}

func QueryById(id uint) *App {
	app := &App{}
	db.First(app, id)
	var (
		svcList       []*Service
		containerList []*Container
		config        = &ContainerConfig{}
		base          = &BaseConfig{}
		configmap     = &ConfigMap{}
		superConfig   = &SuperConfig{}
		volumes       []*Volume
		envs          []*Env
		ports         []*Port
	)
	db.Find(&svcList, Service{AppId: app.ID})
	app.Items = svcList
	for _, svc := range svcList {
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
				BaseConfig:  base,
				ConfigMap:   configmap,
				SuperConfig: superConfig,
			}
		}
	}
	return app
}

func Insert(app *App) {
	svcConfig := app.Items[0].Config
	app.Items[0].Items[0].Config = &ContainerConfig{
		BaseConfig:  svcConfig.BaseConfig,
		ConfigMap:   svcConfig.ConfigMap,
		SuperConfig: svcConfig.SuperConfig,
	}
	if db.NewRecord(app) {
		db.Create(app)
	}
}

func Update(app *App) {

	for _, svc := range app.Items {
		svc.Status = app.AppStatus
		db.Update(svc)
	}
	app.Items = nil
	db.Update(app)
}

func Delete(app *App) {
	db.Delete(app)
	svc := app.Items[0]
	log.Debug(jsonx.ToJson(svc))
	log.Debug(app.ID)
	db.Delete(svc, "app_id=?", app.ID)

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
		db.Delete(c, "container_ip=?", svc.ID)
	}

}
