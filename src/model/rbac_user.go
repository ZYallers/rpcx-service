package model

import (
	"src/libraries/core"
	"src/table"
)

type RbacUser struct {
	core.Model
}

func NewRbacUser() *RbacUser {
	m := &RbacUser{}
	m.DB = m.GetEnjoyThin
	m.Table = table.RbacUserTN
	return m
}
