package user

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"review-order/app/model/db"
	"review-order/app/model/sys"
	"review-order/app/tools"
	"strings"
)

type MenuMeta struct {
	Title   string `json:"title"`   // 标题
	Icon    string `json:"icon"`    // 图标
	NoCache bool   `json:"noCache"` // 是不是缓存
}

type MenuModel struct {
	Path      string      `json:"path"`      // 路由
	Component string      `json:"component"` // 对应vue中的map name
	Name      string      `json:"name"`      // 菜单名称
	Hidden    bool        `json:"hidden"`    // 是否隐藏
	Meta      MenuMeta    `json:"meta"`      // 菜单信息
	Children  []MenuModel `json:"children"`  // 子级菜单
}

type Menu struct{}

// 新增菜单后自动添加菜单下的常规操作
func InitMenu(model sys.Menu) {
	if model.MenuType != 2 {
		return
	}
	add := sys.Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/create", Name: "新增", Sequence: 1, MenuType: 3, Code: model.Code + "Add", OperateType: "add"}
	db.Create(&add)
	del := sys.Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/delete", Name: "删除", Sequence: 2, MenuType: 3, Code: model.Code + "Del", OperateType: "del"}
	db.Create(&del)
	view := sys.Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/detail", Name: "查看", Sequence: 3, MenuType: 3, Code: model.Code + "View", OperateType: "view"}
	db.Create(&view)
	update := sys.Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/update", Name: "编辑", Sequence: 4, MenuType: 3, Code: model.Code + "Update", OperateType: "update"}
	db.Create(&update)
	list := sys.Menu{Status: 1, ParentID: model.ID, URL: model.URL + "/list", Name: "分页api", Sequence: 5, MenuType: 3, Code: model.Code + "List", OperateType: "list"}
	db.Create(&list)
}

func (menu *Menu) List(r *ghttp.Request) {
	page := tools.GetPageIndex(r)
	limit := tools.GetPageLimit(r)
	sort := tools.GetPageSort(r)
	key := tools.GetPageKey(r)
	menuType := r.GetQueryInt64("type")
	parent_id := r.GetQueryInt64("parent_id")
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
		whereOrder = append(whereOrder, db.PageWhereOrder{Where: "name like ? or code like ?", Value: arr})
	}
	if menuType > 0 {
		var arr []interface{}
		arr = append(arr, menuType)
		whereOrder = append(whereOrder, db.PageWhereOrder{Where: "menu_type = ?", Value: arr})
	}
	if parent_id > 0 {
		var arr []interface{}
		arr = append(arr, parent_id)
		whereOrder = append(whereOrder, db.PageWhereOrder{Where: "parent_id = ?", Value: arr})
	}
	var total uint64
	list := []sys.Menu{}
	err := db.GetPage(&sys.Menu{}, &sys.Menu{}, &list, page, limit, &total, whereOrder...)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	tools.ResSuccessPage(r, total, &list)
}

func (menu *Menu) Menubuttonlist(r *ghttp.Request) {
	// 用户ID
	uid := r.GetParam(tools.USER_ID_Key)
	userID := gconv.Uint64(uid)
	menuCode := r.GetQueryString("menucode")
	if menuCode == "" {
		tools.ErrorJson(r, "err")
		return
	}
	btnList := []string{}
	if userID == tools.SUPER_ADMIN_ID {
		//管理员
		btnList = append(btnList, "add")
		btnList = append(btnList, "del")
		btnList = append(btnList, "view")
		btnList = append(btnList, "update")
		btnList = append(btnList, "setrolemenu")
		btnList = append(btnList, "setadminrole")
	} else {
		menu := sys.Menu{}
		err := menu.GetMenuButton(userID, menuCode, &btnList)
		if err != nil {
			tools.ErrorJson(r, gconv.String(err))
			return
		}
	}
	tools.SuccessJson(r, &btnList)

}

func (menu *Menu) Allmenu(r *ghttp.Request) {
	var menus []sys.Menu
	err := db.Find(&sys.Menu{}, &menus, "parent_id asc", "sequence asc")
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	tools.SuccessJson(r, &menus)
}

// 详情
func (Menu) Detail(r *ghttp.Request) {
	id := r.GetQueryUint64("id")
	var menu sys.Menu
	where := sys.Menu{}
	where.ID = id
	_, err := db.First(&where, &menu)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	tools.SuccessJson(r, &menu)
}

// 更新
func (Menu) Update(r *ghttp.Request) {
	model := sys.Menu{}
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

//新增
func (Menu) Create(r *ghttp.Request) {
	menu := sys.Menu{}
	err := r.GetToStruct(&menu)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	err = db.Create(&menu)
	if err != nil {
		tools.ErrorJson(r, "操作失败")
		return
	}
	go InitMenu(menu)
	tools.SuccessJson(r, g.Map{"id": menu.ID})
}

// 删除数据
func (Menu) Delete(r *ghttp.Request) {

	stirngs := r.GetRawString()
	stirngs = strip(strip(stirngs, "["), "]")
	glog.Debug(stirngs)
	ids := strings.Split(stirngs, ",")
	glog.Debug(ids)

	if len(ids) == 0 {
		tools.ErrorJson(r, "err")
		return
	}
	menu := sys.Menu{}
	err := menu.Delete(ids)
	if err != nil {
		tools.ErrorJson(r, gconv.String(err))
		return
	}
	tools.SuccessJson(r)
}

func strip(s_ string, chars_ string) string {
	s, chars := []rune(s_), []rune(chars_)
	length := len(s)
	max := len(s) - 1
	l, r := true, true //标记当左端或者右端找到正常字符后就停止继续寻找
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[rune]bool) //创建字符集，也就是唯一的字符，方便后面判断是否存在
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r { // 如果左端和右端都没找到正常字符，那么表示该字符串没有正常字符
		return ""
	}
	return string(s[start : end+1])
}
