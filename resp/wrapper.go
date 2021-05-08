package resp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tydcwym/edgex_admin/logs"
)

// JSONOutPutWrapper ...
func JSONOutPutWrapper(call func(*gin.Context) *JSONOutput) func(c *gin.Context) {
	return func(c *gin.Context) {
		var output *JSONOutput

		logs.Info("[wraper-request] url=%s, header=%v, body=%v",
			c.Request.URL, c.Request.Header, c.Request.Body)

		start := time.Now()

		defer func() {
			if tErr := recover(); tErr != nil {
				const size = 64 << 10
				buffer := make([]byte, size)
				buffer = buffer[:runtime.Stack(buffer, false)]
				logs.Error("[wrapper-panic] error=%v, stack=%s", tErr, buffer)

				rsp := NewStdResponse(RespCodeServerException, nil)
				output = NewJSONOutput(c, http.StatusInternalServerError, rsp)
			}
			if output == nil {
				logs.Error("[wraper-output-empty] output is empty!")
				rsp := NewStdResponse(RespCodeServerException, nil)
				output = NewJSONOutput(c, http.StatusInternalServerError, rsp)
			}

			output.Write()

			userTime := time.Since(start).Nanoseconds() / 1000
			logs.Info("[wraper-response] useTime=%d, status=%d, resp=%s",
				userTime, output.HTTPStatus, GetMarshalStr(output.Resp))
		}()
		output = call(c)
	}
}

// GetMarshalStr ...
func GetMarshalStr(obj interface{}) string {
	vi := reflect.ValueOf(obj)
	if vi.Kind() == reflect.Ptr && vi.IsNil() {
		return ""
	}
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("json Marshal failed: obj=%v, err=%v", obj, err)
	}
	return string(objBytes)
}
