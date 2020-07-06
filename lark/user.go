package lark

import (
	// "time"
	"strings"

	"fhyx.online/lark-api-go/client"
	"fhyx.online/lark-api-go/gender"
)

// Status 状态
type Status uint

// 状态, 用户状态，bit0(最低位): 1冻结，0未冻结；bit1:1离职，0在职；bit2:1未激活，0已激活
const (
	SNone     Status = 0
	SFreezen  Status = 1
	SLeaved   Status = 2
	SInactive Status = 4
)

func (s Status) Labels() []string {
	as := []string{}
	if s&SFreezen > 0 {
		as = append(as, "freezen")
	}
	if s&SLeaved > 0 {
		as = append(as, "leveaved")
	}
	if s&SInactive > 0 {
		as = append(as, "inactived")
	}
	return as
}

func (s Status) String() string {
	return strings.Join(s.Labels(), ",")
}

// EType 员工类型
type EType uint8

// 员工类型。1:正式员工；2:实习生；3:外包；4:劳务；5:顾问
const (
	ETNone EType = iota
	ETNormal
	ETPractice
	ETOutSourcing
	ETContract
	ETAdvisor
)

//CType 查询类型
type CType uint8

const (
	CEmail CType = iota
	CMobile
)

func (z EType) String() string {
	switch z {
	case ETNormal:
		return "normal"
	case ETPractice:
		return "practice"
	case ETOutSourcing:
		return "outsourcing"
	case ETContract:
		return "contract"
	case ETAdvisor:
		return "advisoor"
	default:
		return "none"
	}
}

// User 用户
// "name":"zhang san",
// "name_py":"zhang san",
// "en_name":"John",
// "employee_id":"a0615a67",
// "employee_no":"235634",
// "open_id":"ou_e03053f0541cecc3269d7a9dc34a0b21",
// "status":2,
// "employee_type": 1,
// "avatar_72": "https://sf3-ttcdn-tos.pstatp.com/img/avatar/62db96e8-c5b6-4077-bb9d-2697d65a29eb~72x72.png",
// "avatar_240": "https://sf3-ttcdn-tos.pstatp.com/img/avatar/62db96e8-c5b6-4077-bb9d-2697d65a29eb~240x240.png",
// "avatar_640": "https://sf3-ttcdn-tos.pstatp.com/img/avatar/62db96e8-c5b6-4077-bb9d-2697d65a29eb~640x640.png",
// "avatar_url": "https://sf3-ttcdn-tos.pstatp.com/img/avatar/62db96e8-c5b6-4077-bb9d-2697d65a29eb~noop.png",
// "gender":1,
// "email":"zhangsan@gmail.com",
// "mobile":"+8615343215730",
// "description": "",
// "country": "CN",
// "city":"Beijing",
// "work_station":"Poly, F6-123",
// "is_tenant_manager":false,
// "join_time":1562342314,
// "update_time":1569140032,
// "leader_employee_id":"a0615a67",
// "leader_open_id":"ou_e03053f0541cecc3269d7a9dc34a0b21",
// "departments":[
//     "od-8c6c97ab9a34c1a649001d7ad36b97a7"
// ],
// "custom_attrs": {
//     "C-6702376000044400907": {
//         "value": "value1"
//     },
//     "C-6702376000048595214": {
//         "value": "value2"
//     }
// }
type User struct {
	Name             string        `json:"name"`                    // 用户名
	NameEN           string        `json:"en_name,omitempty"`       // 英文名
	NamePY           string        `json:"name_py,omitempty"`       // 用户名拼音
	EmployeeID       string        `json:"employee_id"`             // 用户的 employee_id，申请了"获取用户 user_id"权限后返回
	EmployeeNo       string        `json:"employee_no,omitempty"`   // 工号
	EmployeeType     EType         `json:"employee_type,omitempty"` // 员工类型。1:正式员工；2:实习生；3:外包；4:劳务；5:顾问
	AvatarURI        string        `json:"avatar_url,omitempty"`    // 头像，原始大小
	OpenID           string        `json:"open_id,omitempty"`       // 用户的 open_id
	UnionID          string        `json:"union_id,omitempty"`      // 用户的 union_id,申请了"获取用户统一ID"权限后返回
	Mobile           string        `json:"mobile,omitempty"`        // required
	Email            string        `json:"email,omitempty"`         // required
	Gender           gender.Gender `json:"gender,omitempty"`        // 性别
	Status           Status        `json:"userStatus,omitempty"`    // 用户状态，bit0(最低位): 1冻结，0未冻结；bit1:1离职，0在职；bit2:1未激活，0已激活
	Description      string        `json:"description,emitempty"`   // 用户个人签名
	Country          string        `json:"country,omitempty"`       // 用户所在国家
	City             string        `json:"city,omitempty"`          // 用户所在城市
	WorkStation      string        `json:"work_station,omitempty"`  // 工位
	JoinedStamp      int64         `json:"join_time,omitempty"`     // 入职时间
	UpdatedStamp     int64         `json:"update_time,omitempty"`   // 更新时间
	LeaderEmployeeID string        `json:"leader_employee_id,omitempty"`
	LeaderOpenID     string        `json:"leader_open_id,omitempty"`

	Departments []string `json:"departments,omitempty"` // 所在部门，用户可能同时存在于多个部门

	CustomAttrs map[string]interface{} `json:"custom_attrs,omitempty"`
}

type Users []User

// userListResp ...
type userListResp struct {
	HasMore   bool   `json:"hasMore"`
	PageToken string `json:"pageToken,omitempty"`
	Users     []User `json:"data"`
}

type usersSimpResponse struct {
	client.Error

	Data struct {
		HasMore   bool   `json:"has_more"`
		PageToken string `json:"page_token,omitempty"`
		Users     []User `json:"user_list"`
	} `json:"data"`
}

func (usr *usersSimpResponse) Users() Users {
	return usr.Data.Users
}

func (usr *usersSimpResponse) HasMore() bool {
	return usr.Data.HasMore
}

func (usr *usersSimpResponse) PageToken() string {
	return usr.Data.PageToken
}

type usersDetailResponse struct {
	client.Error

	Data struct {
		HasMore   bool   `json:"has_more"`
		PageToken string `json:"page_token,omitempty"`
		Users     []User `json:"user_infos"`
	} `json:"data"`
}

func (udr *usersDetailResponse) Users() Users {
	return udr.Data.Users
}

func (udr *usersDetailResponse) HasMore() bool {
	return udr.Data.HasMore
}

func (udr *usersDetailResponse) PageToken() string {
	return udr.Data.PageToken
}

type UserUp = User

// OAuth2UserInfo 为用户 OAuth2 验证登录后的简单信息
type OAuth2UserInfo struct {
	AccessToken      string `json:"access_token,omitempty"`
	AvatarURI        string `json:"avatar_url,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	Name             string `json:"name,omitempty"`
	NameEn           string `json:"en_name,omitempty"`
	OpenID           string `json:"open_id,omitempty"`
	TenantKey        string `json:"tenant_key,omitempty"`
	RefreshExpiresIn int64  `json:"refresh_expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
}

// OAuth2UserResp ...
type OAuth2UserResp struct {
	client.Error
	User *OAuth2UserInfo `json:"data,omitempty"`
}
