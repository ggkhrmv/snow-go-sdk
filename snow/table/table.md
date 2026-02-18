# Table API

## Delete - Record in table

Path parameters: 

- apiVersion -> optional
- sysId -> Id of record to delete
- tableName

Query parameters:

- sysparmQueryNoDomain (bool, default: false) -> only for sys admins or users with role "query_no_domain_table_api"
  - restrict record search to only domains for which logged in user is configured

Status code: 204

Response: none

## Get - Table

Path parameters: 

- apiVersion -> optional
- tableName

Query parameters:
- name - value pairs -> filter result set (mutually exclusive with sysparm_query)
  - e.g. `&active=true` &state=closed&assigned_to=john.smith

sysparm_display_value	Determines the type of data returned, either the actual values from the database or the display values of the fields. Display values are manipulated based on the actual value in the database and user or system settings and preferences.
If returning display values, the value that is returned is dependent on the field type.
Choice fields: The database value may be a number, but the display value will be more descriptive.
Date fields: The database value is in UTC format, while the display value is based on the user's time zone.
Encrypted text: The database value is encrypted, while the displayed value is unencrypted based on the user's encryption context.
Reference fields: The database value is sys_id, but the display value is a display field of the referenced record.
Data type: String

Valid values:

true: Returns the display values for all fields.
false: Returns the actual values from the database.
all: Returns both actual and display values.
Default: false

Note:
There is no preferred method for setting this parameter. However, specifying the display value may cause performance issues since it is not reading directly from the database and may include referencing other fields and records. For more information on display values and actual values, see Table API FAQs (KB0534905).
sysparm_exclude_reference_link	Flag that indicates whether to exclude Table API links for reference fields.
Valid values:

true: Exclude Table API links for reference fields.
false: Include Table API links for reference fields.
Data type: Boolean

Default: false

sysparm_fields	Comma-separated list of fields to return in the response. Invalid fields are ignored.
Data type: String

Default: Return all fields.

sysparm_limit	Maximum number of records to return. For requests that exceed this number of records, use the sysparm_offset parameter to paginate record retrieval.
This limit is applied before ACL evaluation. If no records return, including records you have access to, rearrange the record order so records you have access to return first.

Note:
Unusually large sysparm_limit values can impact system performance.
Data type: Number

Default: 20

Maximum: 100

sysparm_no_count	Flag that indicates whether to execute a select count(*) query on the table to return the number of rows in the associated table.
Valid values:
true: Do not execute a select count(*).
false: Execute a select count(*).
Data type: Boolean

Default: false

sysparm_offset	Starting record index for which to begin retrieving records. Use this value to paginate record retrieval. This functionality enables the retrieval of all records, regardless of the number of records, in small manageable chunks.
For example, the first time you call this endpoint, sysparm_offset is set to "0". To simply page through all available records, use sysparm_offset=sysparm_offset+sysparm_limit, until you reach the end of all records.

Don't pass a negative number in the sysparm_offset parameter.
Data type: Number

Default: 0

sysparm_query	Encoded query used to filter the result set. You can use a UI filter to obtain a properly encoded query.
Syntax: sysparm_query=<col_name><operator><value>.
<col_name>: Name of the table column to filter against.
<operator>: Supports the following values:
=: Exactly matches <value>.
!=: Does not match <value>.
^: Logically AND multiple query statements.
^OR: Logically OR multiple query statements.
LIKE: <col_name> contains the specified string. Only works for <col_name> fields whose data type is string.
STARTSWITH: <col_name> starts with the specified string. Only works for <col_name> fields whose data type is string.
ENDSWITH: <col_name> ends with the specified string. Only works for <col_name> fields whose data type is string.
<value>: Value to match against.
For more operators, see Operators available for filters and queries.
All parameters are case-sensitive. Queries can contain more than one entry, such as sysparm_query=<col_name><operator><value>[<operator><col_name><operator><value>].

For example:

(sysparm_query=caller_id=javascript:gs.getUserID()^active=true)

Encoded queries also support order by functionality. To sort responses based on certain fields, use the ORDERBY and ORDERBYDESC clauses in sysparm_query.

Syntax:
ORDERBY<col_name>
ORDERBYDESC<col_name>
For example: sysparm_query=active=true^ORDERBYnumber^ORDERBYDESCcategory

