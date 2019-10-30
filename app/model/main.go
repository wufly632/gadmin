package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"review-order/app/model/db"
	"review-order/app/model/sys"
)

func InitDB(config string) {
	var gdb *gorm.DB
	var err error
	gdb, err = gorm.Open("mysql", config)
	if err != nil {
		panic(err)
	}
	gdb.SingularTable(true)
	//gdb.LogMode(true)
	//gdb.SetLogger(log.New(os.Stdout, "\r\n", 0))
	db.DB = gdb
}

func Migration() {
	fmt.Println(db.DB.AutoMigrate(new(sys.Menu)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.Admin)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.RoleMenu)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.Role)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.AdminRole)).Error)
}
