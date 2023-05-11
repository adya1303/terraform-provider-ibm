// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package metricsrouter_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/metricsrouterv3"
)

func TestAccIBMMetricsRouterTargetBasic(t *testing.T) {
	var conf metricsrouterv3.Target
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	destinationCRN := "crn:v1:bluemix:public:sysdig-monitor:us-south:a/11111111111111111111111111111111:22222222-2222-2222-2222-222222222222::"
	nameUpdate := fmt.Sprintf("updated_tf_name_%d", acctest.RandIntRange(10, 100))
	destinationCRNUpdate := "crn:v1:bluemix:public:sysdig-monitor:us-south:a/11111111111111111111111111111111:33333333-3333-3333-3333-333333333333::"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMMetricsRouterTargetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMMetricsRouterTargetConfigBasic(name, destinationCRN),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMMetricsRouterTargetExists("ibm_metrics_router_target.metrics_router_target_instance", conf),
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "name", name),
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "destination_crn", destinationCRN),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMMetricsRouterTargetConfigBasic(nameUpdate, destinationCRNUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "destination_crn", destinationCRNUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_metrics_router_target.metrics_router_target_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIBMMetricsRouterTargetAllArgs(t *testing.T) {
	var conf metricsrouterv3.Target
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	destinationCRN := "crn:v1:bluemix:public:sysdig-monitor:us-south:a/11111111111111111111111111111111:22222222-2222-2222-2222-222222222222::"
	region := "us-south"
	nameUpdate := fmt.Sprintf("updated_tf_name_%d", acctest.RandIntRange(10, 100))
	destinationCRNUpdate := "crn:v1:bluemix:public:sysdig-monitor:us-south:a/11111111111111111111111111111111:33333333-3333-3333-3333-333333333333::"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMMetricsRouterTargetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMMetricsRouterTargetConfig(name, destinationCRN, region),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMMetricsRouterTargetExists("ibm_metrics_router_target.metrics_router_target_instance", conf),
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "name", name),
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "destination_crn", destinationCRN),
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "region", region),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMMetricsRouterTargetConfig(nameUpdate, destinationCRNUpdate, region),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "destination_crn", destinationCRNUpdate),
				),
			},
		},
	})
}

func TestAccIBMMetricsRouterTargetNegative(t *testing.T) {
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	destinationCRN := "crn:v1:bluemix:public:sysdig-monitor:us-south:a/11111111111111111111111111111111:22222222-2222-2222-2222-222222222222::"
	region := "us-south"

	wrongRegion := "Frankfrut"           //A valid by regular expression but no such region exists
	wrongDestinationCRN := ":::::::::::" // "::::sysdig-monitor::::::" must be a sysdig-monitor instance

	dummyName := "$$dummy name"                           //Any random string that does not match the regular expresson ^[a-zA-Z0-9 \-._:]+$ and 1 ≤ length ≤ 1000
	dummyDestinationCRN := "$dummy destination CRN #*^()" //Any random string that does not match the regular expresson ^[a-zA-Z0-9 \-._:/]+$ and 3 ≤ length ≤ 1000
	dummyRegion := "$$dummy region"                       //Any random string that does not match the regular expresson ^[a-zA-Z0-9 \-._:]+$ and 3 ≤ length ≤ 1000

	errorMsg := "should match regexp"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMMetricsRouterTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckIBMMetricsRouterTargetConfig(dummyName, destinationCRN, region),
				ExpectError: regexp.MustCompile(errorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterTargetConfig(name, dummyDestinationCRN, region),
				ExpectError: regexp.MustCompile(errorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterTargetConfig(name, destinationCRN, dummyRegion),
				ExpectError: regexp.MustCompile(errorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterTargetConfig(name, destinationCRN, wrongRegion),
				ExpectError: regexp.MustCompile("CreateTargetWithContext failed Your request has failed because the value of `region` is not valid. Set a valid value in your target request and try again."),
			},
			{
				Config:      testAccCheckIBMMetricsRouterTargetConfig(name, wrongDestinationCRN, region),
				ExpectError: regexp.MustCompile("CreateTargetWithContext failed Your request has failed because the value of destination CRN is not valid. We only support"),
			},
		},
	})
}

