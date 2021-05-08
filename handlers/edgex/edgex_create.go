package edgex

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tydcwym/edgex_admin/caller"
	"github.com/tydcwym/edgex_admin/dal"
	"github.com/tydcwym/edgex_admin/logs"
	"github.com/tydcwym/edgex_admin/middleware/session"
	"github.com/tydcwym/edgex_admin/resp"
)

// CreateEdgexParams ...
type CreateEdgexParams struct {
	UserID      int64
	Username    string
	EdgexName   string `form:"edgex_name" json:"edgex_name" binding:"required"`
	Prefix      string `form:"prefix" json:"prefix" binding:"required"`
	Description string `form:"description" json:"description" binding:"required"`
	Address     string `form:"address" json:"address"`
	Location    string `form:"location" json:"location"`
	Extra       string `form:"extra" json:"extra"`
}

type createEdgexHandler struct {
	Ctx    *gin.Context
	Params CreateEdgexParams
}

func buildCreateEdgexHandler(c *gin.Context) *createEdgexHandler {
	return &createEdgexHandler{
		Ctx: c,
	}
}

// CreateEdgex ...
func CreateEdgex(c *gin.Context) (out *resp.JSONOutput) {

	h := buildCreateEdgexHandler(c)

	// Step1. checkParams
	err := h.CheckParams()
	if err != nil {
		logs.Warn("[CreateEdgex] params-err: err=%v", err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}

	// Step2. createEdgexAndFollow
	err = h.Process()
	if err != nil {
		logs.Warn("[CreateEdgex] params-err: err=%v", err)
		return resp.SampleJSON(c, resp.RespDatabaseError, nil)
	}

	return resp.SampleJSON(c, resp.RespCodeSuccess, nil)
}

func (h *createEdgexHandler) CheckParams() error {

	err := h.Ctx.Bind(&h.Params)
	if err != nil {
		logs.Error("[createEdgexHandler-checkParams] params-err: err=%v", err)
		return err
	}

	if match, _ := regexp.MatchString("^[a-z_-]+$", h.Params.Prefix); !match {
		logs.Error("[createEdgexHandler-checkParams] params-err: prefix=%v", h.Params.Prefix)
		return fmt.Errorf("prefix is invalid: prefix=%v", h.Params.Prefix)
	}

	h.Params.UserID = session.GetSessionUserID(h.Ctx)
	h.Params.Username = session.GetSessionUsername(h.Ctx)
	return nil
}

func (h *createEdgexHandler) Process() (err error) {
	edgex := h.ConvertEdgexItem(h.Params)

	db := caller.EdgexDB.Begin()
	defer func() {
		if err != nil {
			db.Callback()
		} else {
			db.Commit()
		}
	}()
	err = dal.AddEdgex(db, edgex)
	if err != nil {
		logs.Error("[createEdgexHandler-Process] AddEdgex Failed: edgex=%+v, err=%+v", edgex, err)
		return
	}
	item := &dal.EdgexRelatedUser{
		UserID:       h.Params.UserID,
		Username:     h.Params.Username,
		EdgexID:      edgex.ID,
		EdgexName:    edgex.EdgexName,
		Status:       dal.StatusFollow,
		CreatedTime:  time.Now(),
		ModifiedTime: time.Now(),
	}
	err = dal.AddEdgexRelatedUser(db, item)
	if err != nil {
		logs.Error("[createEdgexHandler-Process] AddEdgexRelatedUser Failed: item=%+v, err=%+v", item, err)
		return
	}
	return
}

func (h *createEdgexHandler) ConvertEdgexItem(params CreateEdgexParams) *dal.EdgexServiceItem {
	return &dal.EdgexServiceItem{
		UserID:       params.UserID,
		EdgexName:    params.EdgexName,
		Prefix:       params.Prefix,
		Description:  params.Description,
		Location:     params.Location,
		Extra:        params.Extra,
		Address:      params.Address,
		CreatedTime:  time.Now(),
		ModifiedTime: time.Now(),
	}
}
