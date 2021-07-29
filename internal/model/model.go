package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"lastblog/global"
	"lastblog/pkg/setting"
)

//本文件是公共model

// Model
//公共结构体，里面的字段也是数据库中表内公共字段
type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// NewDBEngine
//在创造这个函数时，项目已经读取完配置信息了，使用读取的信息创造连接
//这里的参数有点像依赖注入
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		//LogMode设置日志模式，' true '为详细日志，' false '为无日志，默认只打印错误日志
		db.LogMode(true)
	}
	db.SingularTable(true) //使用单数表名，就是说user表不会变成users

	//SetMaxIdleConns设置空闲连接池的最大连接数。
	//如果MaxOpenConns大于0但小于新的MaxIdleConns，那么新的MaxIdleConns将被减少以匹配MaxOpenConns限制。
	//当n <= 0时，不保留空闲连接。
	//当前默认的最大空闲连接数是2。这可能会在未来的版本中改变。
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}