func TestAccIBMMetricsRouterTargetNegativeUpdate(t *testing.T) {
	var conf metricsrouterv3.Target
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	destinationCRN := "crn:v1:bluemix:public:sysdig-monitor:us-south:a/11111111111111111111111111111111:22222222-2222-2222-2222-222222222222::"
	region := "us-south"

	wrongRegionUpdate := "Frankfrut"           //A valid by regular expression but no such region exists
	wrongDestinationCRNUpdate := ":::::::::::" // "::::sysdig-monitor::::::" must be a sysdig-monitor instance

	dummyNameUpdate := "$$dummy name"                           //Any random string that does not match the regular expresson ^[a-zA-Z0-9 \-._:]+$ and 1 ≤ length ≤ 1000
	dummyDestinationCRNUpdate := "$dummy destination CRN #*^()" //Any random string that does not match the regular expresson ^[a-zA-Z0-9 \-._:/]+$ and 3 ≤ length ≤ 1000
	dummyRegionUpdate := "$$dummy region"                       //Any random string that does not match the regular expresson ^[a-zA-Z0-9 \-._:]+$ and 3 ≤ length ≤ 1000

	errorMsg := "should match regexp"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMMetricsRouterTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMMetricsRouterTargetConfig(name, destinationCRN, region),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMMetricsRouterTargetExists("ibm_metrics_router_target.metrics_router_target_instance", conf),
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "name", name),
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "destination_crn", destinationCRN),
					resource.TestCheckResourceAttr("ibm_metrics_router_target.metrics_router_target_instance", "region", region),
				),
			},
			{
				Config:      testAccCheckIBMMetricsRouterTargetConfig(dummyNameUpdate, destinationCRN, region),
				ExpectError: regexp.MustCompile(errorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterTargetConfig(name, dummyDestinationCRNUpdate, region),
				ExpectError: regexp.MustCompile(errorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterTargetConfig(name, destinationCRN, dummyRegionUpdate),
				ExpectError: regexp.MustCompile(errorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterTargetConfig(name, destinationCRN, wrongRegionUpdate),
				ExpectError: regexp.MustCompile("CreateTargetWithContext failed Your request has failed because the value of `region` is not valid. Set a valid value in your target request and try again."),
			},
			{
				Config:      testAccCheckIBMMetricsRouterTargetConfig(name, wrongDestinationCRNUpdate, region),
				ExpectError: regexp.MustCompile("CreateTargetWithContext failed Your request has failed because the value of destination CRN is not valid. We only support"),
			},
		},
	})
}

func testAccCheckIBMMetricsRouterTargetConfigBasic(name string, destinationCRN string) string {
	return fmt.Sprintf(`

		resource "ibm_metrics_router_target" "metrics_router_target_instance" {
			name = "%s"
			destination_crn = "%s"
		}
	`, name, destinationCRN)
}

func testAccCheckIBMMetricsRouterTargetConfig(name string, destinationCRN string, region string) string {
	return fmt.Sprintf(`

		resource "ibm_metrics_router_target" "metrics_router_target_instance" {
			name = "%s"
			destination_crn = "%s"
			region = "%s"
		}
	`, name, destinationCRN, region)
}

func testAccCheckIBMMetricsRouterTargetExists(n string, obj metricsrouterv3.Target) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		metricsRouterClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).MetricsRouterV3()
		if err != nil {
			return err
		}

		getTargetOptions := &metricsrouterv3.GetTargetOptions{}

		getTargetOptions.SetID(rs.Primary.ID)

		target, _, err := metricsRouterClient.GetTarget(getTargetOptions)
		if err != nil {
			return err
		}

		obj = *target
		return nil
	}
}

func testAccCheckIBMMetricsRouterTargetDestroy(s *terraform.State) error {
	metricsRouterClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).MetricsRouterV3()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_metrics_router_target" {
			continue
		}

		getTargetOptions := &metricsrouterv3.GetTargetOptions{}

		getTargetOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := metricsRouterClient.GetTarget(getTargetOptions)

		if err == nil {
			return fmt.Errorf("metrics_router_target still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for metrics_router_target (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