This query filters all active records and orders the results in ascending order by number, and then in descending order by category.

If part of the query is invalid, such as by specifying an invalid field name, the instance ignores the invalid part. It then returns rows using only the valid portion of the query. You can control this behavior using the property glide.invalid_query.returns_no_rows. Set this property to true to return no rows on an invalid query.
Note:
The glide.invalid_query.returns_no_rows property controls the behavior of all queries across the instance, such as in lists, scripts (GlideRecord.query()), and web service APIs.
Data type: String

sysparm_query_category	Name of the category to use for queries.
Data type: String

sysparm_query_no_domain	Flag that indicates whether to restrict the record search to only the domains for which the logged in user is configured.
Valid values:

false: Exclude the record if it is in a domain that the currently logged in user is not configured to access.
true: Include the record even if it is in a domain that the currently logged in user is not configured to access.
Data type: Boolean

Default: false

Note:
The sysparm_query_no_domain parameter is available only to system administrators or users who have the query_no_domain_table_api role.
sysparm_suppress_pagination_header	Flag that indicates whether to remove the Link header from the response. The Link header provides various URLs to relative pages in the record set which you can use to paginate the returned record set.
Valid values:

true: Remove the Link header from the response.
false: Do not remove the Link header from the response.
Data type: Boolean

Default: false

sysparm_view	UI view for which to render the data. Determines the fields returned in the response.
Valid values:

desktop
mobile
both
If you also specify the sysparm_fields parameter, it takes precedent.

Data type: String

Status Code: 200

Response body parameters: name-value pairs

    example:
{
  "result": [
    {
      "parent": "",
      "made_sla": "true",
      "watch_list": "",
      "upon_reject": "cancel",
      "sys_updated_on": "2016-01-19 04:52:04",
      "approval_history": "",
      "number": "PRB0000050",
      "sys_updated_by": "glide.maint",
      "opened_by": {
        "link": "https://instance.servicenow.com/api/now/table/sys_user/glide.maint",
        "value": "glide.maint"
      },
      "user_input": "",
      "sys_created_on": "2016-01-19 04:51:19",
      "sys_domain": {
        "link": "https://instance.servicenow.com/api/now/table/sys_user_group/global",
        "value": "global"
      },
      "state": "4",
      "sys_created_by": "glide.maint",
      "knowledge": "false",
      "order": "",
      "closed_at": "2016-01-19 04:52:04",
      "cmdb_ci": {
        "link": "https://instance.servicenow.com/api/now/table/cmdb_ci/55b35562c0a8010e01cff22378e0aea9",
        "value": "55b35562c0a8010e01cff22378e0aea9"
      },
      "delivery_plan": "",
      "impact": "3",
      "active": "false",
      "work_notes_list": "",
      "business_service": "",
      "priority": "4",
      "sys_domain_path": "/",
      "time_worked": "",
      "expected_start": "",
      "rejection_goto": "",
      "opened_at": "2016-01-19 04:49:47",
      "business_duration": "1970-01-01 00:00:00",
      "group_list": "",
      "work_end": "",
      "approval_set": "",
      "wf_activity": "",
      "work_notes": "",
      "short_description": "Switch occasionally drops connections",
      "correlation_display": "",
      "delivery_task": "",
      "work_start": "",
      "assignment_group": "",
      "additional_assignee_list": "",
      "description": "Switch occasionally drops connections",
      "calendar_duration": "1970-01-01 00:02:17",
      "close_notes": "updated firmware",
      "sys_class_name": "problem",
      "closed_by": "",
      "follow_up": "",
      "sys_id": "04ce72c9c0a8016600b5b7f75ac67b5b",
      "contact_type": "phone",
      "urgency": "3",
      "company": "",
      "reassignment_count": "",
      "activity_due": "",
      "assigned_to": "",
      "comments": "",
      "approval": "not requested",
      "sla_due": "",
      "comments_and_work_notes": "",
      "due_date": "",
      "sys_mod_count": "1",
      "sys_tags": "",
      "escalation": "0",
      "upon_approval": "proceed",
      "correlation_id": "",
      "location": ""
    }
  ]
}

## Get - Table - SysId

Path parameters: 

- apiVersion -> optional
- sysId -> Id of record to delete
- tableName

Query parameters:

