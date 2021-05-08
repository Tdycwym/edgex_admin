package edgex

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/tydcwym/edgex_admin/caller"
	"github.com/tydcwym/edgex_admin/dal"
	"github.com/tydcwym/edgex_admin/logs"
	"github.com/tydcwym/edgex_admin/resp"
)

// UpdateEdgexParams ...
type UpdateEdgexParams struct {
	EdgexID     int64  `form:"edgex_id" json:"edgex_id" binding:"required"`
	EdgexName   string `form:"edgex_name" json:"edgex_name"`
	Prefix      string `form:"prefix" json:"prefix"`
	Description string `form:"description" json:"description"`
	Address     string `form:"address" json:"address"`
	Location    string `form:"location" json:"location"`
	Extra       string `form:"extra" json:"extra"`
}

type updateEdgexHandler struct {
	Ctx    *gin.Context
	Params UpdateEdgexParams
}

func buildUpdateEdgexHandler(c *gin.Context) *updateEdgexHandler {
	return &updateEdgexHandler{
		Ctx: c,
	}
}

// UpdateEdgex ...
func UpdateEdgex(c *gin.Context) (out *resp.JSONOutput) {

	h := buildUpdateEdgexHandler(c)

	// Step1. checkParams
	err := h.CheckParams()
	if err != nil {
		logs.Warn("[DeleteEdgex] params-err: err=%v", err)
		return resp.SampleJSON(c, resp.RespCodeParamsError, nil)
	}

	// Step2. update
	err = h.Process()
	if err != nil {
		logs.Warn("[DeleteEdgex] params-err: err=%v", err)
		return resp.SampleJSON(c, resp.RespDatabaseError, nil)
	}

	return resp.SampleJSON(c, resp.RespCodeSuccess, nil)
}

func (h *updateEdgexHandler) CheckParams() error {

	err := h.Ctx.Bind(&h.Params)
	if err != nil {
		logs.Error("[updateEdgexHandler-checkParams] params-err: err=%v", err)
		return err
	}

	if h.Params.EdgexID <= 0 {
		logs.Error("[updateEdgexHandler-checkParams] params-err: edgex_id=%v", h.Params.EdgexID)
		return fmt.Errorf("edgex_id is invalid: edgex_id=%v", h.Params.EdgexID)
	}

	if h.Params.Prefix != "" {
		if match, _ := regexp.MatchString("^[a-z_-]+$", h.Params.Prefix); !match {
			logs.Error("[updateEdgexHandler-checkParams] params-err: prefix=%v", h.Params.Prefix)
			return fmt.Errorf("prefix is invalid: prefix=%v", h.Params.Prefix)
		}
	}

	return nil
}

func (h *updateEdgexHandler) Process() (err error) {

	fieldsMap := h.GetUpdateFieldsMap()
	if len(fieldsMap) == 0 {
		return nil
	}
	err = dal.UpdateEdgex(caller.EdgexDB, h.Params.EdgexID, fieldsMap)
	if err != nil {
		logs.Error("[updateEdgexHandler-process] UpdateEdgex Failed: edgex_id=%+v, fileds=%+v, err=%+v",
			h.Params.EdgexID, fieldsMap, err)
		return
	}
	return
}

func (h *updateEdgexHandler) GetUpdateFieldsMap() (fieldsMap map[string]interface{}) {

	fieldsMap = make(map[string]interface{})

	if h.Params.EdgexName != "" {
		fieldsMap["edgex_name"] = h.Params.EdgexName
	}

	if h.Params.Prefix != "" {
		fieldsMap["prefix"] = h.Params.Prefix
	}

	if h.Params.Description != "" {
		fieldsMap["description"] = h.Params.Description
	}

	if h.Params.Location != "" {
		fieldsMap["location"] = h.Params.Location
	}

	if h.Params.Extra != "" {
		fieldsMap["extra"] = h.Params.Extra
	}

	if h.Params.Address != "" {
		fieldsMap["address"] = h.Params.Address
	}
	return
}
