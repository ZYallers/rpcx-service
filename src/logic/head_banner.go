package logic

import (
	"fmt"
	jsontime "github.com/liamylian/jsontime/v2/v2"
	"sort"
	"src/config/env"
	"src/libraries/core"
	"src/libraries/util/helper"
	"src/model"
	"src/table"
	"strconv"
	"time"
)

type HeadBanner struct {
	core.Redis
}

func NewHeadBanner() *HeadBanner {
	return &HeadBanner{}
}

func (h *HeadBanner) GetNotBanner() table.EtHeadBanner {
	var output table.EtHeadBanner
	_ = h.CacheWithString(env.RedisKey.String.HeadNotBanner, &output, func() (interface{}, bool) {
		var t table.EtHeadBanner
		where := []interface{}{"model=2 AND state=1"}
		model.NewHeadBanner().FindOne(&t, where, "id,model,title,image", "id desc")
		return t, t.Id == 0
	})
	return output
}

func (h *HeadBanner) GetBanner() []table.EtHeadBanner {
	var output []table.EtHeadBanner
	if caches := h.HashGetAll(env.RedisKey.Hash.HeadBanner); len(caches) > 0 {
		for _, value := range caches {
			var str string
			if v, ok := value.(string); !ok || v == "" {
				continue
			} else {
				str = v
			}

			var banner table.EtHeadBanner
			_ = jsontime.ConfigWithCustomTimeFormat.Unmarshal(helper.String2Bytes(str), &banner)
			if banner.Id > 0 && banner.StartTime.Before(time.Now()) && banner.EndTime.After(time.Now()) {
				output = append(output, banner)
			}

			if len(output) > 0 {
				sort.SliceStable(output, func(i, j int) bool {
					if output[i].Sort == output[j].Sort {
						return output[i].Id > output[j].Id
					}
					return output[i].Sort > output[j].Sort
				})
			}
		}
		return output
	}

	where := []interface{}{"model=1 AND state IN(1,-1) AND end_time>=?", helper.NowTime()}
	fd := "id,model,title,image,url,sort,start_time,end_time,state"
	var rows []table.EtHeadBanner
	model.NewHeadBanner().Find(&rows, where, fd, "sort desc", 0, 100)

	data := map[string]interface{}{}
	for _, row := range rows {
		data[strconv.Itoa(row.Id)] = row
		if row.StartTime.Before(time.Now()) && row.EndTime.After(time.Now()) {
			output = append(output, row)
		}
	}

	_, _ = h.DeleteCache(env.RedisKey.Hash.HeadBanner)
	if err := h.HashMultiSet(env.RedisKey.Hash.HeadBanner, data); err == nil {
		expire := h.NoDataExpiration()
		if len(data) > 0 {
			expire = helper.TodayRemainSecond()
		}
		h.GetCache().Expire(env.RedisKey.Hash.HeadBanner, expire)
	}

	return output
}

func (h *HeadBanner) AddBannerCache(id int) error {
	var banner table.EtHeadBanner
	fd := "id,model,title,image,url,sort,start_time,end_time,state"
	model.NewHeadBanner().FindOne(&banner, []interface{}{"id=?", id}, fd, "")
	if banner.Id <= 0 {
		return fmt.Errorf("query specified data in the head_banner(%d) does not exist", id)
	}

	data := map[string]interface{}{strconv.Itoa(id): banner}
	if err := h.HashMultiSet(env.RedisKey.Hash.HeadBanner, data); err != nil {
		return err
	}

	if h.GetCache().TTL(env.RedisKey.Hash.HeadBanner).Val().Seconds() == env.Redis.TTL.Forever {
		h.GetCache().Expire(env.RedisKey.Hash.HeadBanner, helper.TodayRemainSecond())
	}

	return nil
}

func (h *HeadBanner) DeleteBannerCache(id int) error {
	return h.HashMultiDelete(env.RedisKey.Hash.HeadBanner, strconv.Itoa(id))
}

func (h *HeadBanner) DeleteNotBannerCache() error {
	return h.GetCache().Del(env.RedisKey.String.HeadNotBanner).Err()
}
