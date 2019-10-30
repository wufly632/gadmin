package tools

import (
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/net/ghttp"
	"github.com/satori/go.uuid"
)

const (
	MD5_PREFIX     = "wufly"    //MD5加密前缀字符串
	TOKEN_KEY      = "X-Token"  //页面token键名
	USER_ID_Key    = "X-USERID" //页面用户ID键名
	USER_UUID_Key  = "X-UUID"
	SUPER_ADMIN_ID = 0 //超级管理员主键ID
)

func EncryptPassword(data string) string {
	md5_pass, _ := gmd5.EncryptString(data + MD5_PREFIX)
	return md5_pass
}

func GetUUID() string {
	return uuid.NewV4().String()
}

// 获取页码
func GetPageIndex(r *ghttp.Request) uint64 {
	return r.GetQueryUint64("page", 1)
}

// 获取每页记录数
func GetPageLimit(r *ghttp.Request) uint64 {
	limit := r.GetQueryUint64("limit", 20)
	if limit > 500 {
		limit = 20
	}
	return limit
}

// 获取排序信息
func GetPageSort(r *ghttp.Request) string {
	return r.GetQueryString("sort")
}

// 获取搜索关键词信息
func GetPageKey(r *ghttp.Request) string {
	return r.GetQueryString("key")
}
