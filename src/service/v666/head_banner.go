package v666

import (
	"fmt"
	"net/http"
	"src/config/define"
	"src/libraries/core"
	"src/libraries/util/helper"
	"src/logic"
	"src/model"
	"src/table"
	"time"
)

const (
	bannerModel    = 1
	notBannerModel = 2
	notOnlineState = 2
	onlineState    = 1
	offlineState   = -1
)

type HeadBanner struct {
	core.Service
	tag struct {
		List    func() `path:"head/banner/list"`
		Edit    func() `path:"head/banner/edit"`
		Save    func() `path:"head/banner/save"`
		Delete  func() `path:"head/banner/delete"`
		OnLine  func() `path:"head/banner/online"`
		OffLine func() `path:"head/banner/offline"`
		Rotate  func() `path:"head/banner/rotate"`
	}
}

func (h *HeadBanner) Rotate() error {
	m, err := logic.NewHeadModel().LatestModel()
	if err != nil || m.Id <= 0 {
		return h.Json(http.StatusOK, "未配置显示模式")
	}
	var data []table.EtHeadBanner
	switch m.Model {
	case notBannerModel:
		banner := logic.NewHeadBanner().GetNotBanner()
		if banner.Id > 0 {
			data = append(data, banner)
		}
	case bannerModel:
		data = logic.NewHeadBanner().GetBanner()
	}
	return h.Json(http.StatusOK, "", data)
}

func (h *HeadBanner) List() error {
	modelId := h.GetInt("model")
	if modelId == 0 {
		return h.Json(http.StatusInternalServerError, define.ErrMissReqParam)
	}
	title := h.GetString("title")
	state := h.GetInt("state")
	page := h.GetInt("page", 1)
	limit := h.GetInt("limit", 10)

	query := "`model`=?"
	args := []interface{}{modelId}

	switch state {
	case notOnlineState:
		switch modelId {
		case bannerModel:
			query += " AND `state` IN (1,-1) AND `start_time`>?"
			args = append(args, helper.NowTime())
		}
	case onlineState:
		switch modelId {
		case bannerModel:
			query += " AND `state` IN (1,-1) AND ? BETWEEN `start_time` AND `end_time`"
			args = append(args, helper.NowTime())
		case notBannerModel:
			query += " AND `state`=1"
		}
	case offlineState:
		switch modelId {
		case bannerModel:
			query += " AND `state` IN (1,-1) AND `end_time`<?"
			args = append(args, helper.NowTime())
		case notBannerModel:
			query += " AND `state`=-1"
		}
	default:
		query += " AND `state` IN (1,-1)"
	}

	if title != "" {
		query += " AND `title` LIKE ?"
		args = append(args, "%"+title+"%")
	}

	where := []interface{}{query}
	where = append(where, args...)

	count := model.NewHeadBanner().Count(where)
	if count == 0 {
		return h.Json(http.StatusOK, "", define.M{"list": nil, "page": page, "limit": limit, "count": count})
	}

	fd := "id,model,title,image,url,sort,start_time,end_time,state,admin_user_id,update_time"
	offset := limit * (page - 1)
	if offset < 0 {
		offset = 0
	}
	var rows []table.EtHeadBanner
	model.NewHeadBanner().Find(&rows, where, fd, "sort desc", offset, limit)

	var ids []int
	for _, row := range rows {
		ids = append(ids, row.AdminUserId)
	}
	ids = helper.RemoveDuplicateWithInt(ids)
	nicknameMap := logic.NewRbacUser().AdminUserNicknames(ids)
	for i, row := range rows {
		if value, exist := nicknameMap[row.AdminUserId]; exist {
			rows[i].AdminUserNickname = value
		}
	}
	return h.Json(http.StatusOK, "", define.M{"list": rows, "page": page, "limit": limit, "count": count})
}

func (h *HeadBanner) Edit() error {
	id := h.GetInt("id")
	if id <= 0 {
		return h.Json(http.StatusInternalServerError, define.ErrMissReqParam)
	}
	fields := "id,model,title,image,url,sort,start_time,end_time,state"
	where := []interface{}{"id=?", id}
	var row table.EtHeadBanner
	model.NewHeadBanner().FindOne(&row, where, fields, "")
	return h.Json(http.StatusOK, "", row)
}

