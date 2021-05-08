package user

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tdycwym/edgex_admin/caller"
	"github.com/tdycwym/edgex_admin/dal"
	"github.com/tdycwym/edgex_admin/logs"
	"github.com/tdycwym/edgex_admin/middleware/session"
	"github.com/tdycwym/edgex_admin/resp"
	"github.com/tdycwym/edgex_admin/utils"

	"encoding/json"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
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

func sendMessage(c *gin.Context) *resp.JSONOutput {
	/* 必要步骤：
	 * 实例化一个认证对象，入参需要传入腾讯云账户密钥对 secretId 和 secretKey
	 * 本示例采用从环境变量读取的方式，需要预先在环境变量中设置这两个值
	 * 您也可以直接在代码中写入密钥对，但需谨防泄露，不要将代码复制、上传或者分享给他人
	 * CAM 密匙查询: https://console.cloud.tencent.com/cam/capi
	 */
	credential := common.NewCredential(
		// os.Getenv("TENCENTCLOUD_SECRET_ID"),
		// os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		"xxx",
		"xxx",
	)
	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()

	/* SDK 默认使用 POST 方法
	 * 如需使用 GET 方法，可以在此处设置，但 GET 方法无法处理较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"

	/* SDK 有默认的超时时间，非必要请不要进行调整
	 * 如有需要请在代码中查阅以获取最新的默认值 */
	//cpf.HttpProfile.ReqTimeout = 5

	/* SDK 会自动指定域名，通常无需指定域名，但访问金融区的服务时必须手动指定域名
	 * 例如 SMS 的上海金融区域名为 sms.ap-shanghai-fsi.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	/* SDK 默认用 TC3-HMAC-SHA256 进行签名，非必要请不要修改该字段 */
	cpf.SignMethod = "HmacSHA1"

	/* 实例化 SMS 的 client 对象
	 * 第二个参数是地域信息，可以直接填写字符串 ap-guangzhou，或者引用预设的常量 */
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	/* 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	* 您可以直接查询 SDK 源码确定接口有哪些属性可以设置
	 * 属性可能是基本类型，也可能引用了另一个数据结构
	 * 推荐使用 IDE 进行开发，可以方便地跳转查阅各个接口和数据结构的文档说明 */
	request := sms.NewSendSmsRequest()

	/* 基本类型的设置:
	 * SDK 采用的是指针风格指定参数，即使对于基本类型也需要用指针来对参数赋值。
	 * SDK 提供对基本类型的指针引用封装函数
	 * 帮助链接：
	 * 短信控制台：https://console.cloud.tencent.com/smsv2
	 * sms helper：https://cloud.tencent.com/document/product/382/3773
	 */

	/* 短信应用 ID: 在 [短信控制台] 添加应用后生成的实际 SDKAppID，例如1400006666 */
	request.SmsSdkAppid = common.StringPtr("1305849352")

	request.PhoneNumberSet = common.StringPtrs([]string{"+8618651886162", "+8613711333222", "+8613711144422"})

	// 通过 client 对象调用想要访问的接口，需要传入请求对象
	response, err := client.SendSms(request)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, "发送失败")
	}
	// 非 SDK 异常，直接失败。实际代码中可以加入其他的处理
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(response.Response)
	// 打印返回的 JSON 字符串
	fmt.Printf("%s", b)
	return resp.SampleJSON(c, resp.RespCodeSuccess, nil)
}
