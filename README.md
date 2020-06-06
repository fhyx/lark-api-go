# lark-api-go


## Environment for config

```
LARK_CORP_ID=AppId
LARK_CORP_SECRET=AppSecret

LARK_TEST_UID='for unit test only'
```

## Usage

```go

import "fhyx.online/lark-api-go/lark"


api := NewAPI() // or New(appId, appSecret)

deptId := 0
recursive := false
data, err := api.ListDepartment(deptId, recursive)

uid := "yourUID"
at := "uid" // uid,mobile,cuid
user, err := api.GetUser(uid, at)

```

## Links

* https://open.feishu.cn/document


## TODO

* Sync users
* Sync department