sysparm_display_value	Determines the type of data returned, either the actual values from the database or the display values of the fields. Display values are manipulated based on the actual value in the database and user or system settings and preferences.
If returning display values, the value that is returned is dependent on the field type.
Choice fields: The database value may be a number, but the display value will be more descriptive.
Date fields: The database value is in UTC format, while the display value is based on the user's time zone.
Encrypted text: The database value is encrypted, while the displayed value is unencrypted based on the user's encryption context.
Reference fields: The database value is sys_id, but the display value is a display field of the referenced record.
Data type: String

Valid values:

true: Returns the display values for all fields.
false: Returns the actual values from the database.
all: Returns both actual and display values.
Default: false

Note:
There is no preferred method for setting this parameter. However, specifying the display value may cause performance issues since it is not reading directly from the database and may include referencing other fields and records. For more information on display values and actual values, see Table API FAQs (KB0534905).
sysparm_exclude_reference_link	Flag that indicates whether to exclude Table API links for reference fields.
Valid values:

true: Exclude Table API links for reference fields.
false: Include Table API links for reference fields.
Data type: Boolean

Default: false

sysparm_fields	Comma-separated list of fields to return in the response. Invalid fields are ignored.
Data type: String

Default: Return all fields.

sysparm_query_no_domain	Flag that indicates whether to restrict the record search to only the domains for which the logged in user is configured.
Valid values:

false: Exclude the record if it is in a domain that the currently logged in user is not configured to access.
true: Include the record even if it is in a domain that the currently logged in user is not configured to access.
Data type: Boolean

Default: false

Note:
The sysparm_query_no_domain parameter is available only to system administrators or users who have the query_no_domain_table_api role.
sysparm_view	UI view for which to render the data. Determines the fields returned in the response.
Valid values:

desktop
mobile
both
If you also specify the sysparm_fields parameter, it takes precedent.

Data type: String

Status codes: 200, 404

Response body: name-value pairs

example:
{
  "result": {
    "upon_approval": "",
    "location": {
      "link": "https://instance.servicenow.com/api/now/table/cmn_location/105cf7f3c611227501e75e08b14a38ba",
      "value": "105cf7f3c611227501e75e08b14a38ba"
    },
    "expected_start": "",
    "reopen_count": "",
    "close_notes": "",
    "additional_assignee_list": "",
    "impact": "1",
    "urgency": "3",
    "correlation_id": "",
    "sys_tags": "",
    "sys_domain": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user_group/global",
      "value": "global"
    },
    "description": "",
    "group_list": "",
    "priority": "3",
    "delivery_plan": "",
    "sys_mod_count": "4",
    "work_notes_list": "",
    "business_service": "",
    "follow_up": "",
    "closed_at": "",
    "sla_due": "2015-11-11 22:04:15",
    "delivery_task": "",
    "sys_updated_on": "2015-11-01 22:37:27",
    "parent": "",
    "work_end": "",
    "number": "INC0000046",
    "closed_by": "",
    "work_start": "",
    "calendar_stc": "",
    "category": "software",
    "business_duration": "",
    "incident_state": "1",
    "activity_due": "",
    "correlation_display": "",
    "company": "",
    "active": "true",
    "due_date": "",
    "assignment_group": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user_group/8a4dde73c6112278017a6a4baf547aa7",
      "value": "8a4dde73c6112278017a6a4baf547aa7"
    },
    "caller_id": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user/46c6f9efa9fe198101ddf5eed9adf6e7",
      "value": "46c6f9efa9fe198101ddf5eed9adf6e7"
    },
    "knowledge": "false",
    "made_sla": "false",
    "comments_and_work_notes": "",
    "parent_incident": "",
    "state": "1",
    "user_input": "",
    "sys_created_on": "2015-11-01 22:05:30",
    "approval_set": "",
    "reassignment_count": "1",
    "rfc": "",
    "child_incidents": "",
    "opened_at": "2015-11-02 22:04:15",
    "short_description": "Can't access SFA software",
    "order": "",
    "sys_updated_by": "glide.maint",
    "resolved_by": "",
    "notify": "1",
    "upon_reject": "",
    "approval_history": "",
    "problem_id": {
      "link": "https://instance.servicenow.com/api/now/table/problem/a9e4890bc6112276003d7a5a5c774a74",
      "value": "a9e4890bc6112276003d7a5a5c774a74"
    },
    "work_notes": "",
    "calendar_duration": "",
    "close_code": "",
    "sys_id": "a9e30c7dc61122760116894de7bcc7bd",
    "approval": "not requested",
    "caused_by": "",
    "severity": "3",
    "sys_created_by": "admin",
    "resolved_at": "",
    "assigned_to": "",
    "business_stc": "",
    "wf_activity": "",
    "sys_domain_path": "/",
    "cmdb_ci": {
      "link": "https://instance.servicenow.com/api/now/table/cmdb_ci/a9c0c8d2c6112276018f7705562f9cb0",
      "value": "a9c0c8d2c6112276018f7705562f9cb0"
    },
    "opened_by": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user/46c6f9efa9fe198101ddf5eed9adf6e7",
      "value": "46c6f9efa9fe198101ddf5eed9adf6e7"
    },
    "subcategory": "",
    "rejection_goto": "",
    "sys_class_name": "incident",
    "watch_list": "",
    "time_worked": "",
    "contact_type": "phone",
    "escalation": "0",
    "comments": ""
  }
}

