package ic

import (
	"encoding/json"
	"fmt"
	"github.com/invisiblecloud/invisible-collector-go/models"
	"testing"
)

func TestGetCompany(t *testing.T) {
	//ic, _ := NewInvisibleCollector("1", IcAddress)
	//ic.GetCompany()
}

func TestJson(t *testing.T) {
	c := models.MakeCompany()
	c.SetAddress("ad")
	c.SetNotificationsEnabled(true)
	c.SetFieldToNull(models.CompanyCity)

	j, _ := json.Marshal(c)
	fmt.Println(string(j))

	json.Unmarshal(j, &c)
	j, _ = json.Marshal(c)
	fmt.Println(string(j))
}
