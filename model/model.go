package model

import "time"

type APIResponse struct {
	ID            uint      `gorm:"primaryKey"`        // 主键
	Endpoint      string    `gorm:"type:varchar(255)"` // 接口路径，例如 /api/v1/resource
	RequestParams string    `gorm:"type:text"`         // 请求参数，通常以 JSON 格式存储
	ResponseData  string    `gorm:"type:text"`         // 返回数据，以 JSON 格式存储
	CreatedAt     time.Time `gorm:"autoCreateTime"`    // 创建时间
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`    // 更新时间
}

func (APIResponse) TableName() string {
	return "api_responses" // 自定义表名
}
