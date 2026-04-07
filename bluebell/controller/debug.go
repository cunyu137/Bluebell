package controller

import (
	"Bluebell/dao/mysql"
	"Bluebell/dao/redis"

	"github.com/gin-gonic/gin"
)

func PoolStats(c *gin.Context) {
	dbStats := mysql.GetStats()
	redisStats := redis.GetPoolStats()

	c.JSON(200, gin.H{
		"database": gin.H{
			"open_connections": dbStats.OpenConnections,
			"in_use":           dbStats.InUse,
			"idle":             dbStats.Idle,
			"wait_count":       dbStats.WaitCount,
			"wait_duration":    dbStats.WaitDuration.String(),
		},
		"redis": gin.H{
			"total_conns": redisStats.TotalConns,
			"idle_conns":  redisStats.IdleConns,
			"hit_rate":    float64(redisStats.Hits) / float64(redisStats.Hits+redisStats.Misses+1),
		},
	})
}
