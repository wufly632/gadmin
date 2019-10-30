package tools

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

const (
	SUCCESS_CODE = 20000 //成功的状态码
	FAIL_CODE    = 30000 //失败的状态码
)

type ResponsePageData struct {
	Total uint64      `json:"total"`
	Items interface{} `json:"items"`
}

type ResponsePage struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    ResponsePageData `json:"data"`
}

// 标准返回结果数据结构封装。
// 返回固定数据结构的JSON:
// err:  错误码(0:成功, 1:失败, >1:错误码);
// msg:  请求结果信息;
// data: 请求结果,根据不同接口返回结果的数据结构不同;
func Json(r *ghttp.Request, status int, msg string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(g.Map{
		"code":    status,
		"message": msg,
		"data":    responseData,
	})
	r.Exit()
}

func SuccessJson(r *ghttp.Request, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(g.Map{
		"code":    SUCCESS_CODE,
		"message": "",
		"data":    responseData,
	})
	r.Exit()
}

func ErrorJson(r *ghttp.Request, msg string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(g.Map{
		"code":    FAIL_CODE,
		"message": msg,
		"data":    responseData,
	})
	r.Exit()
}

// 响应成功-分页数据
func ResSuccessPage(r *ghttp.Request, total uint64, list interface{}) {
	ret := ResponsePage{Code: SUCCESS_CODE, Message: "ok", Data: ResponsePageData{Total: total, Items: list}}
	r.Response.WriteJson(ret)
}
