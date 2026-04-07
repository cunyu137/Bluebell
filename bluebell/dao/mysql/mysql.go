package mysql

import (
	"Bluebell/settings"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	// ✅ 记录成功日志
	zap.L().Info("数据库连接成功",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.DbName))

	// ✅ 添加日志，确认配置生效
	zap.L().Info("数据库连接池配置",
		zap.Int("MaxOpenConns", cfg.MaxOpenConns),
		zap.Int("MaxIdleConns", cfg.MaxIdleConns),
		zap.Int("实际设置的最大连接数", db.Stats().MaxOpenConnections))
	return
}

func Close() {
	_ = db.Close()
}

func GetStats() sql.DBStats {
	return db.Stats()
}
