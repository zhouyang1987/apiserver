package autorization

func InsertRole(role *Role) error {
	return db.Create(role).Error
}

func QueryRoles(name string, pageCnt, pageNum int) (list []*Role, total int, err error) {
	if name != "" {
		err = db.Where("name like ? ", `%`+name+`%`).Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list).Error
		db.Model(new(Role)).Where("name like ?", name).Count(&total)
	} else {
		err = db.Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list).Error
		db.Model(new(Role)).Count(&total)
	}
	return
}

func DeleteRole(id uint) error {
	return db.Model(new(Role)).Delete(new(Role), id).Error
}
