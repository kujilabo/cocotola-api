package domain

import (
	"fmt"

	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

func NewWorkbookWriter(workbookID WorkbookID) userD.RBACRole {
	return userD.RBACRole(fmt.Sprintf("workbook_%d_writer", uint(workbookID)))
}

func NewWorkbookReader(workbookID WorkbookID) userD.RBACRole {
	return userD.RBACRole(fmt.Sprintf("workbook_%d_reader", uint(workbookID)))
}

func NewWorkbookObject(workbookID WorkbookID) userD.RBACObject {
	return userD.RBACObject(fmt.Sprintf("workbook_%d", uint(workbookID)))
}

var WorkbookObjectPrefix = "workbook_"

var PrivilegeRead = userD.RBACAction("read")
var PrivilegeUpdate = userD.RBACAction("update")
var PrivilegeRemove = userD.RBACAction("remove")
