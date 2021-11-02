/**
 * @Author: Robby
 * @File name: main.go
 * @Create date: 2021-11-02
 * @Function:
 **/

package main

import (
	"log"
	"os"
	mysqlconnect "sso/sso/dao/mysql"
	"sso/sso/logger"
	"sso/sso/settings"

	"go.uber.org/zap"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("参数个数为：%d, 需要带上配置文件，例如运行 ./sso config/config.yaml\n", len(os.Args)-1)
	}

	if err := settings.Init(os.Args[1]); err != nil {
		log.Fatalf("配置文件读取失败：%s\n", err.Error())
	} else {
		log.Println("配置文件读取成功")
	}

	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		log.Fatalf("日志配置失败：%v", err)
	} else {
		log.Println("日志配置成功")
	}
	defer zap.L().Sync()

	if err := mysqlconnect.Init(settings.Conf.MysqlConfig); err != nil {
		log.Fatalf("MySQL连接失败：%v", err)
	} else {
		log.Println("MySQL连接成功")
	}
	defer mysqlconnect.Close()

	settings.Conf.SessionConfig
}
