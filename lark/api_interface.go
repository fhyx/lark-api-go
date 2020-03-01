package lark

// Client ... interface of API client
type Client interface {
	ListDepartment(id string, recursive bool) (data Departments, err error)
	ListUser(deptID string, recursive bool) (data Users, err error)
}
