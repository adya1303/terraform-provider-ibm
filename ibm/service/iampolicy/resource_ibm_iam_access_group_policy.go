// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIAMAccessGroupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMIAMAccessGroupPolicyCreate,
		Read:   resourceIBMIAMAccessGroupPolicyRead,
		Update: resourceIBMIAMAccessGroupPolicyUpdate,
		Delete: resourceIBMIAMAccessGroupPolicyDelete,
		Exists: resourceIBMIAMAccessGroupPolicyExists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				resources, resourceAttributes, err := importAccessGroupPolicy(d, meta)
				if err != nil {
					return nil, fmt.Errorf("[ERROR] Error reading resource ID: %s", err)
				}
				d.Set("resources", resources)
				d.Set("resource_attributes", resourceAttributes)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of access group",
				ForceNew:    true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_access_group_policy",
					"access_group_id"),
			},

			"roles": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Role names of the policy definition",
			},

			"resources": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"account_management", "resource_attributes"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Service name of the policy definition",
						},

						"resource_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of resource instance of the policy definition",
						},

						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region of the policy definition",
						},

						"resource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource type of the policy definition",
						},

						"resource": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource of the policy definition",
						},

						"resource_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the resource group.",
						},

						"service_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Service type of the policy definition",
						},

						"service_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Service group id of the policy definition",
						},

						"attributes": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Set resource attributes in the form of 'name=value,name=value....",
							Elem:        schema.TypeString,
						},
					},
				},
			},

			"resource_attributes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Description:   "Set resource attributes.",
				ConflictsWith: []string{"resources", "account_management"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of attribute.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of attribute.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "stringEquals",
							Description: "Operator of attribute.",
						},
					},
				},
			},
			"account_management": {
				Type:          schema.TypeBool,
				Default:       false,
				Optional:      true,
				Description:   "Give access to all account management services",
				ConflictsWith: []string{"resources", "resource_attributes"},
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"resource_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set access management tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of attribute.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of attribute.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "stringEquals",
							Description: "Operator of attribute.",
						},
					},
				},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Policy",
			},

			"transaction_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Set transactionID for debug",
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"rule_conditions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Rule conditions enforced by the policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key of the condition",
						},
						"operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Operator of the condition",
						},
						"value": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Value of the condition",
						},
					},
				},
			},

			"rule_operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Operator that multiple rule conditions are evaluated over",
			},

			"pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Pattern rule follows for time-based condition",
			},
		},
	}
}

func ResourceIBMIAMAccessGroupPolicyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "access_group_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:access_group", "resolved_to:id"},
			Required:                   true})

	iBMIAMAccessGroupPolicyValidator := validate.ResourceValidator{ResourceName: "ibm_iam_access_group_policy", Schema: validateSchema}
	return &iBMIAMAccessGroupPolicyValidator
}

func resourceIBMIAMAccessGroupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	accessGroupId := d.Get("access_group_id").(string)

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	var policyID string
	ruleConditions, ruleConditionsDefined := d.GetOk("rule_conditions")
	pattern, patternDefined := d.GetOk("pattern")
	if ruleConditionsDefined || patternDefined {
		policyOptions, err := flex.GenerateV2PolicyOptions(d, meta)

		subjectAttribute := &iampolicymanagementv1.V2PolicySubjectAttribute{
			Key:      core.StringPtr("access_group_id"),
			Value:    &accessGroupId,
			Operator: core.StringPtr("stringEquals"),
		}

		policySubject := &iampolicymanagementv1.V2PolicySubject{
			Attributes: []iampolicymanagementv1.V2PolicySubjectAttribute{*subjectAttribute},
		}

		accountIDResourceAttribute := &iampolicymanagementv1.V2PolicyResourceAttribute{
			Key:      core.StringPtr("accountId"),
			Value:    core.StringPtr(userDetails.UserAccount),
			Operator: core.StringPtr("stringEquals"),
		}

		policyResource := &iampolicymanagementv1.V2PolicyResource{
			Attributes: append(policyOptions.Resource.Attributes, *accountIDResourceAttribute),
			Tags:       flex.SetV2PolicyTags(d),
		}

		createPolicyOptions := iamPolicyManagementClient.NewCreateV2PolicyOptions(
			policyOptions.Control,
			"access",
		)

		createPolicyOptions.SetSubject(policySubject)
		createPolicyOptions.SetResource(policyResource)

		if patternDefined {
			createPolicyOptions.SetPattern(pattern.(string))
		}

		if ruleConditionsDefined {
			createPolicyOptions.SetRule(flex.GeneratePolicyRule(d, ruleConditions))
		}

		if description, ok := d.GetOk("description"); ok {
			des := description.(string)
			createPolicyOptions.Description = &des
		}

		if transactionID, ok := d.GetOk("transaction_id"); ok {
			createPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		accessGroupPolicy, resp, err := iamPolicyManagementClient.CreateV2Policy(createPolicyOptions)

		if err != nil {
			return fmt.Errorf("Error creating user policies: %s, %s", err, resp)
		}

		policyID = *accessGroupPolicy.ID
	} else {
		policyOptions, err := flex.GeneratePolicyOptions(d, meta)
		if err != nil {
			return err
		}

		// Keep configuring the policy options by adding subject part
		accessGroupIdSubject := &iampolicymanagementv1.PolicySubject{
			Attributes: []iampolicymanagementv1.SubjectAttribute{
				{
					Name:  core.StringPtr("access_group_id"),
					Value: &accessGroupId,
				},
			},
		}

		accountIdResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
			Name:  core.StringPtr("accountId"),
			Value: &userDetails.UserAccount,
		}

		policyResource := &iampolicymanagementv1.PolicyResource{
			Attributes: append(policyOptions.Resources[0].Attributes, *accountIdResourceAttribute),
			Tags:       flex.SetTags(d),
		}

		createPolicyOptions := iamPolicyManagementClient.NewCreatePolicyOptions(
			"access",
			[]iampolicymanagementv1.PolicySubject{*accessGroupIdSubject},
			policyOptions.Roles,
			[]iampolicymanagementv1.PolicyResource{*policyResource},
		)

		if desc, ok := d.GetOk("description"); ok {
			des := desc.(string)
			createPolicyOptions.Description = &des
		}

		if transactionID, ok := d.GetOk("transaction_id"); ok {
			createPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		accessGroupPolicy, res, err := iamPolicyManagementClient.CreatePolicy(createPolicyOptions)
		if err != nil || accessGroupPolicy == nil {
			return fmt.Errorf("[ERROR] Error creating access group policy: %s\n%s", err, res)
		}

		policyID = *accessGroupPolicy.ID
	}

	getPolicyOptions := &iampolicymanagementv1.GetV2PolicyOptions{
		ID: &policyID,
	}

	if transactionID, ok := d.GetOk("transaction_id"); ok {
		getPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		policy, res, err := iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
		if err != nil || policy == nil {
			if res != nil && res.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		_, _, err = iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
	}
	if err != nil {
		d.SetId(fmt.Sprintf("%s/%s", accessGroupId, policyID))
		return fmt.Errorf("[ERROR] Error fetching access group policy: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", accessGroupId, policyID))

	return resourceIBMIAMAccessGroupPolicyRead(d, meta)
}

func resourceIBMIAMAccessGroupPolicyRead(d *schema.ResourceData, meta interface{}) error {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	accessGroupId := parts[0]
	accessGroupPolicyId := parts[1]

	getPolicyOptions := &iampolicymanagementv1.GetV2PolicyOptions{
		ID: &accessGroupPolicyId,
	}

	if transactionID, ok := d.GetOk("transaction_id"); ok {
		getPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
	}

	accessGroupPolicy := &iampolicymanagementv1.V2PolicyTemplateMetaData{}
	res := &core.DetailedResponse{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		accessGroupPolicy, res, err = iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
		if err != nil || accessGroupPolicy == nil {
			if res != nil && res.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		accessGroupPolicy, res, err = iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
	}
	if err != nil || accessGroupPolicy == nil || res == nil {
		return fmt.Errorf("[ERROR] Error retrieving access group policy: %s\n%s", err, res)
	}

	retrievedAttribute := flex.GetV2PolicySubjectAttribute("access_group_id", *accessGroupPolicy.Subject)
	if accessGroupId != *retrievedAttribute {
		return fmt.Errorf("[ERROR] Policy %s does not belong to access group %s, retrievedAttr: %s", accessGroupPolicyId, accessGroupId, *retrievedAttribute)
	}

	d.Set("access_group_id", accessGroupId)

	roles, err := flex.GetRoleNamesFromPolicyResponse(*accessGroupPolicy, d, meta)
	d.Set("roles", roles)
	d.Set("version", res.Headers.Get("ETag"))

	if _, ok := d.GetOk("resources"); ok {
		d.Set("resources", flex.FlattenV2PolicyResource(*accessGroupPolicy.Resource))
	}
	if _, ok := d.GetOk("resource_attributes"); ok {
		d.Set("resource_attributes", flex.FlattenV2PolicyResourceAttributes(accessGroupPolicy.Resource.Attributes))
	}

	if _, ok := d.GetOk("resource_tags"); ok {
		d.Set("resource_tags", flex.FlattenV2PolicyResourceTags(*accessGroupPolicy.Resource))
	}

	if _, ok := d.GetOk("rule_conditions"); ok {
		d.Set("rule_conditions", flex.FlattenRuleConditions(*accessGroupPolicy.Rule.(*iampolicymanagementv1.V2PolicyRule)))
	}

	if _, ok := d.GetOk("rule_operator"); ok {
		d.Set("rule_operator", *accessGroupPolicy.Rule.(*iampolicymanagementv1.V2PolicyRule).Operator)
	}

	if (&iampolicymanagementv1.V2PolicyResource{}) != accessGroupPolicy.Resource {
		if flex.GetV2PolicyResourceAttribute("serviceType", *accessGroupPolicy.Resource) == "service" {
			d.Set("account_management", false)
		}
		if flex.GetV2PolicyResourceAttribute("serviceType", *accessGroupPolicy.Resource) == "platform_service" {
			d.Set("account_management", true)
		}
	}

	if accessGroupPolicy.Description != nil {
		d.Set("description", *accessGroupPolicy.Description)
	}

	if len(res.Headers["Transaction-Id"]) > 0 && res.Headers["Transaction-Id"][0] != "" {
		d.Set("transaction_id", res.Headers["Transaction-Id"][0])
	}

	return nil
}

func resourceIBMIAMAccessGroupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}
	if d.HasChange("roles") || d.HasChange("resources") || d.HasChange("resource_attributes") || d.HasChange("account_management") || d.HasChange("description") || d.HasChange("resource_tags") || d.HasChange("rule_conditions") || d.HasChange("rule_operator") || d.HasChange("pattern") {
		parts, err := flex.IdParts(d.Id())
		if err != nil {
			return err
		}
		accessGroupId := parts[0]
		accessGroupPolicyId := parts[1]

		userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
		if err != nil {
			return err
		}

		getPolicyOptions := iamPolicyManagementClient.NewGetV2PolicyOptions(
			accessGroupPolicyId,
		)

		if transactionID, ok := d.GetOk("transaction_id"); ok {
			getPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		policy, response, err := iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
		if err != nil || policy == nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("[ERROR] Error retrieving Policy: %s\n%s", err, response)
		}

		accessGroupPolicyETag := response.Headers.Get("ETag")

		if strings.Contains(*policy.Href, "/v2/policies") {
			createPolicyOptions, err := flex.GenerateV2PolicyOptions(d, meta)
			if err != nil {
				return err
			}
			subjectAttribute := &iampolicymanagementv1.V2PolicySubjectAttribute{
				Key:      core.StringPtr("access_group_id"),
				Value:    &accessGroupId,
				Operator: core.StringPtr("stringEquals"),
			}

			policySubject := &iampolicymanagementv1.V2PolicySubject{
				Attributes: []iampolicymanagementv1.V2PolicySubjectAttribute{*subjectAttribute},
			}

			accountIDResourceAttribute := &iampolicymanagementv1.V2PolicyResourceAttribute{
				Key:      core.StringPtr("accountId"),
				Value:    core.StringPtr(userDetails.UserAccount),
				Operator: core.StringPtr("stringEquals"),
			}

			policyResource := &iampolicymanagementv1.V2PolicyResource{
				Attributes: append(createPolicyOptions.Resource.Attributes, *accountIDResourceAttribute),
				Tags:       flex.SetV2PolicyTags(d),
			}

			updatePolicyOptions := iamPolicyManagementClient.NewReplaceV2PolicyOptions(
				accessGroupPolicyId,
				accessGroupPolicyETag,
				createPolicyOptions.Control,
				"access",
			)
			updatePolicyOptions.SetSubject(policySubject)
			updatePolicyOptions.SetResource(policyResource)

			if pattern, ok := d.GetOk("pattern"); ok {
				updatePolicyOptions.SetPattern(pattern.(string))
			}

			if ruleConditions, ok := d.GetOk("rule_conditions"); ok {
				updatePolicyOptions.SetRule(flex.GeneratePolicyRule(d, ruleConditions))
			}

			if description, ok := d.GetOk("description"); ok {
				des := description.(string)
				updatePolicyOptions.Description = &des
			}

			if transactionID, ok := d.GetOk("transaction_id"); ok {
				updatePolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
			}

			_, resp, err := iamPolicyManagementClient.ReplaceV2Policy(updatePolicyOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error updating user policy: %s, %s", err, resp)
			}
		} else {

			_, ruleConditionsDefined := d.GetOk("rule_conditions")
			_, patternDefined := d.GetOk("pattern")
			if ruleConditionsDefined || patternDefined {
				return fmt.Errorf("Cannot use rule_conditions, rule_operator, or pattern when updating v1/policy. Delete existing v1/policy and create using rule_conditions and pattern.")
			}

			var policyOptions iampolicymanagementv1.CreatePolicyOptions
			policyOptions, err = flex.GeneratePolicyOptions(d, meta)
			if err != nil {
				return err
			}

			accessGroupIdSubject := &iampolicymanagementv1.PolicySubject{
				Attributes: []iampolicymanagementv1.SubjectAttribute{
					{
						Name:  core.StringPtr("access_group_id"),
						Value: &accessGroupId,
					},
				},
			}

			accountIdResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
				Name:  core.StringPtr("accountId"),
				Value: &userDetails.UserAccount,
			}

			policyResource := &iampolicymanagementv1.PolicyResource{
				Attributes: append(policyOptions.Resources[0].Attributes, *accountIdResourceAttribute),
				Tags:       flex.SetTags(d),
			}

			updatePolicyOptions := iamPolicyManagementClient.NewReplacePolicyOptions(
				accessGroupPolicyId,
				accessGroupPolicyETag,
				"access",
				[]iampolicymanagementv1.PolicySubject{*accessGroupIdSubject},
				policyOptions.Roles,
				[]iampolicymanagementv1.PolicyResource{*policyResource},
			)

			if desc, ok := d.GetOk("description"); ok {
				des := desc.(string)
				updatePolicyOptions.Description = &des
			}

			if transactionID, ok := d.GetOk("transaction_id"); ok {
				updatePolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
			}

			_, res, err := iamPolicyManagementClient.ReplacePolicy(updatePolicyOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error updating access group policy: %s\n%s", err, res)
			}
		}
	}

	return resourceIBMIAMAccessGroupPolicyRead(d, meta)
}

func resourceIBMIAMAccessGroupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	accessGroupPolicyId := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetV2PolicyOptions(
		accessGroupPolicyId,
	)

	// Get policy to find version in href
	policy, res, err := iamPolicyManagementClient.GetV2Policy(getPolicyOptions)

	if err != nil {
		return err
	}

	if strings.Contains(*policy.Href, "/v2/policies") {
		deletePolicyOptions := iamPolicyManagementClient.NewDeleteV2PolicyOptions(
			accessGroupPolicyId,
		)
		if transactionID, ok := d.GetOk("transaction_id"); ok {
			deletePolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		res, err = iamPolicyManagementClient.DeleteV2Policy(deletePolicyOptions)
	} else {
		deletePolicyOptions := iamPolicyManagementClient.NewDeletePolicyOptions(
			accessGroupPolicyId,
		)

		if transactionID, ok := d.GetOk("transaction_id"); ok {
			deletePolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		res, err = iamPolicyManagementClient.DeletePolicy(deletePolicyOptions)
	}

	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting access group policy: %s\n%s", err, res)
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMAccessGroupPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) < 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of accessGroupID/PolicyID", d.Id())
	}

	accessGroupPolicyId := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetV2PolicyOptions(
		accessGroupPolicyId,
	)

	accessGroupPolicy, resp, err := iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
	if err != nil || accessGroupPolicy == nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting access group policy: %s\n%s", err, resp)
	}

	if accessGroupPolicy != nil && accessGroupPolicy.State != nil && *accessGroupPolicy.State == "deleted" {
		return false, nil
	}

	tempID := fmt.Sprintf("%s/%s", *flex.GetV2PolicySubjectAttribute("access_group_id", *accessGroupPolicy.Subject), *accessGroupPolicy.ID)

	return tempID == d.Id(), nil
}

func importAccessGroupPolicy(d *schema.ResourceData, meta interface{}) (interface{}, interface{}, error) {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return nil, nil, err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil, nil, err
	}
	accgrpPolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetV2PolicyOptions(
		accgrpPolicyID,
	)

	accessGroupPolicy, res, err := iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR] Error retrieving access group policy: %s\n%s", err, res)
	}

	resources := flex.FlattenV2PolicyResource(*accessGroupPolicy.Resource)
	resource_attributes := flex.FlattenV2PolicyResourceAttributes(accessGroupPolicy.Resource.Attributes)
	d.Set("resource_tags", flex.FlattenV2PolicyResourceTags(*accessGroupPolicy.Resource))

	return resources, resource_attributes, nil
}
