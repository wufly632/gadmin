package user

import (
	"github.com/ahmetb/go-linq"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/util/gconv"
	"review-order/app/model/db"
	"review-order/app/model/sys"
	"review-order/app/pkg/casbin"
	"review-order/app/pkg/jwt"
	"review-order/app/tools"
	"time"
)

type User struct{}

type UserData struct {
	Menus        []MenuModel `json:"menus"`        // 菜单
	Introduction string      `json:"introduction"` // 介绍
	Avatar       string      `json:"avatar"`       // 图标
	Name         string      `json:"name"`         // 姓名
}

/**
* 用户登录
 */
func (user *User) Login(c *ghttp.Request) {
	row := c.GetRaw()
	j, err := gjson.DecodeToJson(row)
	if err != nil {
		panic(err)
	}
	username := j.GetString("username")
	password := j.GetString("password")
	if username == "" || password == "" {
		tools.ErrorJson(c, "用户名或密码不能为空")
		return
	}
	password = tools.EncryptPassword(password)
	where := sys.Admin{UserName: username, Password: password}
	userModel := sys.Admin{}
	notFound, err := db.First(&where, &userModel)
	if err != nil {
		if notFound {
			tools.ErrorJson(c, "用户名或密码错误")
			return
		}
		return
	}
	if userModel.Status != 1 {
		tools.ErrorJson(c, "该用户已被禁用")
		return
	}
	// 缓存或者redis
	uuid := tools.GetUUID()
	gcache.Set(uuid, userModel.ID, 3600*24*3650000)

	// token jwt
	userInfo := make(map[string]string)
	userInfo["exp"] = gconv.String(time.Now().Add(time.Hour * time.Duration(1)).Unix()) // 1H
	userInfo["iat"] = gconv.String(time.Now().Unix())
	userInfo["uuid"] = uuid
	token := jwt.CreateToken(userInfo)
	// 发至页面
	resData := make(map[string]string)
	resData["token"] = token
	//casbin 处理
	err = casbin.CasbinAddRoleForUser(userModel.ID)
	if err != nil {
		tools.ErrorJson(c, gconv.String(err))
		return
	}
	tools.SuccessJson(c, resData)
}

func (user *User) Info(c *ghttp.Request) {
	// 用户ID
	uid := c.GetParam(tools.USER_ID_Key)
	if uid == nil {
		tools.ErrorJson(c, "token 无效")
		return
	}
	userID := gconv.Uint64(uid)
	// 根据用户ID获取用户权限菜单
	var menuData []sys.Menu
	var err error
	if userID == tools.SUPER_ADMIN_ID {
		//管理员
		menuData, err = getAllMenu()
		if err != nil {
			tools.ErrorJson(c, gconv.String(err))
			return
		}
		if len(menuData) == 0 {
			menuModelTop := sys.Menu{Status: 1, ParentID: 0, URL: "", Name: "TOP", Sequence: 1, MenuType: 1, Code: "TOP", OperateType: "none"}
			db.Create(&menuModelTop)
			menuModelSys := sys.Menu{Status: 1, ParentID: menuModelTop.ID, URL: "", Name: "系统管理", Sequence: 1, MenuType: 1, Code: "Sys", Icon: "lock", OperateType: "none"}
			db.Create(&menuModelSys)
			menuModel := sys.Menu{Status: 1, ParentID: menuModelSys.ID, URL: "/icon", Name: "图标管理", Sequence: 10, MenuType: 2, Code: "Icon", Icon: "icon", OperateType: "none"}
			db.Create(&menuModel)
			menuModel = sys.Menu{Status: 1, ParentID: menuModelSys.ID, URL: "/menu", Name: "菜单管理", Sequence: 20, MenuType: 2, Code: "Menu", Icon: "documentation", OperateType: "none"}
			db.Create(&menuModel)
			InitMenu(menuModel)
			menuModel = sys.Menu{Status: 1, ParentID: menuModelSys.ID, URL: "/role", Name: "角色管理", Sequence: 30, MenuType: 2, Code: "Role", Icon: "tree", OperateType: "none"}
			db.Create(&menuModel)
			InitMenu(menuModel)
			menuModel = sys.Menu{Status: 1, ParentID: menuModel.ID, URL: "/role/setrole", Name: "分配角色菜单", Sequence: 6, MenuType: 3, Code: "RoleSetrolemenu", Icon: "", OperateType: "setrolemenu"}
			db.Create(&menuModel)
			menuModel = sys.Menu{Status: 1, ParentID: menuModelSys.ID, URL: "/admins", Name: "后台用户管理", Sequence: 40, MenuType: 2, Code: "Admins", Icon: "user", OperateType: "none"}
			db.Create(&menuModel)
			InitMenu(menuModel)
			menuModel = sys.Menu{Status: 1, ParentID: menuModel.ID, URL: "/admins/setrole", Name: "分配角色", Sequence: 6, MenuType: 3, Code: "AdminsSetrole", Icon: "", OperateType: "setadminrole"}
			db.Create(&menuModel)

			menuData, _ = getAllMenu()
		}
	} else {
		menuData, err = getMenusByAdminsid(userID)
		if err != nil {
			tools.ErrorJson(c, gconv.String(err))
			return
		}
	}
	var menus []MenuModel
	if len(menuData) > 0 {
		var topmenuid uint64 = menuData[0].ParentID
		if topmenuid == 0 {
			topmenuid = menuData[0].ID
		}
		menus = setMenu(menuData, topmenuid)
	}
	if len(menus) == 0 && userID == tools.SUPER_ADMIN_ID {
		menus = getSuperAdminMenu()
	}
	resData := UserData{Menus: menus, Name: "wufly"}
	resData.Avatar = "http://127.0.0.1:1689/resource/img/head_go.jpg"
	tools.SuccessJson(c, resData)
}

