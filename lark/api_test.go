package lark

import (
	"os"
	// "sort"
	"testing"

	"go.uber.org/zap"

	"github.com/fhyx/lark-api-go/log"
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

	data, err := api.ListDepartment(true)
	if err != nil {
		t.Fatal(err)
	}

	// sort.Sort(data)

	for _, dept := range data {
		t.Logf("dept %v", dept)
		users, err := api.ListUser(dept.ID, false)
		if err != nil {
			t.Fatal(err)
		}
		// t.Logf("users %v", users)
		for _, user := range users {
			t.Logf("user %v", user)
		}
	}

}

func TestUser(t *testing.T) {
	uid := os.Getenv("lark_TEST_UID")
	user, err := api.GetUser(uid, "uid")
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
