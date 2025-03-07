// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIbmSccProviderTypeDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmSccProviderTypeDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "provider_type_id"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "type"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "description"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "s2s_enabled"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "instance_limit"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "mode"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "data_type"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "icon"),
					resource.TestCheckResourceAttrSet("data.ibm_scc_provider_type.scc_provider_type_instance", "attributes.%"),
				),
			},
		},
	})
}

func testAccCheckIbmSccProviderTypeDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_scc_provider_type" "scc_provider_type_instance" {
			provider_type_id = "afa2476ecfa5f09af248492fe991b4d1"
		}
	`)
}