## Patch

Path parameters: 

- apiVersion -> optional
- sysId -> Id of record to delete
- tableName

Query parameters
Name	Description
sysparm_display_value	Determines the type of data returned, either the actual values from the database or the display values of the fields. Display values are manipulated based on the actual value in the database and user or system settings and preferences.
If returning display values, the value that is returned is dependent on the field type.
Choice fields: The database value may be a number, but the display value will be more descriptive.
Date fields: The database value is in UTC format, while the display value is based on the user's time zone.
Encrypted text: The database value is encrypted, while the displayed value is unencrypted based on the user's encryption context.
Reference fields: The database value is sys_id, but the display value is a display field of the referenced record.
Data type: String

Valid values:

true: Returns the display values for all fields.
false: Returns the actual values from the database.
all: Returns both actual and display values.
Default: false

Note:
There is no preferred method for setting this parameter. However, specifying the display value may cause performance issues since it is not reading directly from the database and may include referencing other fields and records. For more information on display values and actual values, see Table API FAQs (KB0534905).
sysparm_fields	Comma-separated list of fields to return in the response. Invalid fields are ignored.
Data type: String

Default: Return all fields.

sysparm_input_display_value	Flag that indicates whether to set field values using the display value or the actual value. Depending on the different types of fields, the endpoint may manipulate the passed in display values to store the proper values in the database. For example, if you send the display name for a reference field, the endpoint stores the sys_id for that value in the database. For date and time fields, when this parameter is true, the date and time value is adjusted for the current user's timezone. When false, the date and time value is inserted using the GMT timezone.
Valid values:

true: Treats input values as display values and they are manipulated so they are stored properly in the database.
false: Treats input values as actual values and stores them in the database without manipulation.
Data type: Boolean

Default: false - This matches the data type that is returned during data retrieval (GET methods), which is the actual values.

Note:
To set the value of an encrypted field, you must set this parameter to true. If this parameter is not set to true, values submitted to encrypted fields are not saved. Additionally, the requesting user must have the appropriate encryption context prior to submitting the request. Encrypted fields are hidden for users without the appropriate encryption context. For more information on field encryption see Encryption.
sysparm_query_no_domain	Flag that indicates whether to restrict the record search to only the domains for which the logged in user is configured.
Valid values:

false: Exclude the record if it is in a domain that the currently logged in user is not configured to access.
true: Include the record even if it is in a domain that the currently logged in user is not configured to access.
Data type: Boolean

Default: false

Note:
The sysparm_query_no_domain parameter is available only to system administrators or users who have the query_no_domain_table_api role.
sysparm_view	UI view for which to render the data. Determines the fields returned in the response.
Valid values:

desktop
mobile
both
If you also specify the sysparm_fields parameter, it takes precedent.

Data type: String

Request body:
name-value pairs:	
    Field name and the new value for each parameter to update in the specified record.
    Note:
    All fields within a record may not be available for update. For example, fields that have a prefix of "sys_" are typically system parameters that are automatically generated and cannot be updated.

Status codes: 200, 400, 404

