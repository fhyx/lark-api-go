package lark

// ListReq ...
type ListReq struct {
	DeptID    string `json:"deptID,omitempty"`
	Limit     int    `json:"limit"`
	PageToken string `json:"pageToken,omitempty"`
	IncChild  bool   `json:"incChild"`
}

// ListResult ...
type ListResult interface {
	HasMore() bool
	PageToken() string
	Users() Users
}

// IClient ... interface of API client
type IClient interface {
	ListDepartment(recursive bool, id string) (data Departments, err error)
	ListUser(r ListReq) (res ListResult, err error)
	ListContactScope() (*AuthContactResponse, error)
	SyncDepartment(data []DepartmentUp) (res []DeptRespItem, err error)
	SyncUser(user UserUp) error
}
