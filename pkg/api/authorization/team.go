package autorization

func InsertTeam(tearm *Team) error {
	return db.Create(tearm).Error
}

func QueryTeams(name string, pageCnt, pageNum int) (list []*Team, total int, err error) {
	if name != "" {
		err = db.Where("name like ? ", `%`+name+`%`).Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list).Error
		db.Model(new(Team)).Where("name like ?", name).Count(&total)
	} else {
		err = db.Offset(pageCnt * pageNum).Limit(pageCnt).Order("name desc").Find(&list).Error
		db.Model(new(Team)).Count(&total)
	}
	return
}

func DeleteTeam(id uint) error {
	return db.Model(new(Team)).Delete(new(Team), id).Error
}
