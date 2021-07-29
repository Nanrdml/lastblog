package main

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"lastblog/global"
	"lastblog/internal/model"
	"lastblog/internal/routers"
	"lastblog/pkg/logger"
	"lastblog/pkg/setting"
	"log"
	"net/http"
	"time"
)

func init(){
	err := setupSetting()
	if err != nil{
		log.Fatalf("init.setupSetting err: %v",err)
	}

	err = setupDBEngine()
	if err != nil{
		log.Fatalf("init.setupDBEngine err: %v",err)
	}

	err = setupLogger()
	if err != nil{
		log.Fatalf("init.setupLogger err: %v",err)
	}
}

func setupSetting() error{
	//获取到配置信息后将他与全局的结构体绑定
	setting, err := setting.NewSetting()
	if err != nil{
		return err
	}
	err = setting.ReadSection("Server",&global.ServerSetting)
	if err != nil{
		return err
	}

	err = setting.ReadSection("App",&global.AppSetting)
	if err != nil{
		return err
	}

	err = setting.ReadSection("Database",&global.DatabaseSetting)
	if err != nil{
		return err
	}

	//配置文件里没有单位
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	//fmt.Println(global.ServerSetting)
	//fmt.Println(global.AppSetting)
	//fmt.Println(global.DatabaseSetting)
	return nil
}

func setupDBEngine() error{

	//这里需要注意，有一些人会把初始化语句不小心写成：
	//global.DBEngine, err := model.NewDBEngine(global.DatabaseSetting)，这是存在很大问题的，
	//因为 := 会重新声明并创建了左侧的新局部变量，因此在其它包中调用 global.DBEngine 变量时，
	//它仍然是 nil，仍然是达不到可用标准，因为根本就没有赋值到真正需要赋值的包全局变量 global.DBEngine 上。
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil{
		return err
	}
	return nil
}

func setupLogger()error{
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename: global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags)
	return nil
}


// @title 博客系统
// @version 1.0
// @description Go 语言编程之旅：一起用 Go 做项目
// @termsOfService https://github.com/go-programming-tour-book
func main() {
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":8000",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	global.Logger.Infof("%s: go-programming-tour-book/%s", "eddycjy", "blog-service")
	s.ListenAndServe()

}
