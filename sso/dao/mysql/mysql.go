/**
 * @Author: Robby
 * @File name: mysql.go
 * @Create date: 2021-05-18
 * @Function:
 **/

package mysqlconnect

import (
	"fmt"
	"log"
	"sso/sso/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(cfg *settings.MysqlConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)

	defer func() {
		r := recover()
		if r != nil {
			log.Fatalf("连接数据库失败: %s\n", r)
		}
	}()

	// 如果连接不上，这里会panic
	db = sqlx.MustConnect("mysql", dsn)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return

}

func Close() {
	_ = db.Close()
}
