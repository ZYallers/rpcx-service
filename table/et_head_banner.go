package table

import (
	"time"
)

const EtHeadBannerTN = `et_head_banner`

type EtHeadBanner struct {
	Id                int        `json:"id,omitempty" gorm:"column:id;type:int(11);unsigned;not null;primaryKey;autoIncrement"`             // 主键ID
	Model             int        `json:"model,omitempty" gorm:"column:model;type:tinyint(1);unsigned;not null;default:1"`                   // 模式|1banner模式|2无banner模式
	Title             string     `json:"title,omitempty" gorm:"column:title;type:varchar(64);not null;default:''"`                          // 名称
	Image             string     `json:"image,omitempty" gorm:"column:image;type:varchar(128);not null;default:''"`                         // 广告位配图
	Url               string     `json:"url,omitempty" gorm:"column:url;type:varchar(128);not null;default:''"`                             // 广告位跳转地址
	Sort              int        `json:"sort,omitempty" gorm:"column:sort;type:tinyint(3);not null;default:0"`                              // 排序(越大越前)
	StartTime         *time.Time `json:"start_time,omitempty" gorm:"column:start_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`   // 开始时间
	EndTime           *time.Time `json:"end_time,omitempty" gorm:"column:end_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"`       // 结束时间
	State             int        `json:"state,omitempty" gorm:"column:state;type:tinyint(1);not null;default:1"`                            // 状态|1上线中|-1已下线|-2已删除
	AdminUserId       int        `json:"admin_user_id,omitempty" gorm:"column:admin_user_id;type:int(11);unsigned;not null;default:0"`      // 操作人ID
	AdminUserNickname string     `json:"admin_user_nickname,omitempty" gorm:"-"`                                                            // 操作人昵称(额外加的)
	CreateTime        *time.Time `json:"create_time,omitempty" gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"` // 创建时间
	UpdateTime        *time.Time `json:"update_time,omitempty" gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"` // 更新时间
}
