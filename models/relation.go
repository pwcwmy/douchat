package models

import (
	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	OwnerId uint
	TargetId uint
	Type int 
	Desc string
	
}

func (table *Relation) TableName() string {
	return "relation"
}
