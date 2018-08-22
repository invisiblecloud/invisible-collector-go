package invisibleCollector

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetCompany(t *testing.T) {
	//ic, _ := NewInvisibleCollector("1", IcAddress)
	//ic.GetCompany()
}

func TestJson(t *testing.T) {
	c := MakeCompany()
	c.SetAddress("ad")
	c.SetNotificationsEnabled(true)
	c.SetFieldToNull(CompanyCity)

	j, _ := json.Marshal(c)
	fmt.Println(string(j))

	json.Unmarshal(j, &c)
	j, _ = json.Marshal(c)
	fmt.Println(string(j))
}
