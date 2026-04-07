package controller

import (
	"Bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//-----社区相关-----

// 访问
func CommunityHandler(c *gin.Context) {
	//查询所有的社区（community_id,community_name)以列表形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("Logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	//获取社区ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//查询所有的社区（community_id,community_name)以列表形式返回
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("Logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