Response body: name value pairs
example:
request:
curl "https://instance.servicenow.com/api/now/table/incident/ef43c6d40a0a0b5700c77f9bf387afe3" \
--request PATCH \
--header "Accept:application/json" \
--header "Content-Type:application/json" \
--data "{'assigned_to':'681b365ec0a80164000fb0b05854a0cd','urgency':'1','comments':'Elevating urgency, this is a blocking issue'}" \
--user 'username':'password'
response:
{
  "result": {
    "upon_approval": "proceed",
    "location": {
      "link": "https://instance.servicenow.com/api/now/table/cmn_location/108752c8c611227501d4ab0e392ba97f",
      "value": "108752c8c611227501d4ab0e392ba97f"
    },
    "expected_start": "",
    "reopen_count": "",
    "close_notes": "",
    "additional_assignee_list": "",
    "impact": "1",
    "urgency": "1",
    "correlation_id": "",
    "sys_tags": "",
    "sys_domain": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user_group/global",
      "value": "global"
    },
    "description": "",
    "group_list": "",
    "priority": "1",
    "delivery_plan": "",
    "sys_mod_count": "7",
    "work_notes_list": "",
    "business_service": "",
    "follow_up": "",
    "closed_at": "",
    "sla_due": "2017-07-05 05:58:24",
    "delivery_task": "",
    "sys_updated_on": "2016-01-22 14:12:37",
    "parent": "",
    "work_end": "",
    "number": "INC0000050",
    "closed_by": "",
    "work_start": "",
    "calendar_stc": "",
    "category": "hardware",
    "business_duration": "",
    "incident_state": "2",
    "activity_due": "2016-01-22 16:12:37",
    "correlation_display": "",
    "company": {
      "link": "https://instance.servicenow.com/api/now/table/core_company/31bea3d53790200044e0bfc8bcbe5dec",
      "value": "31bea3d53790200044e0bfc8bcbe5dec"
    },
    "active": "true",
    "due_date": "",
    "assignment_group": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user_group/8a5055c9c61122780043563ef53438e3",
      "value": "8a5055c9c61122780043563ef53438e3"
    },
    "caller_id": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user/5b7c200d0a640069006b3845b5d0fa7c",
      "value": "5b7c200d0a640069006b3845b5d0fa7c"
    },
    "knowledge": "false",
    "made_sla": "true",
    "comments_and_work_notes": "",
    "parent_incident": "",
    "state": "2",
    "user_input": "",
    "sys_created_on": "2015-11-02 18:05:40",
    "approval_set": "",
    "reassignment_count": "0",
    "rfc": "",
    "child_incidents": "",
    "opened_at": "2015-11-02 21:58:24",
    "short_description": "Can't access Exchange server - is it down?",
    "order": "",
    "sys_updated_by": "admin",
    "resolved_by": "",
    "notify": "1",
    "upon_reject": "cancel",
    "approval_history": "",
    "problem_id": "",
    "work_notes": "",
    "calendar_duration": "",
    "close_code": "",
    "sys_id": "ef43c6d40a0a0b5700c77f9bf387afe3",
    "approval": "not requested",
    "caused_by": "",
    "severity": "3",
    "sys_created_by": "glide.maint",
    "resolved_at": "",
    "assigned_to": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user/681b365ec0a80164000fb0b05854a0cd",
      "value": "681b365ec0a80164000fb0b05854a0cd"
    },
    "business_stc": "",
    "wf_activity": "",
    "sys_domain_path": "/",
    "cmdb_ci": {
      "link": "https://instance.servicenow.com/api/now/table/cmdb_ci/281190e3c0a8000b003f593aa3f20ca6",
      "value": "281190e3c0a8000b003f593aa3f20ca6"
    },
    "opened_by": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user/glide.maint",
      "value": "glide.maint"
    },
    "subcategory": "",
    "rejection_goto": "",
    "sys_class_name": "incident",
    "watch_list": "",
    "time_worked": "",
    "contact_type": "phone",
    "escalation": "3",
    "comments": ""
  }
}

## Post - Table - Insert singular record

Path parameters: 

- apiVersion -> optional
- tableName

