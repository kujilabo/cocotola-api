package domain

import (
	"fmt"

	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

func NewWorkbookWriter(workbookID WorkbookID) user.RBACRole {
	return user.RBACRole(fmt.Sprintf("workbook_%d_writer", uint(workbookID)))
}

func NewWorkbookReader(workbookID WorkbookID) user.RBACRole {
	return user.RBACRole(fmt.Sprintf("workbook_%d_reader", uint(workbookID)))
}

func NewWorkbookObject(workbookID WorkbookID) user.RBACObject {
	return user.RBACObject(fmt.Sprintf("workbook_%d", uint(workbookID)))
}

var WorkbookObjectPrefix = "workbook_"

var PrivilegeRead = user.RBACAction("read")
var PrivilegeUpdate = user.RBACAction("update")
var PrivilegeRemove = user.RBACAction("remove")
