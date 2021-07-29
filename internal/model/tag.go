package model

import "lastblog/pkg/app"

type Tag struct {
	//继承公共的model类
	*Model
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}


func (t Tag) TableName() string {
	return "blog_tag"
}
