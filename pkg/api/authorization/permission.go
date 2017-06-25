package autorization

func InsertPermission(permission *Permission) error {
	return db.Create(permission).Error
}

func QueryPermissions(name string, pageCnt, pageNum int) (list []*Permission, total int, err error) {
	if name != "" {
		err = db.Where("name like ? ", `%`+name+`%`).Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list).Error
		db.Model(new(Permission)).Where("name like ?", name).Count(&total)
	} else {
		err = db.Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list).Error
		db.Model(new(Permission)).Count(&total)
	}
	return
}

func DeletePermission(id uint) error {
	return db.Model(new(Permission)).Delete(new(Permission), id).Error
}
