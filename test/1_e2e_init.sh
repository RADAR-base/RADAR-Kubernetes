#!/usr/bin/env bash

host=${HOST_NAME:-localhost}
protocol=${PROTOCOL:-http}

# Get secrets from environment. Default to secrets from etc/secrets.yaml.
admin_username=admin
admin_password=${MP_ADMIN_PASSWORD:-`grep "common_admin_password" etc/secrets.yaml | awk '{print $NF}'`}

org_name=${ORGANIZATION_NAME:-MAIN}
project_name=${PROJECT_NAME:-test}
subject_external_id=${SUBJECT_EXTERNAL_ID:-test_user}

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source $SCRIPT_DIR/util.sh

echo "Starting e2e test on $protocol://$host"

echo
echo "Request an admin token for Management Portal"
response=`curl -s "$protocol://$host/managementportal/oauth/token" \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-raw "client_id=ManagementPortalapp&username=$admin_username&password=$admin_password&grant_type=password"`
mpToken=`echo $response | jq -r '.access_token'`
check_success "$mpToken" "mpToken"


# ------ SOURCE TYPES -------
echo
echo "List all source types in Management Portal"
source_types_json=`curl -s "$protocol://$host/managementportal/api/source-types" \
  -H 'Accept: application/json, text/plain, */*' \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $mpToken"`

echo "Get or create Android Source Type"
android_source_type_json=`echo $source_types_json | jq -r '.[] | select(.name == "ANDROID_PHONE")'`
if [ -z "$android_source_type_json" ]; then
  echo
  echo "Create new Android Source Type"
  android_source_type_json=`curl -s "$protocol://$host/managementportal/api/source-types" \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $mpToken" \
    --data @- <<EOF
{
    "producer": "ANDROID",
    "model": "PHONE",
    "catalogVersion": "1.0.0",
    "sourceTypeScope": "PASSIVE",
    "canRegisterDynamically": true,
    "name": "ANDROID_PHONE"
}
EOF
  `
fi
android_source_id=`echo $android_source_type_json | jq -r '.name'`
check_success "$android_source_id" "android_source_id"

echo "Get or create pRMT Source Type"
prmt_source_type_json=`echo $source_types_json | jq -r '.[] | select(.name == "RADAR_pRMT")'`
if [ -z "$prmt_source_type_json" ]; then
  echo
  echo "Create new pRMT Source Type"
  prmt_source_type_json=`curl -s "$protocol://$host/managementportal/api/source-types" \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $mpToken" \
    --data @- <<EOF
{
    "appProvider": "org.radarcns.application.ApplicationServiceProvider",
    "producer": "RADAR",
    "model": "pRMT",
    "catalogVersion": "1.1.0",
    "canRegisterDynamically": true,
    "name": "RADAR_pRMT",
    "sourceTypeScope": "PASSIVE"
}
EOF
  `
fi
prmt_source_id=`echo $prmt_source_type_json | jq -r '.name'`
check_success "$prmt_source_id" "prmt_source_id"

echo "Get or create aRMT Source Type"
armt_source_type_json=`echo $source_types_json | jq -r '.[] | select(.name == "RADAR_aRMT")'`
if [ -z "$armt_source_type_json" ]; then
  echo
  echo "Create new aRMT Source Type"
  armt_source_type_json=`curl -s "$protocol://$host/managementportal/api/source-types" \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $mpToken" \
    --data @- <<EOF
{
    "appProvider": "org.radarcns.application.ApplicationServiceProvider",
    "producer": "RADAR",
    "model": "aRMT",
    "catalogVersion": "1.5.0",
    "canRegisterDynamically": true,
    "name": "RADAR_aRMT",
    "sourceTypeScope": "ACTIVE"
}
EOF
`
fi
armt_source_id=`echo $armt_source_type_json | jq -r '.name'`
check_success "$armt_source_id" "armt_source_id"


