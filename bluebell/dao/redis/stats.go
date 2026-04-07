package redis

import (
	"time"

	"go.uber.org/zap"
)

func StartMonitor() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			stats := client.PoolStats()

			zap.L().Info("Redis连接池状态",
				zap.Uint32("总连接数", stats.TotalConns),
				zap.Uint32("空闲连接", stats.IdleConns),
				zap.Uint32("过期连接", stats.StaleConns),
				zap.Uint32("命中次数", stats.Hits),
				zap.Uint32("未命中", stats.Misses))

			// 如果未命中率太高，说明连接池配置有问题
			if stats.Misses > 0 && stats.Hits > 0 {
				hitRate := float64(stats.Hits) / float64(stats.Hits+stats.Misses)
				if hitRate < 0.8 {
					zap.L().Warn("Redis连接池命中率过低",
						zap.Float64("命中率", hitRate))
				}
			}
		}
	}()
}
