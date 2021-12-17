package model

import (
	"src/libraries/core"
	"src/table"
)

type HeadBanner struct {
	core.Model
}

func NewHeadBanner() *HeadBanner {
	m := &HeadBanner{}
	m.DB = m.GetEnjoyThin
	m.Table = table.EtHeadBannerTN
	return m
}
