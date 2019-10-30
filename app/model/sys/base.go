package sys

import (
	"fmt"
	"review-order/app/model/db"
)

func TableName(name string) string {
	return fmt.Sprintf("%s%s%s", db.GetTablePrefix(), "sys_", name)
}
