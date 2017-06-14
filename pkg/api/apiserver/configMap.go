package apiserver

import ()

// func InsertConfigs(c *ConfigMap) {
// 	db.Model(new(ConfigMap)).Create(c)
// }

// func DeleteConfigMap() {

// }

// func QueryConfigs(namespace, appName string, pageCnt, pageNum int) (list []*Config, total int) {
// 	return
// }

func UpdateConfig(c *Config) {
	db.Model(c).Update(c)
}

func DeleteConfig(id uint) {
	db.Model(new(Config)).Delete(new(Config), id)
	db.Model(new(ConfigMap)).Delete(new(ConfigMap), "config_id=? ", id)
}

func DeleteConfigItem(id uint) {
	db.Model(new(ConfigMap)).Delete(new(ConfigMap), id)
}

func InsertConfig(c *Config) {
	db.Model(c).Create(c)
}

func InsertConfigItem(c *ConfigMap) {
	db.Model(c).Create(c)
}

func QueryConfigById(id uint) *Config {
	cfg := &Config{}
	db.Model(new(ConfigMap)).First(cfg, id)

	var cfgmaps []*ConfigMap
	db.Model(new(ConfigMap)).Find(&cfgmaps, "config_id=?", cfg.ID)
	cfg.ConfigMaps = cfgmaps
	return cfg
}

func QueryConfigs(configName string, pageCnt, pageNum int) (list []*Config, total int) {
	if configName != "" {
		db.Where("name like ? ", `%`+configName+`%`).Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list)
		db.Model(new(Config)).Where("name like ?", configName).Count(&total)
	} else {
		db.Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list)
		db.Model(new(Config)).Count(&total)
	}

	for _, cfg := range list {
		var cfgmaps []*ConfigMap
		db.Model(new(ConfigMap)).Find(&cfgmaps, "config_id=?", cfg.ID)
		cfg.ConfigMaps = cfgmaps
	}

	return
}
