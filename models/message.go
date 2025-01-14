package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId string
	TargetId string
	Type string // 消息类型 群聊 私聊 广播
	Media int // 消息类型 文字 图片 音频
	Content string
	Pic string
	Url string
	Desc string
	Amount int // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}
