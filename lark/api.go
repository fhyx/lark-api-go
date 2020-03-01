package lark

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fhyx/lark-api-go/client"
)

const (
	uriAPPToken    = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal/"
	uriTenantToken = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/"

	uriAuthorize = "https://open.feishu.cn/open-apis/authen/v1/access_token"

	uriUserGet        = "https://open.feishu.cn/open-apis/authen/v1/user_info"
	uriUserGetID      = "https://open.feishu.cn/open-apis/user/v1/batch_get_id"
	uriUserBatchGet   = "https://open.feishu.cn/open-apis/contact/v1/user/batch_get"
	uriUserListDetail = "https://open.feishu.cn/open-apis/contact/v1/department/user/detail/list"
	uriUserListSimp   = "https://open.feishu.cn/open-apis/contact/v1/department/user/list"
	uriUserBulk       = "https://open.feishu.cn/open-apis/contact/v2/user/batch_add"
	uriUesrAdd        = "https://open.feishu.cn/open-apis/contact/v1/user/add"

	uriTaskStatus = "https://open.feishu.cn/open-apis/contact/v2/task/get" // batch task status

	uriDeptList = "https://open.feishu.cn/open-apis/contact/v1/department/simple/list"
	uriDeptSync = "https://open.feishu.cn/open-apis/contact/v2/department/batch_add" // 批量添加部门
)

// API ...
type API struct {
	corpID     string
	corpSecret string
	ca         *client.Client
	ct         *client.Client
}

func NewAPI() *API {
	return New(os.Getenv("LARK_CORP_ID"), os.Getenv("LARK_CORP_SECRET"))
}

// New ...
func New(corpID, corpSecret string) *API {
	if corpID == "" || corpSecret == "" {
		log.Printf("corpID or corpSecret are empty or not found")
	}
	ca := client.NewClient(uriAPPToken)
	ca.SetContentType("application/json")
	ca.SetCorp(corpID, corpSecret)
	ct := client.NewClient(uriTenantToken)
	ct.SetContentType("application/json")
	ct.SetCorp(corpID, corpSecret)
	return &API{
		corpID:     corpID,
		corpSecret: corpSecret,
		ca:         ca,
		ct:         ct,
	}
}

func (a *API) CorpID() string {
	return a.corpID
}

type authCodeReq struct {
	AccessToken string `json:"app_access_token"`
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
}

// AuthorizeCode ...
func (a *API) AuthorizeCode(code string) (ou *OAuth2UserInfo, err error) {
	var token string
	token, err = a.ca.GetAuthToken()
	if err != nil {
		return
	}

	var req = authCodeReq{AccessToken: token, GrantType: "authorization_code", Code: code}
	var buf []byte
	buf, err = json.Marshal(&req)
	if err != nil {
		logger().Infow("unmarshal fail", "err", err)
		return
	}

	var resp = new(OAuth2UserResp)
	err = a.ca.PostJSON(uriAuthorize, buf, resp)
	if err != nil {
		logger().Infow("authen fail", "resp", resp, "err", err)
		return
	}
	ou = resp.User
	return
}

func uriForUserGet(uid, at string) string {

	switch at {
	// case "uid":
	// 	return fmt.Sprintf("%s?userId=%s", uriUserGetID, uid)
	case "email":
		return fmt.Sprintf("%s?emails=%s", uriUserGetID, uid)
	case "mobile":
		return fmt.Sprintf("%s?mobiles=%s", uriUserGetID, uid)
	default:
		return fmt.Sprintf("%s?user_id=%s", uriUserGetID, uid)
	}
}

// GetUser get user with uid,mobile,cuid
func (a *API) GetUser(uid, at string) (*User, error) {
	user := new(User)
	err := a.ca.GetJSON(uriForUserGet(uid, at), user)
	if err != nil {
		logger().Infow("get user fail", "at:"+at, uid, "err", err)
		return nil, err
	}
	return user, nil
}

// ListUser ...
func (a *API) ListUser(deptID string, recursive bool) (data Users, err error) {
	offset := 0
	limit := 50
	uri := fmt.Sprintf("%s?department_id=%s&offset=%d&page_size=%d", uriUserListDetail, deptID, offset, limit)
	if recursive {
		uri += "&fetch_child=true"
	}

	var ret usersDetailResponse
	err = a.ca.GetJSON(uri, &ret)

	if err == nil {
		data = ret.Data.Users
	}

	return
}

// ListDepartment ...
func (a *API) ListDepartment(id string, recursive bool) (data Departments, err error) {

	offset := 0
	limit := 20
	uri := fmt.Sprintf("%s?department_id=%s&offset=%d&page_size=%d", uriDeptList, id, offset, limit)
	if recursive {
		uri += "&fetch_child=true"
	}

	var ret departmentResponse
	err = a.ca.GetJSON(uri, &ret)

	if err == nil {
		data = ret.Data.Departments
	}

	// if recursive && id == "0" {
	// 	for _, dept := range data {
	// 		var child Departments
	// 		child, err = a.ListDepartment(dept.ID, true)
	// 		if err != nil {
	// 			return
	// 		}
	// 		data = append(data, child...)
	// 	}
	// }

	return
}

// SyncDepartment ...
func (a *API) SyncDepartment(data []DepartmentUp) (res []DeptRespItem, err error) {
	var req deptBatchReq
	req.Data = data

	var buf []byte
	buf, err = json.Marshal(&req)
	if err != nil {
		return
	}
	var resp deptBatchResp
	err = a.ca.PostJSON(uriDeptSync, buf, &resp)
	if err != nil {
		logger().Infow("sync department fail", "err", err)
		return
	}
	res = resp.Data
	logger().Infow("sync department ok", "resp", resp)
	return
}

// GetTaskStatus ...
func (a *API) GetTaskStatus(taskID string) (res []DeptRespItem, err error) {

	uri := uriTaskStatus + "?task_id=" + taskID
	var resp deptBatchResp
	err = a.ca.GetJSON(uri, &resp)
	if err != nil {
		logger().Infow("status department fail", "err", err)
		return
	}
	res = resp.Data
	logger().Infow("status department ok", "resp", resp)
	return
}

// SyncUser ...
func (a *API) SyncUser(data []UserUp) (res []UserRespItem, err error) {
	var req userBatchReq
	req.Data = data

	var buf []byte
	buf, err = json.Marshal(&req)
	if err != nil {
		return
	}
	var resp userBatchResp
	err = a.ca.PostJSON(uriUserBulk, buf, &resp)
	if err != nil {
		logger().Infow("sync User fail", "err", err)
		return
	}
	res = resp.Data
	logger().Infow("sync User ok", "resp", resp)
	return
}
