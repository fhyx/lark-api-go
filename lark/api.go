package lark

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"fhyx.online/lark-api-go/client"
)

type Error = client.Error

// consts
const (
	uriAPPTokenInnel    = "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal/"
	uriTenantTokenInnel = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/"

	uriAuthorize = "https://open.feishu.cn/open-apis/authen/v1/access_token"

	uriContactScope = "https://open.feishu.cn/open-apis/contact/v1/scope/get"

	uriUserAuthGet = "https://open.feishu.cn/open-apis/authen/v1/user_info"
	uriUserGetID   = "https://open.feishu.cn/open-apis/user/v1/batch_get_id"

	uriUserBatchGet   = "https://open.feishu.cn/open-apis/contact/v1/user/batch_get"
	uriUserListDetail = "https://open.feishu.cn/open-apis/contact/v1/department/user/detail/list"
	uriUserListSimp   = "https://open.feishu.cn/open-apis/contact/v1/department/user/list"
	uriUserBulk       = "https://open.feishu.cn/open-apis/contact/v2/user/batch_add"
	uriUesrAdd        = "https://open.feishu.cn/open-apis/contact/v1/user/add"

	uriTaskStatus = "https://open.feishu.cn/open-apis/contact/v2/task/get" // batch task status

	uriDeptSimpList = "https://open.feishu.cn/open-apis/contact/v1/department/simple/list"      // 获取子部门列表
	uriDeptBatchGet = "https://open.feishu.cn/open-apis/contact/v1/department/detail/batch_get" // 批量获取部门详情
	uriDeptSync     = "https://open.feishu.cn/open-apis/contact/v2/department/batch_add"        // 批量添加部门
)

var (
	_ IClient = (*API)(nil)
)

// API ...
type API struct {
	corpID     string
	corpSecret string
	ca         *client.Client
	ct         *client.Client
}

// NewAPI return new api instance with ([corpID, [corpSecret]])
func NewAPI(strs ...string) *API {
	corpID := os.Getenv("LARK_CORP_ID")
	corpSecret := os.Getenv("LARK_CORP_SECRET")
	if len(strs) > 0 && len(strs[0]) > 0 {
		corpID = strs[0]
		if len(strs) > 1 && len(strs[1]) > 0 {
			corpSecret = strs[1]
		}
	}

	if corpID == "" || corpSecret == "" {
		log.Printf("corpID or corpSecret are empty or not found")
	}
	ca := client.NewClient(uriAPPTokenInnel)
	ca.SetContentType("application/json")
	ca.SetCorp(corpID, corpSecret)
	ct := client.NewClient(uriTenantTokenInnel)
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

func uriForUserGet(uid string, at CType) string {
	switch at {
	case CEmail:
		return fmt.Sprintf("%s?emails=%s", uriUserGetID, uid)
	case CMobile:
		return fmt.Sprintf("%s?mobiles=%s", uriUserGetID, uid)
	default:
		return ""
	}
}

func (a *API) ListContactScope() (*AuthContactResponse, error) {
	cr := new(AuthContactResponse)
	err := a.ca.GetJSON(uriContactScope, cr)
	if err != nil {
		logger().Infow("list contact scope fail", "err", err)
		return nil, err
	}
	return cr, nil
}

// GetUser get user with emails,mobile
func (a *API) GetUser(uid string, at CType) (*User, error) {
	user := new(User)
	err := a.ca.GetJSON(uriForUserGet(uid, at), user)
	if err != nil {
		logger().Infow("get user fail", "at:", at, uid, "err", err)
		return nil, err
	}
	return user, nil
}

// ListUser ...
func (a *API) ListUser(lr ListReq) (ListResult, error) {
	if lr.Limit < 1 {
		lr.Limit = 1
	}
	uri := fmt.Sprintf("%s?department_id=%s&page_token=%s&page_size=%d", uriUserListDetail, lr.DeptID, lr.PageToken, lr.Limit)
	if lr.IncChild && lr.Limit > 1 {
		uri += "&fetch_child=true"
	}
	logger().Debugw("listUser", "req", lr)

	var ret = new(usersDetailResponse)
	err := a.ca.GetJSON(uri, ret)
	if err != nil {
		logger().Infow("getJSON fail", "uri", uri, "lr", lr, "err", err)
		return nil, err
	}

	return ret, err
}

// GetsDepartments ...
func (a *API) GetsDepartments(ids []string) (data Departments, err error) {
	var ret departmentResponse
	param := url.Values{"department_ids": ids}
	uri := fmt.Sprintf("%s?%s", uriDeptBatchGet, param.Encode())
	err = a.ca.GetJSON(uri, &ret)
	if err == nil {
		data = ret.Data.Departments
	}
	return
}

// ListDepartment child Department ...
func (a *API) ListDepartment(recursive bool, id string) (data Departments, err error) {
	var pageToken string
	limit := 20
	if "" == id {
		id = "0"
	}
queryF:
	uri := fmt.Sprintf("%s?department_id=%s&page_token=%s&page_size=%d", uriDeptSimpList, id, pageToken, limit)
	if recursive {
		uri += "&fetch_child=true"
	}

	var ret departmentResponse
	err = a.ca.GetJSON(uri, &ret)

	if err == nil {
		data = append(data, ret.Data.Departments...)
	}
	if ret.Data.HasMore && len(ret.Data.PageToken) > 0 {
		pageToken = ret.Data.PageToken
		logger().Infow("has more", "pageToken", pageToken)
		goto queryF
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
func (a *API) SyncUser(user UserUp) error {
	buf, err := json.Marshal(user)
	if err != nil {
		return err
	}
	var resp Error
	err = a.ca.PostJSON(uriUserBulk, buf, &resp)
	if err != nil {
		logger().Infow("sync User fail", "err", err)
		return err
	}
	logger().Infow("sync User ok", "resp", resp)
	return nil
}
