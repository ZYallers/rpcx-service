package model

import (
	"src/libraries/core"
	"src/table"
)

type HeadModel struct {
	core.Model
}

func NewHeadModel() *HeadModel {
	m := &HeadModel{}
	m.DB = m.GetEnjoyThin
	m.Table = table.EtHeadModelTN
	return m
}
