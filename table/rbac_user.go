package table

const RbacUserTN = `rbac_user`

type RbacUser struct {
	Id       int    `json:"id,omitempty" gorm:"column:id;type:int(11);not null;AUTO_INCREMENT"`
	Username string `json:"username,omitempty" gorm:"column:username;type:varchar(20);not null"`               // 用户名
	Password string `json:"password,omitempty" gorm:"column:password;type:varchar(32);not null"`               // 密码
	Nickname string `json:"nickname,omitempty" gorm:"column:nickname;type:varchar(20);not null"`               // 昵称
	Email    string `json:"email,omitempty" gorm:"column:email;type:varchar(25);not null"`                     // Email
	Mobile   string `json:"mobile,omitempty" gorm:"column:mobile;type:char(11);not null;default:''"`           // 手机号码
	RoleId   int    `json:"role_id,omitempty" gorm:"column:role_id;type:int(11) DEFAULT;null"`                 // 角色ID
	DeptId   int    `json:"dept_id,omitempty" gorm:"column:dept_id;type:int(11);not null;default:0"`           // 部门ID
	Leader   int    `json:"leader,omitempty" gorm:"column:leader;type:tinyint(1);unsigned;not null;default:0"` // 1-组长
	Status   int    `json:"status,omitempty" gorm:"column:status;type:int(11);not null;default:1"`             // 状态(1:正常|0:停用|-1暂未审核通过)
	Explain  string `json:"explain,omitempty" gorm:"column:explain;type:text;not null"`                        // 申请说明
}