// 递归菜单
func setMenu(menus []sys.Menu, parentID uint64) (out []MenuModel) {
	var menuArr []sys.Menu
	linq.From(menus).Where(func(c interface{}) bool {
		return c.(sys.Menu).ParentID == parentID
	}).OrderBy(func(c interface{}) interface{} {
		return c.(sys.Menu).Sequence
	}).ToSlice(&menuArr)
	if len(menuArr) == 0 {
		return
	}
	noCache := false
	for _, item := range menuArr {
		menu := MenuModel{
			Path:      item.URL,
			Component: item.Code,
			Name:      item.Code,
			Meta:      MenuMeta{Title: item.Name, Icon: item.Icon, NoCache: noCache},
			Children:  []MenuModel{}}
		if item.MenuType == 3 {
			menu.Hidden = true
		}
		//查询是否有子级
		menuChildren := setMenu(menus, item.ID)
		if len(menuChildren) > 0 {
			menu.Children = menuChildren
		}
		if item.MenuType == 2 {
			// 添加子级首页，有这一级NoCache才有效
			menuIndex := MenuModel{
				Path:      "index",
				Component: item.Code,
				Name:      item.Code,
				Meta:      MenuMeta{Title: item.Name, Icon: item.Icon, NoCache: noCache},
				Children:  []MenuModel{}}
			menu.Children = append(menu.Children, menuIndex)
			menu.Name = menu.Name + "index"
			menu.Meta = MenuMeta{}
		}
		out = append(out, menu)
	}
	return
}

//查询所有菜单
func getAllMenu() (menus []sys.Menu, err error) {
	db.Find(&sys.Menu{}, &menus, "parent_id asc", "sequence asc")
	return
}

//查询登录用户权限菜单
func getMenusByAdminsid(adminsid uint64) (ret []sys.Menu, err error) {
	menu := sys.Menu{}
	var menus []sys.Menu
	err = menu.GetMenuByAdminsid(adminsid, &menus)
	if err != nil || len(menus) == 0 {
		return
	}
	allmenu, err := getAllMenu()
	if err != nil || len(allmenu) == 0 {
		return
	}
	menuMapAll := make(map[uint64]sys.Menu)
	for _, item := range allmenu {
		menuMapAll[item.ID] = item
	}
	menuMap := make(map[uint64]sys.Menu)
	for _, item := range menus {
		menuMap[item.ID] = item
	}
	for _, item := range menus {
		_, exists := menuMap[item.ParentID]
		if exists {
			continue
		}
		setMenuUp(menuMapAll, item.ParentID, menuMap)
	}
	for _, m := range menuMap {
		ret = append(ret, m)
	}
	linq.From(ret).OrderBy(func(c interface{}) interface{} {
		return c.(sys.Menu).ParentID
	}).ToSlice(&ret)
	return
}

//获取超级管理员初使菜单
func getSuperAdminMenu() (out []MenuModel) {
	menuTop := MenuModel{
		Path:      "/sys",
		Component: "Sys",
		Name:      "Sys",
		Meta:      MenuMeta{Title: "系统管理", NoCache: false},
		Children:  []MenuModel{}}
	menuModel := MenuModel{
		Path:      "/icon",
		Component: "Icon",
		Name:      "Icon",
		Meta:      MenuMeta{Title: "图标管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/menu",
		Component: "Menu",
		Name:      "Menu",
		Meta:      MenuMeta{Title: "菜单管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/role",
		Component: "Role",
		Name:      "Role",
		Meta:      MenuMeta{Title: "角色管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/admins",
		Component: "Admins",
		Name:      "Admins",
		Meta:      MenuMeta{Title: "用户管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	out = append(out, menuTop)
	return
}

// 向上查找父级菜单
func setMenuUp(menuMapAll map[uint64]sys.Menu, menuid uint64, menuMap map[uint64]sys.Menu) {
	menuModel, exists := menuMapAll[menuid]
	if exists {
		mid := menuModel.ID
		_, exists = menuMap[mid]
		if !exists {
			menuMap[mid] = menuModel
			setMenuUp(menuMapAll, menuModel.ParentID, menuMap)
		}
	}
}
