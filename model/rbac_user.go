package model

import (
	"github.com/ZYallers/rpcx-service/define"
	"github.com/ZYallers/rpcx-service/table"
)

type RbacUser struct {
	define.Model
}

func NewRbacUser() *RbacUser {
	m := &RbacUser{}
	m.DB = m.EnjoyThin
	m.Table = table.RbacUserTN
	return m
}
