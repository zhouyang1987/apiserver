package build

import (
	"errors"
	"time"

	"apiserver/pkg/storage/mysqld"
	"apiserver/pkg/util/jsonx"
	"apiserver/pkg/util/log"
)

type Build struct {
	Id         int       `json:"id" xorm:"pk not null autoincr"`
	AppName    string    `json:"app_name" xorm:"not null varchar(255)"`
	Version    string    `json:"version" xorm:"not null varchar(255)"`
	Remark     string    `json:"remark" xorm:"not null varchar(255)"`
	BaseImage  string    `json:"baseImage" xorm:"not null varchar(255)"`
	Image      string    `json:"image" xorm:"not null varchar(255)"`
	Tarball    string    `json:"tarball" xorm:"not null varchar(255)"`
	Registry   string    `json:"registry" xorm:"not null varchar(255)"`
	Repositroy string    `json:"repository" xorm:"not null varchar(255)"`
	Branch     string    `json:"branch" xorm:"not null varchar(255)"`
	Status     int       `json:"status" xorm:"not null int(1)"`
	Create_At  time.Time `json:"create_at" xorm:"created not null"`
}

var (
	engine = mysqld.GetEngine()
)

func init() {
	engine.ShowSQL(true)
	if err := engine.Sync(new(Build)); err != nil {
		log.Fatalf("Sync fail :%s", err.Error())
	}
}

func (this *Build) String() string {
	buildStr := jsonx.ToJson(this)
	return buildStr
}

func (this *Build) Insert() error {
	_, err := engine.Insert(this)
	if err != nil {
		return err
	}
	return nil
}

func (this *Build) Delete() error {
	_, err := engine.Id(this.Id).Delete(this)
	if err != nil {
		return err
	}
	return nil
}

func (this *Build) Update() error {
	_, err := engine.Id(this.Id).Update(this)
	if err != nil {
		return err
	}
	return nil
}

func (this *Build) QueryOne() (*Build, error) {
	has, err := engine.Id(this.Id).Get(this)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("current app not exsit")
	}
	return this, nil
}

func (this *Build) QuerySet() ([]*Build, error) {
	buildSet := []*Build{}
	err := engine.Where("1 and 1 order by create_at desc").Find(&buildSet)
	if err != nil {
		return nil, err
	}
	return buildSet, nil
}
