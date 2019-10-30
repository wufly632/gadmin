package boot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"review-order/app/model"
)

func init() {
	// 加载配置
	c := g.Config()
	// glog配置
	logpath := c.GetString("setting.logpath")
	glog.SetPath(logpath)
	glog.SetStdoutPrint(true)

	// 数据库配置
	config := c.GetString("database.link")
	initDB(config)

	g.Server().SetPort(1689)
}

func initDB(config string) {
	model.InitDB(config)
	model.Migration()
}
