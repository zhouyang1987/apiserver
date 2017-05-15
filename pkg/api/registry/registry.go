package registry

import (
	"errors"

	"apiserver/pkg/storage/mysqld"
	"apiserver/pkg/util/jsonx"
	"apiserver/pkg/util/log"
)

var (
	engine = mysqld.GetEngine()
)

type Image struct {
	Name string      `json:"name"`
	Fest []*Manifest `json:"manifest"`
}

type Manifest struct {
	Name          string `json:"name" xorm:"pk not null varchar(255)"`
	Tag           string `json:"tag" xorm:"pk not null varchar(255)"`
	Architecture  string `json:"architecture" xorm:"varchar(255)"`
	Os            string `json:"os" xorm:"varchar(255)"`
	Author        string `json:"author" xorm:"varchar(255)"`
	Id            string `json:"id" xorm:"varchar(255)"`
	ParentId      string `json:"parent" xorm:"varchar(255)"`
	Created       string `json:"created" xorm:"varchar(255)"`
	DockerVersion string `json:"docker_version" xorm:"varchar(255)"`
	Pull          string `json:"pull" xorm:"varchar(255)"`
}

func init() {
	engine.ShowSQL(true)
	if err := engine.Sync(new(Manifest)); err != nil {
		log.Fatalf("Sync fail :%s", err.Error())
	}
}

func (manifest *Manifest) String() string {
	manifestStr := jsonx.ToJson(manifest)
	return manifestStr
}

func (manifest *Manifest) Insert() error {
	_, err := engine.Insert(manifest)
	if err != nil {
		log.Debugf("插入数据库失败：%v", err)
		return err
	}
	return nil
}

func (manifest *Manifest) Delete() error {
	_, err := engine.Delete(manifest)
	if err != nil {
		return err
	}
	return nil
}

func (manifest *Manifest) Update() error {
	_, err := engine.Update(manifest)
	if err != nil {
		return err
	}
	return nil
}

func (manifest *Manifest) QueryOne() (*Manifest, error) {
	has, err := engine.Get(manifest)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("current app not exsit")
	}
	return manifest, nil
}

func (manifest *Manifest) QuerySet(where map[string]interface{}) (fests []*Manifest, total int64, err error) {
	log.Infof("%#v", where)
	pageCnt := where["pageCnt"].(int)
	pageNum := where["pageNum"].(int)
	if where["name"].(string) != "" {
		name := where["name"].(string)
		if err = engine.Where("name=?", name).Limit(pageCnt, pageCnt*pageNum).Desc("name").Find(&fests); err != nil {
			return
		}
		/*if total, err = engine.Where("name=?", name).Count(manifest); err != nil {
			return
		}*/
	} else {
		if err = engine.Limit(pageCnt, pageCnt*pageNum).Desc("name").Find(&fests); err != nil {
			return
		}
		/*if total, err = engine.Count(manifest); err != nil {
			return
		}*/
	}
	return
}

func (manifest *Manifest) Exsit() (bool, error) {
	has, err := engine.Get(manifest)
	if err != nil {
		return false, err
	}
	return has, nil
}
