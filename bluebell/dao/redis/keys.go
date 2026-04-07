package redis

const (
	KeyPrefix              = "bluebell:"
	KeyPostTimeZSet        = "post:time"  //zset:帖子及发帖时间
	KeyPostScoreZSet       = "post:score" //zset:帖子及头片分数
	KeyPostVotedZSetPrefix = "post:vote"  //zset:记录用户及投票类型 参数是post id
	KeyCommunitySetPrefix  = "community"  //set:保存每个分区下帖子的id

)

// 给redisKey加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
