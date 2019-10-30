package user

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"review-order/app/model/db"
	"review-order/app/model/sys"
	"review-order/app/pkg/casbin"
	"review-order/app/tools"
)

type Role struct{}

// 所有角色
func (Role) Allrole(r *ghttp.Request) {
	var list []sys.Role
	err := db.Find(&sys.Role{}, &list, "parent_id asc", "sequence asc")
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	tools.SuccessJson(r, &list)
}

// 分页数据
func (Role) List(r *ghttp.Request) {
	page := tools.GetPageIndex(r)
	limit := tools.GetPageLimit(r)
	sort := tools.GetPageSort(r)
	key := tools.GetPageKey(r)
	parent_id := r.GetQueryUint64("parent_id")
	var whereOrder []db.PageWhereOrder
	order := "ID DESC"
	if len(sort) >= 2 {
		orderType := sort[0:1]
		order = sort[1:len(sort)]
		if orderType == "+" {
			order += " ASC"
		} else {
			order += " DESC"
		}
	}
	whereOrder = append(whereOrder, db.PageWhereOrder{Order: order})
	if key != "" {
		v := "%" + key + "%"
		var arr []interface{}
		arr = append(arr, v)
		whereOrder = append(whereOrder, db.PageWhereOrder{Where: "name like ?", Value: arr})
	}
	if parent_id > 0 {
		var arr []interface{}
		arr = append(arr, parent_id)
		whereOrder = append(whereOrder, db.PageWhereOrder{Where: "parent_id = ?", Value: arr})
	}
	var total uint64
	list := []sys.Role{}
	err := db.GetPage(&sys.Role{}, &sys.Role{}, &list, page, limit, &total, whereOrder...)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	tools.ResSuccessPage(r, total, &list)
}

//新增
func (Role) Create(r *ghttp.Request) {
	model := sys.Role{}
	jsonRequest := r.GetRaw()
	err := json.Unmarshal(jsonRequest, &model)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	err = db.Create(&model)
	if err != nil {
		tools.ErrorJson(r, "操作失败")
		return
	}
	tools.SuccessJson(r, g.Map{"id": model.ID})
}

// 详情
func (Role) Detail(r *ghttp.Request) {
	id := r.GetQueryUint64("id")
	var model sys.Role
	where := sys.Role{}
	where.ID = id
	_, err := db.First(&where, &model)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	tools.SuccessJson(r, &model)
}

// 更新
func (Role) Update(r *ghttp.Request) {
	model := sys.Role{}
	jsonRequest := r.GetRaw()
	err := json.Unmarshal(jsonRequest, &model)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	err = db.Save(&model)
	if err != nil {
		tools.ErrorJson(r, "操作失败")
		return
	}
	tools.SuccessJson(r)
}

// 删除数据
func (Role) Delete(r *ghttp.Request) {
	var ids []uint64
	jsonRequest := r.GetRaw()
	err := json.Unmarshal(jsonRequest, &ids)
	if err != nil || len(ids) == 0 {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	role := sys.Role{}
	err = role.Delete(ids)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	go casbin.CsbinDeleteRole(ids)
	tools.SuccessJson(r)
}

// 获取角色下的菜单ID列表
func (Role) Rolemenuidlist(r *ghttp.Request) {
	roleid := r.GetQueryUint64("roleid")
	menuIDList := []uint64{}
	where := sys.RoleMenu{RoleID: roleid}
	err := db.PluckList(&sys.RoleMenu{}, &where, &menuIDList, "menu_id")
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	tools.SuccessJson(r, &menuIDList)
}
