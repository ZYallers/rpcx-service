package logic

import (
	"github.com/ZYallers/rpcx-service/define"
	"github.com/ZYallers/rpcx-service/env"
	"github.com/ZYallers/rpcx-service/model"
	"github.com/ZYallers/rpcx-service/table"
)

type HeadModel struct {
	define.Logic
}

func NewHeadModel() *HeadModel {
	h := &HeadModel{}
	h.Client = h.Cache
	return h
}

func (h *HeadModel) LatestModel() (table.EtHeadModel, error) {
	var output table.EtHeadModel
	err := h.CacheWithString(env.RedisKey.String["LatestHeadModel"], &output, env.Redis.CommonExpiration, func() (interface{}, bool) {
		var row table.EtHeadModel
		fd := "id,model,admin_user_id,update_time"
		model.NewHeadModel().FindOne(&row, nil, fd, "id desc")
		return row, row.Id == 0
	})
	return output, err
}