Query parameters
Name	Description
sysparm_display_value	Determines the type of data returned, either the actual values from the database or the display values of the fields. Display values are manipulated based on the actual value in the database and user or system settings and preferences.
If returning display values, the value that is returned is dependent on the field type.
Choice fields: The database value may be a number, but the display value will be more descriptive.
Date fields: The database value is in UTC format, while the display value is based on the user's time zone.
Encrypted text: The database value is encrypted, while the displayed value is unencrypted based on the user's encryption context.
Reference fields: The database value is sys_id, but the display value is a display field of the referenced record.
Data type: String

Valid values:

true: Returns the display values for all fields.
false: Returns the actual values from the database.
all: Returns both actual and display values.
Default: false

Note:
There is no preferred method for setting this parameter. However, specifying the display value may cause performance issues since it is not reading directly from the database and may include referencing other fields and records. For more information on display values and actual values, see Table API FAQs (KB0534905).
sysparm_exclude_reference_link	Flag that indicates whether to exclude Table API links for reference fields.
Valid values:

true: Exclude Table API links for reference fields.
false: Include Table API links for reference fields.
Data type: Boolean

Default: false

sysparm_fields	Comma-separated list of fields to return in the response. Invalid fields are ignored.
Data type: String

Default: Return all fields.

sysparm_input_display_value	Flag that indicates whether to set field values using the display value or the actual value. Depending on the different types of fields, the endpoint may manipulate the passed in display values to store the proper values in the database. For example, if you send the display name for a reference field, the endpoint stores the sys_id for that value in the database. For date and time fields, when this parameter is true, the date and time value is adjusted for the current user's timezone. When false, the date and time value is inserted using the GMT timezone.
Valid values:

true: Treats input values as display values and they are manipulated so they are stored properly in the database.
false: Treats input values as actual values and stores them in the database without manipulation.
Data type: Boolean

Default: false - This matches the data type that is returned during data retrieval (GET methods), which is the actual values.

Note:
To set the value of an encrypted field, you must set this parameter to true. If this parameter is not set to true, values submitted to encrypted fields are not saved. Additionally, the requesting user must have the appropriate encryption context prior to submitting the request. Encrypted fields are hidden for users without the appropriate encryption context. For more information on field encryption see Encryption.
sysparm_view	UI view for which to render the data. Determines the fields returned in the response.
Valid values:

desktop
mobile
both
If you also specify the sysparm_fields parameter, it takes precedent.

Data type: String

Request body: name-value pairs
    Field name and the associated value for each parameter to define in the specified record.
    Note:
    All fields within a record may not be available for update. For example, fields that have a prefix of "sys_" are typically system parameters that are automatically generated and cannot be updated.
    Fields that are not specified and not auto generated by the system are set to the associated data type's null value.

Status codes: 201, 400, 404

Response: name-value pairs

example:
request:
curl "https://instance.servicenow.com/api/now/table/incident" \
--request POST \
--header "Accept:application/json" \
--header "Content-Type:application/json" \
--data "{'short_description':'Unable to connect to office wifi','assignment_group':'287ebd7da9fe198100f92cc8d1d2154e','urgency':'2','impact':'2'}" \
--user 'username':'password'

