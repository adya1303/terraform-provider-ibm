// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIbmSccLatestReportsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmSccLatestReportsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_scc_latest_reports.scc_latest_reports_instance", "id"),
				),
			},
		},
	})
}

func testAccCheckIbmSccLatestReportsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_scc_latest_reports" "scc_latest_reports_instance" {
			sort = "profile_name"
		}
	`)
}
