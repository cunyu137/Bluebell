package mysql

import (
	//"encoding/json"
	"time"

	"go.uber.org/zap"
)

// StartMonitor 启动连接池监控
func StartMonitor() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			stats := db.Stats()

			// 记录关键指标
			zap.L().Info("MySQL连接池状态",
				zap.Int("最大连接数", stats.MaxOpenConnections),
				zap.Int("当前打开数", stats.OpenConnections),
				zap.Int("使用中", stats.InUse),
				zap.Int("空闲", stats.Idle),
				zap.Uint64("等待次数", uint64(stats.WaitCount)),
				zap.Duration("等待总时长", stats.WaitDuration))

			// 如果等待次数在增加，说明连接池不够用
			if stats.WaitCount > 0 {
				zap.L().Warn("连接池等待队列增长",
					zap.Uint64("wait_count", uint64(stats.WaitCount)))
			}
		}
	}()
}
