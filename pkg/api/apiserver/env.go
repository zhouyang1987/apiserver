package apiserver

import (
	"errors"
)

func (e *Env) Insert() error {
	_, err := engine.Insert(e)
	if err != nil {
		return err
	}
	return nil
}

func (e *Env) Delete() error {
	_, err := engine.Delete(e)
	if err != nil {
		return err
	}
	return nil
}

func (e *Env) Update() error {
	_, err := engine.Update(e)
	if err != nil {
		return err
	}
	return nil
}

func (e *Env) QueryOne() (*Env, error) {
	has, err := engine.Get(e)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("current env not exsit")
	}
	return e, nil
}

func (e *Env) QuerySet(where map[string]interface{}) (fests []*Env, total int64, err error) {
	pageCnt := where["pageCnt"].(int)
	pageNum := where["pageNum"].(int)
	if where["name"].(string) != "" {
		name := where["name"].(string)
		if err = engine.Where("name=?", name).Limit(pageCnt, pageCnt*pageNum).Desc("name").Find(&fests); err != nil {
			return
		}
		if total, err = engine.Distinct("name").Where("name=?", name).Count(Env{}); err != nil {
			return
		}
	} else {
		if err = engine.Limit(pageCnt, pageCnt*pageNum).Desc("name").Find(&fests); err != nil {
			return
		}
		if total, err = engine.Distinct("name").Count(Env{}); err != nil {
			return
		}
	}
	return
}

func (e *Env) Exsit() (bool, error) {
	has, err := engine.Get(e)
	if err != nil {
		return false, err
	}
	return has, nil
}
