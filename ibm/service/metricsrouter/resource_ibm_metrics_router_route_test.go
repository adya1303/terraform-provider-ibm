// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package metricsrouter_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/metricsrouterv3"
)

func TestAccIBMMetricsRouterRouteBasic(t *testing.T) {
	var conf metricsrouterv3.Route
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	filterValue := "value"
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	updatedFilterValue := "value"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMMetricsRouterRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMMetricsRouterRouteConfigBasic(name, filterValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMMetricsRouterRouteExists("ibm_metrics_router_route.metrics_router_route_instance", conf),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "name", name),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operand", "location"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operator", "is"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.value.0", filterValue),
				),
			},
			{
				Config: testAccCheckIBMMetricsRouterRouteConfigBasic(nameUpdate, updatedFilterValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operand", "location"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operator", "is"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.value.0", updatedFilterValue),
				),
			},
			{
				ResourceName:      "ibm_metrics_router_route.metrics_router_route_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIBMMetricsRouterRouteAllArgs(t *testing.T) {
	var conf metricsrouterv3.Route
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	operand := "location"
	operator := "is"
	value := []string{"us-east"}

	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	operandUpdate := "resource"
	operatorUpdate := "in"
	valueUpdate := []string{"resource-1", "resource-2"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMMetricsRouterRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMMetricsRouterRouteConfigAllArgs("metrics_router_target_instance", "us-south", name, operand, operator, value),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMMetricsRouterRouteExists("ibm_metrics_router_route.metrics_router_route_instance", conf),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "name", name),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operand", operand),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operator", operator),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.value.0", value[0]),
					resource.TestCheckResourceAttrPair("ibm_metrics_router_target.metrics_router_target_instance", "id", "ibm_metrics_router_route.metrics_router_route_instance", "rules.0.target_ids.0"),
				),
			},
			{
				Config: testAccCheckIBMMetricsRouterRouteConfigAllArgs("metrics_router_target_instance", "us-south", nameUpdate, operandUpdate, operatorUpdate, valueUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operand", operandUpdate),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operator", operatorUpdate),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.value.0", valueUpdate[0]),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.value.1", valueUpdate[1]),
					resource.TestCheckResourceAttrPair("ibm_metrics_router_target.metrics_router_target_instance", "id", "ibm_metrics_router_route.metrics_router_route_instance", "rules.0.target_ids.0"),
				),
			},
			{
				ResourceName:      "ibm_metrics_router_route.metrics_router_route_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIBMMetricsRouterRouteNegative(t *testing.T) {
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	operand := "location"
	operator := "is"
	value := []string{"us-east"}

	dummyName := "$$dummy name"      //Any random string that does not match the regular expresson ^[a-zA-Z0-9 \-._:]+$ and 1 ≤ length ≤ 1000
	unallowedOperator := "srvc_name" //Any string other than [location, service_name, service_instance, resource_type, resource]
	unallowedOperand := "at"         //Any string other than ["is", "in"]

	regexpErrorMsg := "should match regexp"
	allowedErrorMsg := "must contain a value from"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMMetricsRouterRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckIBMMetricsRouterRouteConfigAllArgs("mr-target", "us-south", dummyName, operand, operator, value),
				ExpectError: regexp.MustCompile(regexpErrorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterRouteConfigAllArgs("mr-target", "us-south", name, unallowedOperand, operator, value),
				ExpectError: regexp.MustCompile(allowedErrorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterRouteConfigAllArgs("mr-target", "us-south", name, operand, unallowedOperator, value),
				ExpectError: regexp.MustCompile(allowedErrorMsg),
			},
		},
	})
}

