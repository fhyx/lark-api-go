package lark

// IClient ... interface of API client
type IClient interface {
	ListDepartment(recursive bool, ids ...string) (data Departments, err error)
	ListUser(deptID string, recursive bool) (data Users, err error)
	ListContactScope() (*AuthContactResponse, error)
}
