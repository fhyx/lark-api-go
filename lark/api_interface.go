package lark

// ListReq ...
type ListReq struct {
	DeptID    string   `json:"deptID,omitempty"`
	Limit     int      `json:"limit"`
	PageToken string   `json:"pageToken,omitempty"`
	IncChild  bool     `json:"incChild"`
	OpenIDs   []string `json:"openIDs"`
	IsSimple  bool     `json:"isSimple,omitempty"`
}

// ListResult ...
type ListResult interface {
	HasMore() bool
	PageToken() string
	Users() Users
}

// ContactScoper ...
type ContactScoper interface {
	GetDepartments() []string
	GetEmployeeIDs() []string
	GetOpenIDs() []string
}

// IClient ... interface of API client
type IClient interface {
	ListDepartment(recursive bool, id string) (data Departments, err error)
	ListUser(r ListReq) (res ListResult, err error)
	ListContactScope() (ContactScoper, error)
	SyncDepartment(data []DepartmentUp) (res []DeptRespItem, err error)
	SyncUser(user UserUp) error
}
