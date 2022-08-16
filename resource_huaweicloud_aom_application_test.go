package cmdb

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"testing"
)

func getAppResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	config := &config.Config{AccessKey: acceptance.HW_ACCESS_KEY, SecretKey: acceptance.HW_SECRET_KEY}
	c, diagErr := httpclient_go.NewHttpclientGo(config)
	if diagErr != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud client: %s", diagErr)
	}

	url := "https://aom.cn-north-7.myhuaweicloud.com/v1/applications/" + state.Primary.Attributes["id"]

	c.WithMethod(httpclient_go.MethodGet).WithUrl(url)
	response, err := c.Do()
	body, diagErr := c.CheckDeletedDiag(nil, err, response, "")
	if body == nil {
		return nil, fmt.Errorf("error creating HuaweiCloud client: %s", diagErr)
	}

	rlt := &entity.BizAppVo{}
	err = json.Unmarshal(body, rlt)

	if err == nil {
		return nil, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}

	return rlt, nil
}

func TestAccDcsInstances_basic(t *testing.T) {
	var instance entity.BizAppVo
	var instanceName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_application.app_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getAppResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesAomApp_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "description", "application description"),
					resource.TestCheckResourceAttr(resourceName, "displayName", instanceName),
					resource.TestCheckResourceAttr(resourceName, "eps_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "register_type", "CONSOLE"),
				),
			},
			{
				Config: tesAomApp_updated(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "port", "6389"),
					resource.TestCheckResourceAttr(resourceName, "description", "application description"),
				),
			},
		},
	})
}

func tesAomApp_basic(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_aom_application" "test" {
  name = "%s"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  description        = "application description"
  displayName        = huaweicloud_aom_application.app_1.name
  name               = huaweicloud_aom_application.app_1.name
  eps_id             = "0"
  register_type      = "CONSOLE"
}`, instanceName)
}

func tesAomApp_updated(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_aom_application" "test" {
  name = "%s"
}

resource "huaweicloud_aom_application" "app_1" {
  description        = "application description_update"
  displayName        = huaweicloud_aom_application.app_1.name
  name               = huaweicloud_aom_application.app_1.name
  eps_id             = "0"
  register_type      = "CONSOLE"
}`, instanceName)
}

func testAccDcsV1Instance_epsId(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name                  = "vpc-default"
  enterprise_project_id = "0"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"

  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1]
    save_days   = 1
  }
  enterprise_project_id = "%s"
}`, instanceName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDcsV1Instance_tiny(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 0.125
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = "redis.ha.xu1.tiny.r2.128"
  
  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1]
    save_days   = 1
  }
}`, instanceName)
}

func testAccDcsV1Instance_single(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = "redis.single.xu1.large.2"
}`, instanceName)
}

func testAccDcsV1Instance_whitelists(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 2
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = "redis.ha.xu1.large.r2.2"
  
  backup_policy {
    backup_type = "auto"
    begin_at    = "00:00-01:00"
    period_type = "weekly"
    backup_at   = [1]
    save_days   = 1
  }

  whitelists {
    group_name = "test-group1"
    ip_address = ["192.168.10.100", "192.168.0.0/24"]
  }
  whitelists {
    group_name = "test-group2"
    ip_address = ["172.16.10.100", "172.16.0.0/24"]
  }
}`, instanceName)
}
