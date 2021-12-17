package logic

import (
	"src/libraries/core"
	"src/libraries/util/helper"
	"src/model"
	"src/table"
)

type RbacUser struct {
	core.Redis
}

func NewRbacUser() *RbacUser {
	return &RbacUser{}
}

func (r *RbacUser) AdminUserNicknames(ids []int) map[int]string {
	res := map[int]string{}
	ids = helper.RemoveWithInt(ids, 0)
	var rows []table.RbacUser
	model.NewRbacUser().Find(&rows, []interface{}{"id IN (?)", ids}, "id,nickname", "", 0, 0)
	if rows != nil {
		for _, row := range rows {
			res[row.Id] = row.Nickname
		}
	}
	return res
}
