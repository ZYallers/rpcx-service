package logic

import (
	"src/config/env"
	"src/libraries/core"
	"src/model"
	"src/table"
)

type HeadModel struct {
	core.Redis
}

func NewHeadModel() *HeadModel {
	return &HeadModel{}
}

func (h *HeadModel) LatestModel() (table.EtHeadModel, error) {
	var output table.EtHeadModel
	err := h.CacheWithString(env.RedisKey.String.LatestHeadModel, &output, func() (interface{}, bool) {
		var row table.EtHeadModel
		fd := "id,model,admin_user_id,update_time"
		model.NewHeadModel().FindOne(&row, nil, fd, "id desc")
		return row, row.Id == 0
	})
	return output, err
}
