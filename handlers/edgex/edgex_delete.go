package edgex

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tydcwym/edgex_admin/caller"
	"github.com/tydcwym/edgex_admin/dal"
	"github.com/tydcwym/edgex_admin/logs"
	"github.com/tydcwym/edgex_admin/resp"
)

// DeleteEdgexParams ...
type DeleteEdgexParams struct {
	EdgexID int64 `form:"edgex_id" json:"edgex_id" binding:"required"`
}

type deleteEdgexHandler struct {
	Ctx    *gin.Context
	Params DeleteEdgexParams
}

func buildDeleteEdgexHandler(c *gin.Context) *deleteEdgexHandler {
	return &deleteEdgexHandler{
		Ctx: c,
	}
}

// DeleteEdgex ...
func DeleteEdgex(c *gin.Context) (out *resp.JSONOutput) {

	h := buildDeleteEdgexHandler(c)

	// Step1. checkParams
	err := h.CheckParams()
	if err != nil {
		logs.Warn("[DeleteEdgex] params-err: err=%v", err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}

	// Step2. update deleted
	err = h.Process()
	if err != nil {
		logs.Warn("[DeleteEdgex] params-err: err=%v", err)
		return resp.SampleJSON(c, resp.RespDatabaseError, nil)
	}

	return resp.SampleJSON(c, resp.RespCodeSuccess, nil)
}

func (h *deleteEdgexHandler) CheckParams() error {

	err := h.Ctx.Bind(&h.Params)
	if err != nil {
		logs.Error("[deleteEdgexHandler-checkParams] params-err: err=%v", err)
		return err
	}

	if h.Params.EdgexID <= 0 {
		logs.Error("[deleteEdgexHandler-checkParams] params-err: edgex_id=%v", h.Params.EdgexID)
		return fmt.Errorf("edgex_id is invalid: edgex_id=%v", h.Params.EdgexID)
	}
	return nil
}

func (h *deleteEdgexHandler) Process() (err error) {

	fieldsMap := map[string]interface{}{"deleted": 1}

	err = dal.UpdateEdgex(caller.EdgexDB, h.Params.EdgexID, fieldsMap)
	if err != nil {
		logs.Error("[deleteEdgexHandler-process] UpdateEdgex Failed: edgex_id=%+v, fileds=%+v, err=%+v",
			h.Params.EdgexID, fieldsMap, err)
		return
	}
	return
}
