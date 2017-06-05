package registry

import (
	"errors"

	"apiserver/pkg/storage/mysqld"
	"apiserver/pkg/util/jsonx"
	"apiserver/pkg/util/log"
)

var (
	db = mysqld.GetDB()
)

type Image struct {
	Name   string      `json:"name"`
	TagLen int         `json:"tagLen"`
	Fest   []*Manifest `json:"manifest"`
}

type Manifest struct {
	ID            uint
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	Architecture  string `json:"architecture"`
	Os            string `json:"os"`
	Author        string `json:"author"`
	Id            string `json:"id"`
	ParentId      string `json:"parent"`
	Created       string `json:"created"`
	DockerVersion string `json:"docker_version"`
	Pull          string `json:"pull"`
}

func init() {
	db.SingularTable(true)
	db.CreateTable(&Manifest{})
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
	pageCnt := where["pageCnt"].(int)
	pageNum := where["pageNum"].(int)
	if where["name"].(string) != "" {
		name := where["name"].(string)
		if err = engine.Where("name=?", name).Limit(pageCnt, pageCnt*pageNum).Desc("name").Find(&fests); err != nil {
			return
		}
		if total, err = engine.Distinct("name").Where("name=?", name).Count(Manifest{}); err != nil {
			return
		}
	} else {
		if err = engine.Limit(pageCnt, pageCnt*pageNum).Desc("name").Find(&fests); err != nil {
			return
		}
		if total, err = engine.Distinct("name").Count(Manifest{}); err != nil {
			return
		}
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
