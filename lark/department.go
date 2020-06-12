package lark

import (
	"fhyx.online/lark-api-go/client"
)

// Department 部门
// "department_info": {
//     "id":"od-c042a4980ba8e1466050e3e8da2378fe",
//     "leader_employee_id":"612a67ef",
//     "leader_open_id":"ou_05065996251935ada9c2b0ecc50be91e",
//     "chat_id": "oc_405333f8fc89c3262865b014ccbbb274",
//     "member_count": 79,
//     "name": "市场部",
//     "parent_id": "0",
//     "status": 1
// }
type Department struct {
	ID               string `json:"id"`
	ParentID         string `json:"parent_id"`
	Name             string `json:"name"`
	NameEN           string `json:"name_en,omitempty"`
	ChatID           string `json:"chat_id"`
	LeaderEmployeeID string `json:"leader_employee_id,omitempty"`
	LeaderOpenID     string `json:"leader_open_id,omitempty"`
	Status           int    `json:"status,omitempty"`
	MemberCount      int    `json:"member_count,omitempty"`
}

type Departments []Department

// DepartmentUp 部门更新请求对象
// {
//     "name":"市场部",
//     "parent_id":"od-455efa262dc736b3e45a8b17fe945293",
//     "id":"tt_123456",
//     "leader_employee_id":"2fab234c",
//     "leader_open_id":"ou_4a2eb24a52b27c0b7fc6fd04162c0246",
//     "create_group_chat":true
// }

type DepartmentUp struct {
	ID              string `json:"id"`
	ParentID        string `json:"parent_id"`
	Name            string `json:"name"`
	LeaderUserID    string `json:"leader_user_id,omitempty"`
	LeaderOpenID    string `json:"leader_open_id,omitempty"`
	CreateGroupChat bool   `json:"create_group_chat,omitempty"`
}

type deptBatchReq struct {
	Data []DepartmentUp `json:"departments"`
}

type DeptRespItem struct {
	TaskID []string `json:"task_id"`
}

type deptBatchResp struct {
	client.Error

	Data []DeptRespItem `json:"data"`
}

type deptStatusUp struct {
	CorpDeptID int `json:"corpDeptCode,string"`
}

type deptStatusReq struct {
	Data []deptStatusUp `json:"deptInfo"`
}

// default sort
// func (z Departments) Len() int      { return len(z) }
// func (z Departments) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
// func (z Departments) Less(i, j int) bool {
// 	return z[i].ParentID == 0 || z[i].ParentID < z[j].ParentID ||
// 		z[i].Level < z[j].Level || z[i].ParentID == z[j].ParentID && z[i].OrderNo > z[j].OrderNo
// }

func (z Departments) WithID(id string) *Department {
	for _, dept := range z {
		if dept.ID == id {
			return &dept
		}
	}
	return nil
}

// departmentResponse
// "offset": 100,
// "limit": 25,
// "totalCount": 327,
// "departmentInfo": []
type departmentResponse struct {
	client.Error

	Data struct {
		HasMore     bool   `json:"has_more"`
		PageToken   string `json:"page_token"`
		Departments `json:"department_infos"`
	} `json:"data"`
}

// FilterDepartment Deprecated with Departments.WithID()
// func FilterDepartment(data []Department, id int) (*Department, error) {
// 	for _, dept := range data {
// 		if dept.ID == id {
// 			return &dept, nil
// 		}
// 	}
// 	return nil, ErrNotFound
// }