func TestAccIBMMetricsRouterRouteNegativeUpdate(t *testing.T) {
	var conf metricsrouterv3.Route
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	operand := "location"
	operator := "is"
	value := []string{"us-east"}

	dummyName := "$$dummy name"      //Any random string that does not match the regular expresson ^[a-zA-Z0-9 \-._:]+$ and 1 ≤ length ≤ 1000
	unallowedOperator := "srvc_name" //Any string other than [location, service_name, service_instance, resource_type, resource]
	unallowedOperand := "at"         //Any string other than ["is", "in"]

	regexpErrorMsg := "should match regexp"
	allowedErrorMsg := "must contain a value from"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMMetricsRouterRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMMetricsRouterRouteConfigAllArgs("mr-target", "us-south", name, operand, operator, value),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMMetricsRouterRouteExists("ibm_metrics_router_route.metrics_router_route_instance", conf),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "name", name),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operand", operand),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operator", operator),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.value.0", value[0]),
				),
			},
			{
				Config:      testAccCheckIBMMetricsRouterRouteConfigAllArgs("mr-target", "us-south", dummyName, operand, operator, value),
				ExpectError: regexp.MustCompile(regexpErrorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterRouteConfigAllArgs("mr-target", "us-south", name, unallowedOperand, operator, value),
				ExpectError: regexp.MustCompile(allowedErrorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterRouteConfigAllArgs("mr-target", "us-south", name, operand, unallowedOperator, value),
				ExpectError: regexp.MustCompile(allowedErrorMsg),
			},
			{
				Config:      testAccCheckIBMMetricsRouterRouteConfigAllArgs("mr-target", "us-east", name, operand, operator, value),
				ExpectError: regexp.MustCompile(`DeleteTargetWithContext failed Your request has failed because the target id [a-zA-Z0-9 \-._:]+ is being used by a route.`),
			},
		},
	})
}

func TestAccIBMMetricsRouterRouteRuleAND(t *testing.T) {
	var conf metricsrouterv3.Route
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMMetricsRouterRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMMetricsRouterRouteConfigRuleAND(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMMetricsRouterRouteExists("ibm_metrics_router_route.metrics_router_route_instance", conf),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operand", "location"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operator", "is"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.value.0", "us-south"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.1.operand", "resource"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.1.operator", "in"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.1.value.0", "resource-1"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.1.value.1", "resource-2"),
				),
			},
		},
	})
}

func TestAccIBMMetricsRouterRouteRuleOR(t *testing.T) {
	var conf metricsrouterv3.Route
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMMetricsRouterRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMMetricsRouterRouteConfigRuleOR(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMMetricsRouterRouteExists("ibm_metrics_router_route.metrics_router_route_instance", conf),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operand", "service_name"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.operator", "is"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.0.inclusion_filters.0.value.0", "kubernetes"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.1.inclusion_filters.0.operand", "resource_type"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.1.inclusion_filters.0.operator", "is"),
					resource.TestCheckResourceAttr("ibm_metrics_router_route.metrics_router_route_instance", "rules.1.inclusion_filters.0.value.0", "worker"),
				),
			},
		},
	})
}

func testAccCheckIBMMetricsRouterRouteConfigBasic(name, filter_value string) string {
	return fmt.Sprintf(`
		resource "ibm_metrics_router_target" "metrics_router_target_instance" {
			name = "my-mr-target"
			destination_crn = "crn:v1:bluemix:public:sysdig-monitor:us-south:a/0be5ad401ae913d8ff665d92680664ed:22222222-2222-2222-2222-222222222222::"
		}

		resource "ibm_metrics_router_route" "metrics_router_route_instance" {
			name = "%s"
			rules {
				target_ids = [ ibm_metrics_router_target.metrics_router_target_instance.id ]
				inclusion_filters {
					operand = "location"
					operator = "is"
					value = [ "%s" ]
				}
			}
		}
	`, name, filter_value)
}

