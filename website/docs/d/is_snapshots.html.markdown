---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : snapshots"
description: |-
  Reads IBM Cloud snapshots.
---
# ibm_is_snapshots

Import the details of an existing IBM Cloud infrastructure snapshot collection as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax, see [viewing snapshots](https://cloud.ibm.com/docs/vpc?topic=vpc-snapshots-vpc-view).

**Note:** 
VPC infrastructure services are a regional specific based endpoint, by default targets to `us-south`. Please make sure to target right region in the provider block as shown in the `provider.tf` file, if VPC service is created in region other than `us-south`.

**provider.tf**

```terraform
provider "ibm" {
  region = "eu-gb"
}
```


## Example usage

```terraform

data "ibm_is_snapshots" "example" {
}

```

## Argument reference
Review the argument references that you can specify for your data source. 

- `name` - (Optional, String) Filter snapshot collection by name of the snapshot.
- `resource_group` - (Optional, String) Filter snapshot collection by resource group id of the snapshot.
- `source_image` - (Optional, String) Filter snapshot collection by source image of the snapshot.
- `source_volume` - (Optional, String) Filter snapshot collection by source volume of the snapshot.
- `backup_policy_plan_tag` - Filters the collection to resources with the exact tag value.
- `backup_policy_plan_id` - Filters the collection to backup policy jobs with the backup plan with the specified identifier
- `snapshot_copies_id` - Filters the collection to snapshots with copies with the specified identifier.
- `snapshot_copies_name` - Filters the collection to snapshots with copies with the exact specified name.
- `snapshot_copies_crn` - Filters the collection to snapshots with copies with the specified CRN.
- `snapshot_copies_remote_region_name` - Filters the collection to snapshots with copies with the exact remote region name.
- `source_snapshot_id` - Filters the collection to resources with the source snapshot with the specified identifier
- `source_snapshot_remote_region_name` - Filters the collection to snapshots with a source snapshot with the exact remote region name.
- `snapshot_source_volume_remote_region_name` - Filters the collection to snapshots with a source volume with the exact remote region name.


## Attribute reference
In addition to all argument reference list, you can access the following attribute references after your data source is created.

- `snapshots` - (List) List of snapshots in the IBM Cloud Infrastructure.
  
  Nested scheme for `snapshots`:
  - `access_tags`  - (Array of Strings) Access management tags associated with the snapshot.
  - `id` - (String) The unique identifier for this snapshot.
  - `backup_policy_plan` - (List) If present, the backup policy plan which created this snapshot.
  
   Nested scheme for `backup_policy_plan`:
    - `deleted` - (List) If present, this property indicates the referenced resource has been deleted and provides some supplementary information.
   
      Nested scheme for `deleted`:
      - `more_info` - (String) Link to documentation about deleted resources.
    - `href` - (String) The URL for this backup policy plan.
    - `id` - (String) The unique identifier for this backup policy plan.
    - `name` - (String) The unique user defined name for this backup policy plan. If unspecified, the name will be a hyphenated list of randomly selected words.
    - `resource_type` - (String) The type of resource referenced.
  - `bootable` - (Bool) Indicates if a boot volume attachment can be created with a volume created from this snapshot.
  - `copies` - (List) The copies of this snapshot in other regions.
  
      Nested scheme for `copies`:
      - `crn` - (String) The CRN for the copied snapshot.
      - `deleted` - (List) If present, this property indicates the referenced resource has been deleted and provides some supplementary information.

          Nested scheme for `deleted`:
          - `more_info` - (String) Link to documentation about deleted resources.
      - `href` - (String) The URL for the copied snapshot.
      - `id` - (String) The unique identifier for the copied snapshot.
      - `name` - (String) The name for the copied snapshot. The name is unique across all snapshots in the copied snapshot's native region.
      - `remote` - (List) If present, this property indicates the referenced resource is remote to this region,and identifies the native region.
        Nested scheme for `remote`:
        - `href` - (String) The URL for this region.
        - `name` - (String) The globally unique name for this region.
      - `resource_type` - (String) The resource type.
  - `clones` - (List) The list of zones where clones of this snapshot exist.
  - `crn` - (String) The CRN for this snapshot.
  - `encryption` - (String) The type of encryption used on the source volume. Supported values are **provider_managed**, **user_managed** ]).
  - `encryption_key` - (String) The CRN of the `Key Protect Root Key` or `Hyper Protect Crypto Services Root Key` for this resource. The root key used to wrap the data encryption key for the source volume. This property will be present for volumes with an encryption type of `user_managed`.
  - `href` - (String) The URL for this snapshot.
  - `lifecycle_state` - (String) The lifecycle state of this snapshot. Supported values are **deleted**, **deleting**, **failed**, **pending**, **stable**, **updating**, **waiting**, **suspended**.
  - `minimum_capacity` - (Integer) The minimum capacity of a volume created from this snapshot. When a snapshot is created, this will be set to the capacity of the source_volume.
  - `operating_system` - (String) The globally unique name for the operating system included in this image.  
  - `resource_type` - (String) The resource type.
  - `size` - (Integer) The size of this snapshot rounded up to the next gigabyte.
  - `source_image` - (String) If present, the unique identifier for the image from which the data on this volume was most directly provisioned.
  - `source_snapshot` - (String) If present, the source snapshot this snapshot was created from.
    
     Nested scheme for `source_snapshot`:
      - `crn` - (String) The CRN of the source snapshot.
      - `deleted` - (List) If present, this property indicates the referenced resource has been deleted and provides some supplementary information.
    
          Nested scheme for `deleted`:
          - `more_info` - (String) Link to documentation about deleted resources.
      - `href` - (String) The URL for the source snapshot.
      - `id` - (String) The unique identifier for the source snapshot.
      - `name` - (String) The name for the source snapshot. The name is unique across all snapshots in the source snapshot's native region.
      - `remote` - (List) If present, this property indicates the referenced resource is remote to this region,and identifies the native region.
    
          Nested scheme for `remote`:
          - `href` - (String) The URL for this region.
          - `name` - (String) The globally unique name for this region.
      - `resource_type` - (String) The resource type.
  - `captured_at` - (String) The date and time that this snapshot was captured.
  - `tags` - (String) Tags associated with the snapshot.