func (h *HeadBanner) Save() error {
	title := h.GetString("title")
	image := h.GetString("image")
	url := h.GetString("url")
	sort := h.GetInt("sort")
	startTime := h.GetString("start_time")
	endTime := h.GetString("end_time")
	state := h.GetInt("state", 1)

	modelId := h.GetInt("model")
	adminUserId := h.GetInt("admin_user_id")

	switch modelId {
	case bannerModel:
		if title == "" || image == "" || url == "" || sort <= 0 || startTime == "" || endTime == "" ||
			adminUserId <= 0 || modelId <= 0 {
			return h.Json(http.StatusInternalServerError, define.ErrMissReqParam)
		}
	case notBannerModel:
		if title == "" || image == "" || adminUserId <= 0 || modelId <= 0 {
			return h.Json(http.StatusInternalServerError, define.ErrMissReqParam)
		}
	}

	id := h.GetInt("id")
	start, _ := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)

	if modelId == notBannerModel && state == onlineState {
		if err := h.checkOnLineNotBanner(); err != nil {
			return h.Json(http.StatusBadGateway, err)
		}
	}

	value := table.EtHeadBanner{
		Id:          id,
		Model:       modelId,
		Title:       title,
		Image:       image,
		Url:         url,
		Sort:        sort,
		StartTime:   &start,
		EndTime:     &end,
		State:       state,
		AdminUserId: adminUserId,
	}
	res, err := model.NewHeadBanner().Save(&value, value.Id)
	if err != nil {
		return h.Json(http.StatusInternalServerError, err)
	}
	data := res.(*table.EtHeadBanner)
	if data.Id == 0 {
		return h.Json(http.StatusInternalServerError, define.ErrOperationFailed)
	}

	switch modelId {
	case bannerModel:
		_ = logic.NewHeadBanner().AddBannerCache(data.Id)
	case notBannerModel:
		_ = logic.NewHeadBanner().DeleteNotBannerCache()
	}

	typeStr := "add"
	intro := "新增banner模式"
	if id > 0 {
		typeStr = "update"
		intro = "编辑banner模式"
	}
	h.Record(define.Record{Type: typeStr, TableName: table.EtHeadBannerTN, DataId: data.Id, Intro: intro})
	return h.Json(http.StatusOK, "保存成功", data)
}

func (h *HeadBanner) Delete() error {
	id := h.GetInt("id")
	if id <= 0 {
		return h.Json(http.StatusInternalServerError, define.ErrMissReqParam)
	}
	fd := "id,model"
	where := []interface{}{"id=?", id}
	var banner table.EtHeadBanner
	model.NewHeadBanner().FindOne(&banner, where, fd, "")
	if banner.Id == 0 {
		return fmt.Errorf("query specified data in the head_banner(%d) does not exist", id)
	}

	err := model.NewHeadBanner().Update([]interface{}{"id=?", id}, table.EtHeadBanner{State: -2})
	if err != nil {
		return h.Json(http.StatusInternalServerError, err)
	}

	switch banner.Model {
	case bannerModel:
		_ = logic.NewHeadBanner().DeleteBannerCache(id)
	case notBannerModel:
		_ = logic.NewHeadBanner().DeleteNotBannerCache()
	}

	h.Record(define.Record{Type: "update", TableName: table.EtHeadBannerTN, DataId: id, Intro: "删除配置"})
	return h.Json(http.StatusOK, "删除成功")
}

func (h *HeadBanner) checkOnLineNotBanner() error {
	var banner table.EtHeadBanner
	model.NewHeadBanner().FindOne(&banner, []interface{}{"model=2 AND state=1"}, "id", "")
	if banner.Id > 0 {
		return fmt.Errorf("当前已有配置正在上线中，配置ID为%d，请下线该配置后再进行操作！", banner.Id)
	}
	return nil
}

func (h *HeadBanner) OnLine() error {
	id := h.GetInt("id")
	if id <= 0 {
		return h.Json(http.StatusInternalServerError, define.ErrMissReqParam)
	}

	if err := h.checkOnLineNotBanner(); err != nil {
		return h.Json(http.StatusBadGateway, err)
	}

	err := model.NewHeadBanner().Update([]interface{}{"id=?", id}, table.EtHeadBanner{State: 1})
	if err != nil {
		return h.Json(http.StatusInternalServerError, err)
	}

	_ = logic.NewHeadBanner().DeleteNotBannerCache()

	h.Record(define.Record{Type: "update", TableName: table.EtHeadBannerTN, DataId: id, Intro: "上线配置"})
	return h.Json(http.StatusOK, "上线成功")
}

func (h *HeadBanner) OffLine() error {
	id := h.GetInt("id")
	if id <= 0 {
		return h.Json(http.StatusInternalServerError, define.ErrMissReqParam)
	}
	err := model.NewHeadBanner().Update([]interface{}{"id=?", id}, table.EtHeadBanner{State: -1})
	if err != nil {
		return h.Json(http.StatusInternalServerError, err)
	}

	_ = logic.NewHeadBanner().DeleteNotBannerCache()

	h.Record(define.Record{Type: "update", TableName: table.EtHeadBannerTN, DataId: id, Intro: "下线配置"})
	return h.Json(http.StatusOK, "下线成功")
}
