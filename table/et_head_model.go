package table

import (
	"time"
)

const EtHeadModelTN = `et_head_model`

type EtHeadModel struct {
	Id                int        `json:"id,omitempty" gorm:"column:id;type:int(11);unsigned;not null;AUTO_INCREMENT"`                       // 主键ID
	Model             int        `json:"model,omitempty" gorm:"column:model;type:tinyint(1);unsigned;not null;default:1"`                   // 模式|1banner模式|2无banner模式
	AdminUserId       int        `json:"admin_user_id,omitempty" gorm:"column:admin_user_id;type:int(11);unsigned;not null;default:0"`      // 操作人ID
	AdminUserNickname string     `json:"admin_user_nickname,omitempty" gorm:"-"`                                                            // 操作人昵称(额外加的)
	CreateTime        *time.Time `json:"create_time,omitempty" gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"` // 创建时间
	UpdateTime        *time.Time `json:"update_time,omitempty" gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP"` // 更新时间
}
