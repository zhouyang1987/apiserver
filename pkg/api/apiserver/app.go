package apiserver

import (
	"errors"
)

import (
	"apiserver/pkg/storage/mysqld"
)

var (
	engine = mysqld.GetEngine()
)

func (e *App) Insert() error {
	_, err := engine.Insert(e)
	if err != nil {
		return err
	}
	return nil
}

func (e *App) Delete() error {
	_, err := engine.Delete(e)
	if err != nil {
		return err
	}
	return nil
}

func (e *App) Update() error {
	_, err := engine.Update(e)
	if err != nil {
		return err
	}
	return nil
}

func (e *App) QueryOne() (*App, error) {
	has, err := engine.Get(e)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("current env not exsit")
	}
	return e, nil
}

func (e *App) QuerySet(where map[string]interface{}) (fests []*App, total int64, err error) {
	pageCnt := where["pageCnt"].(int)
	pageNum := where["pageNum"].(int)
	if where["name"].(string) != "" {
		name := where["name"].(string)
		if err = engine.Where("name=?", name).Limit(pageCnt, pageCnt*pageNum).Desc("name").Find(&fests); err != nil {
			return
		}
		if total, err = engine.Distinct("name").Where("name=?", name).Count(App{}); err != nil {
			return
		}
	} else {
		if err = engine.Limit(pageCnt, pageCnt*pageNum).Desc("name").Find(&fests); err != nil {
			return
		}
		if total, err = engine.Distinct("name").Count(App{}); err != nil {
			return
		}
	}
	return
}

func (e *App) Exsit() (bool, error) {
	has, err := engine.Get(e)
	if err != nil {
		return false, err
	}
	return has, nil
}