# ------ ORGANIZATION -------
echo
echo "Get or create a new organization with name $org_name"
organizations_json=`curl -s -X GET 'http://localhost/managementportal/api/organizations' \
  --header 'Accept: application/json, text/plain, */*' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer '$mpToken`
organization_json=`echo $organizations_json | jq -r '.[] | select(.name == "'$org_name'")'`
if [ -z "$organization_json" ]; then
  echo "Creating a new organization"
  organization_json=`curl -s --request POST \
    --url "$protocol://$host/managementportal/api/organizations" \
    --header 'Accept: application/json, text/plain, */*' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Bearer '$mpToken \
    --data @- <<EOF
{
    "name": "$org_name",
    "description": "TEST",
    "location": "TEST"
}
EOF
  `
fi
check_success "$organization_json" "organization_json"


# ------ PROJECT -------
echo
echo "Get or create a new project with name $project_name"
projects_json=`curl -s -X GET 'http://localhost/managementportal/api/projects' \
  --header 'Accept: application/json, text/plain, */*' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer '$mpToken`
project_json=`echo $projects_json | jq -r '.[] | select(.projectName == "'$project_name'")'`
if [ -z "$project_json" ]; then
  echo "Creating a new project"
  project_json=`curl -s "$protocol://$host/managementportal/api/projects" \
    -H 'Accept: application/json, text/plain, */*' \
    -H 'Content-Type: application/json' \
    --header 'Authorization: Bearer '$mpToken \
    --data @- <<EOF
{
    "organization": $organization_json,
    "projectName": "$project_name",
    "description": "$project_name",
    "location": "",
    "sourceTypes": [
      $armt_source_type_json
    ],
    "startDate": null,
    "endDate": null
}
EOF
  `
fi
check_success "$project_json" "project_json"


# ------ aRMT PROJECT SOURCE -------
echo
echo "Get or create a new aRMT source for the $project_name project"
project_sources_json=`curl -s -X GET "http://localhost/managementportal/api/projects/$project_name/sources" \
  --header 'Accept: application/json, text/plain, */*' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer '$mpToken`
armt_project_source_json=`echo $project_sources_json | jq -r '.[] | select(.sourceType.name == "RADAR_aRMT")'`
if [ -z "$armt_project_source_json" ]; then
  echo "Creating a new aRMT project source"
  armt_project_source_json=`curl -s "$protocol://$host/managementportal/api/sources" \
    -H 'Accept: application/json, text/plain, */*' \
    -H 'Content-Type: application/json' \
    --header "Authorization: Bearer $mpToken" \
    --data @- <<EOF
{
    "sourceName": "aRMT-test-source-$project_name",
    "assigned": false,
    "sourceType": $armt_source_type_json,
    "project": $project_json
}
EOF
  `
fi
check_success "$armt_project_source_json" "armt_project_source_json"
armt_project_source_id=`echo $armt_project_source_json | jq -r '.sourceId'`


# ------ SUBJECT -------
echo
echo "Get or create a new subject with external ID $subject_external_id"
subjects_json=`curl -s -X GET "http://localhost/managementportal/api/projects/$project_name/subjects" \
  --header 'Accept: application/json, text/plain, */*' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer '$mpToken`
subject_json=`echo $subjects_json | jq -r '.[] | select(.externalId == "'$subject_external_id'")'`
if [ -z "$subject_json" ]; then
  echo "Creating a new subject"
  subject_json=`curl -s "$protocol://$host/managementportal/api/subjects" \
    -H 'Accept: application/json, text/plain, */*' \
    -H 'Content-Type: application/json' \
    --header 'Authorization: Bearer '$mpToken \
    --data @- << EOF
{
    "project": $project_json,
    "sources": [
        $armt_project_source_json
    ],
    "status": 1,
    "externalId": "$subject_external_id",
    "group": null
}
EOF
  `
fi
subject_id=`echo $subject_json | jq -r '.login'`
check_success "$subject_id" "subject_id"
