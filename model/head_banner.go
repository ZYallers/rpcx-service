package model

import (
	"github.com/ZYallers/rpcx-service/define"
	"github.com/ZYallers/rpcx-service/table"
)

type HeadBanner struct {
	define.Model
}

func NewHeadBanner() *HeadBanner {
	m := &HeadBanner{}
	m.DB = m.EnjoyThin
	m.Table = table.EtHeadBannerTN
	return m
}
