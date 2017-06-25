package autorization

func InsertUser(user *User) error {
	return db.Create(user).Error
}

func QueryUsers(name string, pageCnt, pageNum int) (list []*User, total int, err error) {
	if name != "" {
		err = db.Where("name like ? ", `%`+name+`%`).Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list).Error
		db.Model(new(User)).Where("name like ?", name).Count(&total)
	} else {
		err = db.Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list).Error
		db.Model(new(User)).Count(&total)
	}
	return
}

func DeleteUser(id uint) error {
	return db.Model(new(User)).Delete(new(User), id).Error
}
