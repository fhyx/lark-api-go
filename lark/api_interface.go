package lark

// IClient ... interface of API client
type IClient interface {
	ListDepartment(recursive bool, id string) (data Departments, err error)
	ListUser(deptID string, recursive bool) (data Users, err error)
	ListContactScope() (*AuthContactResponse, error)
	SyncDepartment(data []DepartmentUp) (res []DeptRespItem, err error)
}
