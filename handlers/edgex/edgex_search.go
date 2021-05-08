package edgex

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tydcwym/edgex_admin/constdef"
	"github.com/tydcwym/edgex_admin/dal"
	"github.com/tydcwym/edgex_admin/logs"
	"github.com/tydcwym/edgex_admin/middleware/session"
	"github.com/tydcwym/edgex_admin/model"
	"github.com/tydcwym/edgex_admin/resp"
	"github.com/tydcwym/edgex_admin/utils"
)

const (
	ActionAll    = "all"    // 全部
	ActionMe     = "me"     // 我的创建
	ActionFollow = "follow" // 我的关注
)

// SearchEdgexParams ...
type SearchEdgexParams struct {
	Action   string `form:"action" json:"action"`
	UserID   int64
	Username string
	Keyword  string `form:"keyword" json:"keyword"`
	Status   int    `form:"status" json:"status"`
	Offset   int    `form:"offset" json:"offset"`
	Count    int    `form:"count" json:"count"`
}

type searchEdgexHandler struct {
	Ctx       *gin.Context
	Params    SearchEdgexParams
	EdgexList []*model.EdgexInfo
}

func buildSearchEdgexHandler(c *gin.Context) *searchEdgexHandler {
	return &searchEdgexHandler{
		Ctx:       c,
		EdgexList: make([]*model.EdgexInfo, 0),
	}
}

// SearchEdgex ...
func SearchEdgex(c *gin.Context) (out *resp.JSONOutput) {

	h := buildSearchEdgexHandler(c)

	// Step1. checkParams
	err := h.CheckParams()
	if err != nil {
		logs.Error("[SearchEdgex] params-err: err=%v", err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}

	// Step2. search
	err = h.Process()
	if err != nil {
		logs.Error("[SearchEdgex] params-err: err=%v", err)
		return resp.SampleJSON(c, resp.RespDatabaseError, nil)
	}

	return resp.SampleJSON(c, resp.RespCodeSuccess, h.EdgexList)
}

func (h *searchEdgexHandler) CheckParams() error {

	err := h.Ctx.Bind(&h.Params)
	if err != nil {
		logs.Error("[searchEdgexHandler-checkParams] params-err: err=%v", err)
		return err
	}

	h.Params.UserID = session.GetSessionUserID(h.Ctx)
	h.Params.Username = session.GetSessionUsername(h.Ctx)

	if h.Params.Action == "" {
		h.Params.Action = ActionAll
	}

	if !utils.InStringSlice(h.Params.Action, []string{ActionAll, ActionMe, ActionFollow}) {
		return fmt.Errorf("action is invalid: action=%s", h.Params.Action)
	}

	if (h.Params.Action == ActionMe || h.Params.Action == ActionFollow) && h.Params.UserID == 0 {
		return fmt.Errorf("params error: action=%s but user_id=0", h.Params.Action)
	}

	if h.Params.Count == 0 {
		h.Params.Count = 10
	}
	return nil
}

func (h *searchEdgexHandler) Process() (err error) {

	var (
		edgexIDs []int64
		userIDs  []int64
		keyword  = h.Params.Keyword
	)

	followMap, err := dal.GetFollowMapByUserID(h.Params.UserID)
	if err != nil {
		return
	}

	switch h.Params.Action {
	case ActionMe:
		userIDs = []int64{h.Params.UserID}
	case ActionFollow:
		for edgexID := range followMap {
			edgexIDs = append(edgexIDs, edgexID)
		}
		if err != nil {
			return
		}
		if len(edgexIDs) == 0 {
			return nil
		}
	}

	// 指定edgex_id搜索
	if edgexID, err := strconv.ParseInt(keyword, 10, 64); err == nil && edgexID > 0 {
		if h.Params.Action == ActionFollow && !utils.InInt64Slice(edgexID, edgexIDs) {
			return nil
		}
		edgexIDs = []int64{edgexID}
		keyword = ""
	}

	edgexList, err := dal.GetEdgexList(edgexIDs, userIDs, keyword, h.Params.Status, h.Params.Offset, h.Params.Count)

	if err != nil {
		logs.Error("[searchEdgexHandler-Process] GetEdgexList failed: err=%v", err)
		return err
	}

	h.Pack(edgexList, followMap)
	return
}

func (h *searchEdgexHandler) Pack(edgexList []*dal.EdgexServiceItem, followMap map[int64]bool) {

	h.EdgexList = make([]*model.EdgexInfo, 0)

	for _, item := range edgexList {
		h.EdgexList = append(h.EdgexList, &model.EdgexInfo{
			EdgexID:          item.ID,
			EdgexName:        item.EdgexName,
			UserID:           item.UserID,
			UserName:         h.Params.Username,
			Prefix:           item.Prefix,
			Address:          item.Address,
			Status:           item.Status,
			CreatedTimeStamp: item.CreatedTime.Unix(),
			CreatedTime:      item.CreatedTime.Format(constdef.TimeFormat),
			Description:      item.Description,
			Location:         item.Location,
			Extra:            item.Extra,
			IsFollow:         followMap[item.ID],
		})
	}
}
