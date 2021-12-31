package logic

import (
	"github.com/ZYallers/rpcx-framework/helper"
	"github.com/ZYallers/rpcx-service/define"
	"github.com/ZYallers/rpcx-service/model"
	"github.com/ZYallers/rpcx-service/table"
)

type RbacUser struct {
	define.Logic
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