response:
{
  "result": {
    "upon_approval": "proceed",
    "location": "",
    "expected_start": "",
    "reopen_count": "0",
    "close_notes": "",
    "additional_assignee_list": "",
    "impact": "2",
    "urgency": "2",
    "correlation_id": "",
    "sys_tags": "",
    "sys_domain": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user_group/global",
      "value": "global"
    },
    "description": "",
    "group_list": "",
    "priority": "3",
    "delivery_plan": "",
    "sys_mod_count": "0",
    "work_notes_list": "",
    "business_service": "",
    "follow_up": "",
    "closed_at": "",
    "sla_due": "",
    "delivery_task": "",
    "sys_updated_on": "2016-01-22 14:28:24",
    "parent": "",
    "work_end": "",
    "number": "INC0010002",
    "closed_by": "",
    "work_start": "",
    "calendar_stc": "",
    "category": "inquiry",
    "business_duration": "",
    "incident_state": "1",
    "activity_due": "",
    "correlation_display": "",
    "company": "",
    "active": "true",
    "due_date": "",
    "assignment_group": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user_group/287ebd7da9fe198100f92cc8d1d2154e",
      "value": "287ebd7da9fe198100f92cc8d1d2154e"
    },
    "caller_id": "",
    "knowledge": "false",
    "made_sla": "true",
    "comments_and_work_notes": "",
    "parent_incident": "",
    "state": "1",
    "user_input": "",
    "sys_created_on": "2016-01-22 14:28:24",
    "approval_set": "",
    "reassignment_count": "0",
    "rfc": "",
    "child_incidents": "0",
    "opened_at": "2016-01-22 14:28:24",
    "short_description": "Unable to connect to office wifi",
    "order": "",
    "sys_updated_by": "admin",
    "resolved_by": "",
    "notify": "1",
    "upon_reject": "cancel",
    "approval_history": "",
    "problem_id": "",
    "work_notes": "",
    "calendar_duration": "",
    "close_code": "",
    "sys_id": "c537bae64f411200adf9f8e18110c76e",
    "approval": "not requested",
    "caused_by": "",
    "severity": "3",
    "sys_created_by": "admin",
    "resolved_at": "",
    "assigned_to": "",
    "business_stc": "",
    "wf_activity": "",
    "sys_domain_path": "/",
    "cmdb_ci": "",
    "opened_by": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user/6816f79cc0a8016401c5a33be04be441",
      "value": "6816f79cc0a8016401c5a33be04be441"
    },
    "subcategory": "",
    "rejection_goto": "",
    "sys_class_name": "incident",
    "watch_list": "",
    "time_worked": "",
    "contact_type": "phone",
    "escalation": "0",
    "comments": ""
  }
}

## Put - Table - Update record with request body

Path parameters: 

- apiVersion -> optional
- sysId -> Id of record to delete
- tableName

Query parameters
Name	Description
sysparm_display_value	Determines the type of data returned, either the actual values from the database or the display values of the fields. Display values are manipulated based on the actual value in the database and user or system settings and preferences.
If returning display values, the value that is returned is dependent on the field type.
Choice fields: The database value may be a number, but the display value will be more descriptive.
Date fields: The database value is in UTC format, while the display value is based on the user's time zone.
Encrypted text: The database value is encrypted, while the displayed value is unencrypted based on the user's encryption context.
Reference fields: The database value is sys_id, but the display value is a display field of the referenced record.
Data type: String

Valid values:

true: Returns the display values for all fields.
false: Returns the actual values from the database.
all: Returns both actual and display values.
Default: false

Note:
There is no preferred method for setting this parameter. However, specifying the display value may cause performance issues since it is not reading directly from the database and may include referencing other fields and records. For more information on display values and actual values, see Table API FAQs (KB0534905).
sysparm_exclude_reference_link	Flag that indicates whether to exclude Table API links for reference fields.
Valid values:

true: Exclude Table API links for reference fields.
false: Include Table API links for reference fields.
Data type: Boolean

Default: false

sysparm_fields	Comma-separated list of fields to return in the response. Invalid fields are ignored.
Data type: String

Default: Return all fields.

sysparm_input_display_value	Flag that indicates whether to set field values using the display value or the actual value.
Valid values:

true: Treats input values as display values and they are manipulated so they can be stored properly in the database.
false: Treats input values as actual values and stored them in the database without manipulation.
Data type: Boolean

Default: false

Note:
If this parameter is set to true, pay attention to input values, especially date values, as these are interpreted as being supplied via the user time zone preference and are transformed into UTC format.
To set the value of an encrypted field, you must set this parameter to true. If this parameter is not set to true, values submitted to encrypted fields are not saved. Additionally, the requesting user must have the appropriate encryption context prior to submitting the request. Encrypted fields are hidden for users without the appropriate encryption context. For more information on display values and actual values, see Table API FAQs (KB0534905). For more information on field encryption see Field Encryption.
sysparm_query_no_domain	Flag that indicates whether to restrict the record search to only the domains for which the logged in user is configured.
Valid values:

false: Exclude the record if it is in a domain that the currently logged in user is not configured to access.
true: Include the record even if it is in a domain that the currently logged in user is not configured to access.
Data type: Boolean

Default: false

Note:
The sysparm_query_no_domain parameter is available only to system administrators or users who have the query_no_domain_table_api role.
sysparm_view	UI view for which to render the data. Determines the fields returned in the response.
Valid values:

desktop
mobile
both
If you also specify the sysparm_fields parameter, it takes precedent.

Data type: String

