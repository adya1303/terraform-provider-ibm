---
layout: "ibm"
page_title: "IBM : ibm_scc_profile_attachment"
description: |-
  Manages scc_profile_attachment.
subcategory: "Security and Compliance Center"
---

# ibm_scc_profile_attachment

Create, update, and delete profile attachments with this resource.

## Example Usage

```hcl
resource "ibm_scc_profile_attachment" "scc_profile_attachment_instance" {
  profiles_id = ibm_scc_profile.scc_profile_instance.id
}
```

## Argument Reference

You can specify the following arguments for this resource.

* `profile_id` - (Required, Forces new resource, String) The profile ID.
  * Constraints: The maximum length is `36` characters. The minimum length is `36` characters. The value must match regular expression `/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/`.
* `scope` - (List) The scope payload for the multi cloud feature.
  * Constraints: The maximum length is `8` items. The minimum length is `0` items.
Nested schema for **scope**:
	* `environment` - (String) The environment that relates to this scope.
	  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters. The value must match regular expression `/[A-Za-z0-9]+/`.
	* `properties` - (List) The properties supported for scoping by this environment.
	  * Constraints: The maximum length is `8` items. The minimum length is `0` items.
	Nested schema for **properties**:
		* `name` - (String) The name of the property.
		  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters. The value must match regular expression `/[A-Za-z0-9]+/`.
		* `value` - (String) The value of the property.
		  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters. The value must match regular expression `/[A-Za-z0-9]+/`.
* `notifications` - (List) The request payload of the attachment notifications.
Nested schema for **notifications**:
	* `controls` - (List) The failed controls.
	Nested schema for **controls**:
		* `failed_control_ids` - (List) The failed control IDs.
		  * Constraints: The list items must match regular expression `/^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-4[0-9A-Fa-f]{3}-[89ABab][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}$|^$/`. The maximum length is `512` items. The minimum length is `0` items.
		* `threshold_limit` - (Integer) The threshold limit.
	* `enabled` - (Boolean) enabled notifications.
* `schedule` - (String) The schedule of an attachment evaluation.
  * Constraints: Allowable values are: `daily`, `every_7_days`, `every_30_days`.
* `name` - (String) The name of the attachment.
  * Constraints: The maximum length is `128` characters. The minimum length is `2` characters. The value must match regular expression `/^[a-zA-Z0-9-]*$/`.

## Attribute Reference

After your resource is created, you can read values from the listed arguments and the following attributes.

* `id` - The unique identifier of the scc_profile_attachment.
* `account_id` - (String) The account ID that is associated to the attachment.
  * Constraints: The maximum length is `32` characters. The minimum length is `32` characters. The value must match regular expression `/^[a-zA-Z0-9-]*$/`.
* `attachment_id` - (String) The ID of the attachment.
  * Constraints: The maximum length is `32` characters. The minimum length is `32` characters. The value must match regular expression `/^[a-zA-Z0-9-]*$/`.
* `attachment_parameters` - (List) The profile parameters for the attachment.
  * Constraints: The maximum length is `512` items. The minimum length is `0` items.
Nested schema for **attachment_parameters**:
	* `assessment_id` - (String) The implementation ID of the parameter.
	  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters. The value must match regular expression `/[A-Za-z0-9]+/`.
	* `assessment_type` - (String) The type of the implementation.
	  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters. The value must match regular expression `/[A-Za-z0-9]+/`.
	* `parameter_display_name` - (String) The parameter display name.
	  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters. The value must match regular expression `/^[a-zA-Z0-9_,'\\s\\-]*$/`.
	* `parameter_name` - (String) The parameter name.
	  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters. The value must match regular expression `/^[a-zA-Z0-9_]*$/`.
	* `parameter_type` - (String) The parameter type.
	  * Constraints: Allowable values are: `string`, `numeric`, `general`, `boolean`, `string_list`, `ip_list`, `timestamp`.
	* `parameter_value` - (String) The value of the parameter.
	  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters. The value must match regular expression `/^[a-zA-Z0-9_,'"\\s\\-\\[\\]]+$/`.
* `created_by` - (String) The user who created the attachment.
  * Constraints: The maximum length is `255` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9-\\.:,_\\s]*$/`.
* `created_on` - (String) The date when the attachment was created.
* `description` - (String) The description for the attachment.
  * Constraints: The maximum length is `256` characters. The minimum length is `0` characters. The value must match regular expression `/^[a-zA-Z0-9_,'\\s\\-]*$/`.
* `instance_id` - (String) The instance ID of the account that is associated to the attachment.
  * Constraints: The maximum length is `36` characters. The minimum length is `36` characters. The value must match regular expression `/^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-4[0-9A-Fa-f]{3}-[89ABab][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}$|^$/`.
* `last_scan` - (List) The details of the last scan of an attachment.
Nested schema for **last_scan**:
	* `id` - (String) The ID of the last scan of an attachment.
	  * Constraints: The maximum length is `36` characters. The minimum length is `36` characters. The value must match regular expression `/^[a-zA-Z0-9-]*$/`.
	* `status` - (String) The status of the last scan of an attachment.
	  * Constraints: Allowable values are: `in_progress`, `completed`.
	* `time` - (String) The time when the last scan started.
* `next_scan_time` - (String) The start time of the next scan.
* `status` - (String) The status of an attachment evaluation.
  * Constraints: Allowable values are: `enabled`, `disabled`.
* `updated_by` - (String) The user who updated the attachment.
  * Constraints: The maximum length is `255` characters. The minimum length is `1` character. The value must match regular expression `/^[a-zA-Z0-9-\\.:,_\\s]*$/`.
* `updated_on` - (String) The date when the attachment was updated.


## Import

You can import the `ibm_scc_profile_attachment` resource by using `id`.
The `id` property can be formed from `profiles_id`, and `attachment_id` in the following format:

```
<profiles_id>/<attachment_id>
```
* `profiles_id`: A string. The profile ID.
* `attachment_id`: A string. The attachment ID.

# Syntax
```
$ terraform import ibm_scc_profile_attachment.scc_profile_attachment <profiles_id>/<attachment_id>
```
