package user

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tydcwym/edgex_admin/caller"
	"github.com/tydcwym/edgex_admin/dal"
	"github.com/tydcwym/edgex_admin/logs"
	"github.com/tydcwym/edgex_admin/middleware/session"
	"github.com/tydcwym/edgex_admin/resp"
	"github.com/tydcwym/edgex_admin/utils"
)

type LoginParams struct {
	UserID   int64  `form:"user_id" json:"user_id"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) *resp.JSONOutput {
	// Step1. 查看用户是否已登陆
	if session.GetSessionUserID(c) > 0 {
		return resp.SampleJSON(c, resp.RespCodeSuccess, "用户已登陆")
	}

	// Step2. 参数校验
	params := &LoginParams{}
	err := c.Bind(&params)
	if err != nil || (params.UserID <= 0 && params.Username == "") {
		logs.Error("[Login] request-params error: params=%+v, err=%v", params, err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}

	var (
		userInfo *dal.EdgexUser
		dbErr    error
	)

	// Step3. 查看用户是否存在
	if params.UserID > 0 {
		userInfo, dbErr = dal.GetEdgexUserByID(params.UserID)
	} else if params.Username != "" {
		userInfo, dbErr = dal.GetEdgexUserByName(params.Username)
	}

	if dbErr != nil {
		logs.Error("[Login] get userInfo failed: params=%+v, err=%v", params, err)
		return resp.SampleJSON(c, resp.RespDatabaseError, nil)
	}
	if userInfo == nil {
		logs.Error("[Login] user is Not Exsit: params=%+v, userInfo=%+v", params, userInfo)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}

	// Step4. 密码比对
	err = utils.Compare(userInfo.Password, params.Password)
	if err != nil {
		logs.Error("[Login] password is invalid:  params=%+v, err=%v", params, err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}

	// step5. session save
	session.SaveAuthSession(c, userInfo.ID, userInfo.Username)
	return resp.SampleJSON(c, resp.RespCodeSuccess, nil)
}

// RegisterParams ...
type RegisterParams struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Register(c *gin.Context) *resp.JSONOutput {
	// Step1. 参数校验
	params := &RegisterParams{}
	err := c.Bind(&params)
	if err != nil {
		logs.Error("[Register] request-params error: params=%+v, err=%v", params, err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}

	// Step2. 查看用户是否存在
	userInfo, dbErr := dal.GetEdgexUserByName(params.Username)
	if dbErr != nil {
		logs.Error("[Register] get user failed: username=%s, err=%v", params.Username, dbErr)
		return resp.SampleJSON(c, resp.RespDatabaseError, nil)
	}
	if userInfo != nil {
		return resp.SampleJSON(c, resp.RespCodeUserExsit, nil)
	}

	// Step3. 添加用户
	user := &dal.EdgexUser{
		Username:     params.Username,
		Password:     params.Password,
		CreatedTime:  time.Now(),
		ModifiedTime: time.Now(),
	}
	dbErr = dal.AddEdgexUser(caller.EdgexDB, user)
	if dbErr != nil {
		return resp.SampleJSON(c, resp.RespDatabaseError, nil)
	}
	return resp.SampleJSON(c, resp.RespCodeSuccess, nil)
}

func Logout(c *gin.Context) *resp.JSONOutput {
	userID := session.GetSessionUserID(c)
	if userID == 0 {
		return resp.SampleJSON(c, resp.RespCodeParamsError, "用户未登录")
	}
	session.ClearAuthSession(c)
	return resp.SampleJSON(c, resp.RespCodeSuccess, nil)
}
