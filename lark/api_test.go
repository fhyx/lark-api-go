package lark

import (
	"os"
	// "sort"
	"testing"

	"go.uber.org/zap"

	"fhyx.online/lark-api-go/log"
)

var (
	api *API
)

func TestMain(m *testing.M) {
	_logger, _ := zap.NewDevelopment()
	defer _logger.Sync() // flushes buffer, if any
	sugar := _logger.Sugar()
	log.SetLogger(sugar)

	api = NewAPI()
	os.Exit(m.Run())
}

// TestAPIDepartment test api // lark_CORP_ID= lark_CORP_SECRET=
func TestAPIDepartment(t *testing.T) {

	data, err := api.ListDepartment(true, "")
	if err != nil {
		t.Fatal(err)
	}

	// sort.Sort(data)

	for _, dept := range data {
		t.Logf("dept %v", dept)
		data, err := api.ListUser(ListReq{DeptID: dept.ID})
		if err != nil {
			t.Fatal(err)
		}
		// t.Logf("users %v", users)
		for _, user := range data.Users() {
			t.Logf("user %v", user)
		}
	}

}

func TestUser(t *testing.T) {
	uid := os.Getenv("lark_TEST_EMAIL")
	user, err := api.GetUser(uid, CEmail)
	if err != nil {
		t.Fatal(err)
	}
	logger().Infow("got", "user", user)
}

func TestContactScope(t *testing.T) {
	cr, err := api.ListContactScope()
	if err != nil {
		t.Fatal(err)
	}
	logger().Infow("get", "contact scope", cr)
}
