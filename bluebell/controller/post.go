package controller

import (
	"Bluebell/logic"
	"Bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 创建帖子的处理函数
func CreatPostHandler(c *gin.Context) {
	//获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	//从c中取得当前请求用户的ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorId = int64(userID)

	//创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failes", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, nil)

}

// 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {
	//从URL获取参数
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//根据ID取出帖子数据
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

// 获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {
	//获取数据
	page, size := getPageInfo(c)

	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// 根据前端传来的参数动态的获取帖子列表
// 按创建时间 或者 分数 排序
// 1.获取参数
// 2.去redis查询id列表
// 3.根据id去数据库查询详细信息
func GetPostListHandler2(c *gin.Context) {
	//GET请求参数： /api/v1/posts2?page=1&size=10&order=time
	//获取数据
	//初始化结构体参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetPostListNew(p) //更新合二为一

	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// func GetCommunityPostListHandler(c *gin.Context) {

// 	p := &models.ParamCommunityPostList{
// 		ParamPostList: models.ParamPostList{
// 		Page:  1,
// 		Size:  10,
// 		Order: models.OrderTime,
// 		},

// 	}
// 	if err := c.ShouldBindQuery(p); err != nil {
// 		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
// 		ResponseError(c, CodeInvalidParam)
// 		return
// 	}

// 	data, err := logic.GetCommunityPostList2(p)
// 	if err != nil {
// 		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
// 		ResponseError(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuccess(c, data)
// }
