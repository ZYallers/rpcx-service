package v666

import (
	"github.com/ZYallers/rpcx-framework/define"
	env2 "github.com/ZYallers/rpcx-framework/env"
	"github.com/ZYallers/rpcx-framework/util/mtsc"
	"github.com/ZYallers/rpcx-service/env"
	"github.com/ZYallers/rpcx-service/logic"
	"github.com/ZYallers/rpcx-service/model"
	"github.com/ZYallers/rpcx-service/table"
	"net/http"
)

type HeadModel struct {
	mtsc.Service
	tag struct {
		LatestModel func() `path:"head/model/latest"`
		SwitchModel func() `path:"head/model/switch"`
	}
}

func (h *HeadModel) LatestModel() error {
	data, err := logic.NewHeadModel().LatestModel()
	if err != nil {
		return h.Json(http.StatusInternalServerError, err)
	}
	if data.AdminUserId > 0 {
		nickname := logic.NewRbacUser().AdminUserNicknames([]int{data.AdminUserId})
		data.AdminUserNickname = nickname[data.AdminUserId]
	}
	return h.Json(http.StatusOK, "", data)
}

func (h *HeadModel) SwitchModel() error {
	modelId := h.GetInt("model")
	adminUserId := h.GetInt("admin_user_id")
	if modelId <= 0 || adminUserId <= 0 {
		return h.Json(http.StatusInternalServerError, env2.ErrMissReqParam)
	}

	data := table.EtHeadModel{Model: modelId, AdminUserId: adminUserId}
	res, err := model.NewHeadModel().Save(&data)
	if err != nil {
		return h.Json(http.StatusInternalServerError, err)
	}

	resId := res.(*table.EtHeadModel).Id
	if resId == 0 {
		return h.Json(http.StatusInternalServerError, env2.ErrOperationFailed)
	}

	_, _ = logic.NewHeadModel().DeleteCache(env.RedisKey.String["LatestHeadModel"])

	h.Record(define.Record{Type: "add", TableName: table.EtHeadModelTN, DataId: resId, Intro: "切换显示模式"})
	return h.Json()
}