Request body: name-value pairs
    Name-value pairs for the field(s) to update in the associated table. For example, to update the short description file, enter a name-value pair similar to the following: --data "{\"short_description\": \"my short desc\" }" \.

status codes: 200, 400, 404

response: name-value pairs

example:
request:
curl "https://instance.servicenow.com/api/now/table/incident/ef43c6d40a0a0b5700c77f9bf387afe3" \
--request PUT \
--header "Accept:application/json" \
--header "Content-Type:application/json" \
--data "{'assigned_to':'681b365ec0a80164000fb0b05854a0cd','urgency':'1','comments':'Elevating urgency, this is a blocking issue'}" \
--user 'username':'password'

response:
{
  "result": {
    "upon_approval": "proceed",
    "location": {
      "link": "https://instance.servicenow.com/api/now/table/cmn_location/108752c8c611227501d4ab0e392ba97f",
      "value": "108752c8c611227501d4ab0e392ba97f"
    },
    "expected_start": "",
    "reopen_count": "",
    "close_notes": "",
    "additional_assignee_list": "",
    "impact": "1",
    "urgency": "1",
    "correlation_id": "",
    "sys_tags": "",
    "sys_domain": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user_group/global",
      "value": "global"
    },
    "description": "",
    "group_list": "",
    "priority": "1",
    "delivery_plan": "",
    "sys_mod_count": "7",
    "work_notes_list": "",
    "business_service": "",
    "follow_up": "",
    "closed_at": "",
    "sla_due": "2017-07-05 05:58:24",
    "delivery_task": "",
    "sys_updated_on": "2016-01-22 14:12:37",
    "parent": "",
    "work_end": "",
    "number": "INC0000050",
    "closed_by": "",
    "work_start": "",
    "calendar_stc": "",
    "category": "hardware",
    "business_duration": "",
    "incident_state": "2",
    "activity_due": "2016-01-22 16:12:37",
    "correlation_display": "",
    "company": {
      "link": "https://instance.servicenow.com/api/now/table/core_company/31bea3d53790200044e0bfc8bcbe5dec",
      "value": "31bea3d53790200044e0bfc8bcbe5dec"
    },
    "active": "true",
    "due_date": "",
    "assignment_group": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user_group/8a5055c9c61122780043563ef53438e3",
      "value": "8a5055c9c61122780043563ef53438e3"
    },
    "caller_id": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user/5b7c200d0a640069006b3845b5d0fa7c",
      "value": "5b7c200d0a640069006b3845b5d0fa7c"
    },
    "knowledge": "false",
    "made_sla": "true",
    "comments_and_work_notes": "",
    "parent_incident": "",
    "state": "2",
    "user_input": "",
    "sys_created_on": "2015-11-02 18:05:40",
    "approval_set": "",
    "reassignment_count": "0",
    "rfc": "",
    "child_incidents": "",
    "opened_at": "2015-11-02 21:58:24",
    "short_description": "Can't access Exchange server - is it down?",
    "order": "",
    "sys_updated_by": "admin",
    "resolved_by": "",
    "notify": "1",
    "upon_reject": "cancel",
    "approval_history": "",
    "problem_id": "",
    "work_notes": "",
    "calendar_duration": "",
    "close_code": "",
    "sys_id": "ef43c6d40a0a0b5700c77f9bf387afe3",
    "approval": "not requested",
    "caused_by": "",
    "severity": "3",
    "sys_created_by": "glide.maint",
    "resolved_at": "",
    "assigned_to": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user/681b365ec0a80164000fb0b05854a0cd",
      "value": "681b365ec0a80164000fb0b05854a0cd"
    },
    "business_stc": "",
    "wf_activity": "",
    "sys_domain_path": "/",
    "cmdb_ci": {
      "link": "https://instance.servicenow.com/api/now/table/cmdb_ci/281190e3c0a8000b003f593aa3f20ca6",
      "value": "281190e3c0a8000b003f593aa3f20ca6"
    },
    "opened_by": {
      "link": "https://instance.servicenow.com/api/now/table/sys_user/glide.maint",
      "value": "glide.maint"
    },
    "subcategory": "",
    "rejection_goto": "",
    "sys_class_name": "incident",
    "watch_list": "",
    "time_worked": "",
    "contact_type": "phone",
    "escalation": "3",
    "comments": ""
  }
}