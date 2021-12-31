package model

import (
	"github.com/ZYallers/rpcx-service/define"
	"github.com/ZYallers/rpcx-service/table"
)

type HeadModel struct {
	define.Model
}

func NewHeadModel() *HeadModel {
	m := &HeadModel{}
	m.DB = m.EnjoyThin
	m.Table = table.EtHeadModelTN
	return m
}