func testAccCheckIBMMetricsRouterRouteConfigAllArgs(targetName, targetRegion, routeName, operand, operator string, value []string) string {
	return fmt.Sprintf(`
		resource "ibm_metrics_router_target" "metrics_router_target_instance" {
			name = "%s"
			destination_crn = "crn:v1:bluemix:public:sysdig-monitor:us-south:a/0be5ad401ae913d8ff665d92680664ed:22222222-2222-2222-2222-222222222222::"
			region = "%s"
		}

		resource "ibm_metrics_router_route" "metrics_router_route_instance" {
			name = "%s"
			rules {
				target_ids = [ ibm_metrics_router_target.metrics_router_target_instance.id ]
				inclusion_filters {
					operand = "%s"
					operator = "%s"
					value = ["%s"]
				}
			}
		}
	`, targetName, targetRegion, routeName, operand, operator, strings.Join(value, "\", \""))
}

func testAccCheckIBMMetricsRouterRouteConfigRuleAND(name string) string {
	return fmt.Sprintf(`
	variable "inclusion_filters" {
		type = list
		default = [
			{"operand"="location", "operator"="is", "value"=["us-south"]},
		  	{"operand"="resource", "operator"="in", "value"=["resource-1", "resource-2"]}
		]
	}

	resource "ibm_metrics_router_target" "metrics_router_target_instance1" {
		name = "my-mr-target"
		destination_crn = "crn:v1:bluemix:public:sysdig-monitor:us-south:a/0be5ad401ae913d8ff665d92680664ed:22222222-2222-2222-2222-222222222222::"
		region = "us-south"
	}

	resource "ibm_metrics_router_target" "metrics_router_target_instance2" {
		name = "my-mr-target"
		destination_crn = "crn:v1:bluemix:public:sysdig-monitor:us-south:a/0be5ad401ae913d8ff665d92680664ed:22222222-2222-2222-2222-222222222222::"
		region = "us-east"
	}

	resource "ibm_metrics_router_route" "metrics_router_route_instance" {
		name = "%s"
		rules {
			target_ids = [ibm_metrics_router_target.metrics_router_target_instance1.id,  ibm_metrics_router_target.metrics_router_target_instance2.id]
			dynamic "inclusion_filters" {
				for_each = var.inclusion_filters
		        content{
                    operand = inclusion_filters.value.operand
                	operator = inclusion_filters.value.operator
             		value = inclusion_filters.value.value
             	}
			}
		}
	}`, name)
}

func testAccCheckIBMMetricsRouterRouteConfigRuleOR(name string) string {
	return fmt.Sprintf(`
	resource "ibm_metrics_router_target" "metrics_router_target_instance" {
		name = "my-mr-target"
		destination_crn = "crn:v1:bluemix:public:sysdig-monitor:us-south:a/0be5ad401ae913d8ff665d92680664ed:22222222-2222-2222-2222-222222222222::"
		region = "us-south"
	}

	resource "ibm_metrics_router_route" "metrics_router_route_instance" {
		name = "%s"
		rules {
			target_ids = [ibm_metrics_router_target.metrics_router_target_instance.id]
			inclusion_filters {
				operand = "service_name"
				operator = "is"
				value = ["kubernetes"]
			}
		}
		rules {
			target_ids = [ibm_metrics_router_target.metrics_router_target_instance.id]
			inclusion_filters {
				operand = "resource_type"
				operator = "is"
				value = ["worker"]
			}
		}
	}`, name)
}

func testAccCheckIBMMetricsRouterRouteExists(n string, obj metricsrouterv3.Route) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		metricsRouterClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).MetricsRouterV3()
		if err != nil {
			return err
		}

		getRouteOptions := &metricsrouterv3.GetRouteOptions{}

		getRouteOptions.SetID(rs.Primary.ID)

		route, _, err := metricsRouterClient.GetRoute(getRouteOptions)
		if err != nil {
			return err
		}

		obj = *route
		return nil
	}
}

func testAccCheckIBMMetricsRouterRouteDestroy(s *terraform.State) error {
	metricsRouterClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).MetricsRouterV3()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_metrics_router_route" {
			continue
		}

		getRouteOptions := &metricsrouterv3.GetRouteOptions{}

		getRouteOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := metricsRouterClient.GetRoute(getRouteOptions)

		if err == nil {
			return fmt.Errorf("metrics_router_route still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for metrics_router_route (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
