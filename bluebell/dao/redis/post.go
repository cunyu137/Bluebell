package redis

import (
	"Bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1

	//ZRevRange 按分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id

	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	//确定查询索引的起始点
	return getIDsFromKey(key, p.Page, p.Size)

}

// 根据ids查询每篇帖子投赞成票的数据
// func GetPostVoteData(ids []string) (data []int64, err error) {

// 	data = make([]int64, 0, len(ids))
// 	for _, id := range ids {
// 		key := getRedisKey(KeyPostVotedZSetPrefix + id)
// 		//查找每篇帖子中分数是1的数量
// 		v := client.ZCount(key, "1", "1").Val()
// 		data = append(data, v)
// 	}
// 	return

// }

// 按社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {

	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	//使用zinterstore把分区的帖子set和帖子分数zset 生成一个新zset
	//针对新的zset 按之前的逻辑读取数据

	//社区的key
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		//不存在
		pipline := client.Pipeline()
		pipline.ZInterStore(key, redis.ZStore{
			Aggregate: "MaX",
		}, cKey, orderKey)
		pipline.Expire(key, 69*time.Second) //设置超时时间
		_, err := pipline.Exec()
		if err != nil {
			return nil, err
		}

	}
	//从redis获取id

	return getIDsFromKey(key, p.Page, p.Size)

}

func GetPostVoteData(ids []string) (data []int64, err error) {
	if len(ids) == 0 {
		return []int64{}, nil
	}

	// 使用 Pipeline 批量发送命令
	pipe := client.Pipeline()
	cmds := make([]*redis.IntCmd, 0, len(ids))

	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		cmd := pipe.ZCount(key, "1", "1")
		cmds = append(cmds, cmd)
	}

	// 执行所有命令（一次网络往返）
	_, err = pipe.Exec()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// 收集结果
	data = make([]int64, 0, len(ids))
	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil && err != redis.Nil {
			// 如果某个key不存在，返回0
			data = append(data, 0)
			continue
		}
		data = append(data, val)
	}

	return data, nil
}
