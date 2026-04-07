package logic

import (
	"Bluebell/dao/redis"
	"Bluebell/models"
	"strconv"
)

//投票功能
//1.用户投票的数据

//简化版投票分数

func VoteForPost(userID int64, p *models.ParamVoteData) error {
	//投票的限制
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
	//更新帖子的分数

	//记录用户为该帖子投票的数据
}
