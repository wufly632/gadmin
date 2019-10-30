package user

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"review-order/app/model/db"
	"review-order/app/model/sys"
	"review-order/app/tools"
)

type Admins struct{}

// 分页数据
func (Admins) List(r *ghttp.Request) {
	page := tools.GetPageIndex(r)
	limit := tools.GetPageLimit(r)
	sort := tools.GetPageSort(r)
	key := tools.GetPageKey(r)
	status := r.GetQueryUint("status")
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
		arr = append(arr, v)
		whereOrder = append(whereOrder, db.PageWhereOrder{Where: "user_name like ? or real_name like ?", Value: arr})
	}
	if status > 0 {
		var arr []interface{}
		arr = append(arr, status)
		whereOrder = append(whereOrder, db.PageWhereOrder{Where: "status = ?", Value: arr})
	}
	var total uint64
	list := []sys.Admin{}
	err := db.GetPage(&sys.Admin{}, &sys.Admin{}, &list, page, limit, &total, whereOrder...)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	tools.ResSuccessPage(r, total, &list)
}

// 详情
func (Admins) Detail(r *ghttp.Request) {
	id := r.GetQueryUint64("id")
	var model sys.Admin
	where := sys.Admin{}
	where.ID = id
	_, err := db.First(&where, &model)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	model.Password = ""
	tools.SuccessJson(r, &model)
}

// 更新
func (Admins) Update(r *ghttp.Request) {
	/*model := sys.Admin{}
	err := c.Bind(&model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	where := sys.Admins{}
	where.ID = model.ID
	modelOld := sys.Admins{}
	_, err = models.First(&where, &modelOld)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	model.UserName = modelOld.UserName
	model.Password = modelOld.Password
	err = models.Save(&model)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	common.ResSuccessMsg(c)*/
}

//新增
func (Admins) Create(r *ghttp.Request) {
	/*model := sys.Admin{}
	err := c.Bind(&model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	model.Password = hash.Md5String(common.MD5_PREFIX + model.Password)
	err = models.Create(&model)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	common.ResSuccess(c, gin.H{"id": model.ID})*/
}

// 删除数据
func (Admins) Delete(r *ghttp.Request) {
	/*var ids []uint64
	err := c.Bind(&ids)
	if err != nil || len(ids) == 0 {
		common.ResErrSrv(c, err)
		return
	}
	admin := sys.Admins{}
	err = admin.Delete(ids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccessMsg(c)*/
}

// 获取用户下的角色ID列表
func (Admins) Adminsroleidlist(r *ghttp.Request) {
	adminsid := r.GetQueryUint64("adminsid")
	roleList := []uint64{}
	where := sys.AdminRole{AdminsID: adminsid}
	err := db.PluckList(&sys.AdminRole{}, &where, &roleList, "role_id")
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	tools.SuccessJson(r, &roleList)
}

// 分配用户角色权限
func (Admins) SetRole(r *ghttp.Request) {
	/*adminsid := r.GetQueryUint64("adminsid")
	var roleids []uint64
	err := c.Bind(&roleids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	ar := sys.AdminRole{}
	err = ar.SetRole(adminsid, roleids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	go common.CsbinAddRoleForUser(adminsid)
	common.ResSuccessMsg(c)*/
}
