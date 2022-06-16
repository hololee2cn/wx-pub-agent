package entity

// User 用户表
type User struct {
	// 主键ID
	ID int `json:"id" gorm:"id"`
	// 用户open id
	OpenID string `json:"open_id" gorm:"open_id"`
	// 用户名
	Name string `json:"name" gorm:"name"`
	// 联系号码
	Phone string `json:"phone" gorm:"phone"`
	// 创建时间
	CreateTime int64 `json:"create_time" gorm:"create_time"`
	// 删除时间
	DeleteTime int64 `json:"delete_time" gorm:"delete_time"`
}

func (u User) TableName() string {
	return "user"
}
