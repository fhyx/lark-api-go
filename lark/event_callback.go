package lark

import (
	"encoding/json"
)

// CallbackReq ...
// {
//     "challenge": "ajls384kdjx98XX", // 应用需要原样返回的值
//     "token": "xxxxxx",              // 约定的校验Token
//     "type": "url_verification"      // 表示这是一个验证请求
// }
type CallbackReq struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

// UnmarsalCallback ...
func UnmarsalCallback(s string) (cr *CallbackReq, err error) {
	cr = new(CallbackReq)
	err = json.Unmarshal([]byte(s), cr)
	return
}

// CallbackResp ...
type CallbackResp struct {
	Challenge string `json:"challenge"`
}

// EventCallback ...
// {
//      "uuid":"c4ca4238a0b923820dcc509a6f75849b",
//      "token": "41a9425ea7df4536a7623e38fa321bae", //校验Token
//      "ts": "1502199207.7171419", //时间戳
//      "type": "event_callback",   //事件回调此处固定为event_callback
//      "event": json
// }
type EventCallback struct {
	UUID  string          `json:"uuid"`
	Token string          `json:"token"`
	Stamp int64           `json:"ts"`
	Type  string          `json:"type"`
	Event json.RawMessage `json:"event"`
}

// Event ...
type Event struct {
	Type      string `json:"type"`
	APPID     string `json:"app_id,omitempty"`
	TenantKey string `json:"tenant_key,omitempty"`
}

// EventCallbackForUser 通讯录用户相关变更 {
//      "type": "user_add",    //事件类型，包括user_add,user_update,user_leave
//      "app_id": "cli_xxx",
//      "tenant_key": "xxx",  //企业标识
//      "open_id":"xxx" ,
//      "employee_id":"xxx"   //企业自建应用返回
// }
type EventCallbackForUser struct {
	Type       string `json:"type"`
	APPID      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
	OpenID     string `json:"open_id"`
	EmployeeID string `json:"employee_id"`
}

// EventCallbackForDepartment 通讯录部门相关变更 {
//     "type": "dept_add",  //事件类型，包括 dept_add,dept_update,dept_delete
//     "app_id": "cli_xxx",
//     "tenant_key": "xxx",           //企业标识
//     "open_department_id":"od-xxx"  //部门的Id
// }
type EventCallbackForDepartment struct {
	Type             string `json:"type"`
	APPID            string `json:"app_id"`
	TenantKey        string `json:"tenant_key"`
	OpenDepartmentID string `json:"open_department_id"`
}

// EventCallbackForUserStatus 用户状态变更 {
// 	"app_id": "cli_9c8609450f78d102",
// 	"before_status": {
// 		"is_active": false,        // 账号是否已激活
// 		"is_frozen": false,       // 账号是否冻结
// 		"is_resigned": false    // 是否离职
// 	},
// 	"change_time": "2020-02-21 16:28:48",
// 	"current_status": {
// 		"is_active": true,
// 		"is_frozen": false,
// 		"is_resigned": false
// 	},
// 	"employee_id": "ca51d83b",
// 	"open_id": "ou_2d2c0399b53d06fd195bb393cd1e38f2",
// 	"tenant_key": "xxx",
// 	"type": "user_status_change"
// }
type EventCallbackForUserStatus struct {
	EventCallbackForUser
	BeforeStatus struct {
		IsActive bool `json:"is_active"`
		IsFrozen bool `json:"is_frozen"`
		IsLeaved bool `json:"is_resigned"`
	} `json:"before_status"`
	ChangedStatus struct {
		IsActive bool `json:"is_active"`
		IsFrozen bool `json:"is_frozen"`
		IsLeaved bool `json:"is_resigned"`
	} `json:"current_status"`
	ChangeTime string `json:"change_time"`
}
