#!/usr/bin/env bash

# TODO remove this when minio mc is provided in the PATH
export PATH=$PATH:$HOME/minio-binaries/

host=${HOST:-localhost}
protocol=${PROTOCOL:-http}

admin_username=admin
admin_password=`grep "common_admin_password" etc/secrets.yaml | awk '{print $NF}'`
s3_access_key=`grep "s3_access_key" etc/secrets.yaml | awk '{print $NF}'`
s3_secret_key=`grep "s3_secret_key" etc/secrets.yaml | awk '{print $NF}'`

project_name=TEST
org_name=$project_name
subject_external_id=test_user

check_success() {
  if [ -z "$1" ] || [ "$1" = "null" ]; then
    echo "Error: $2 is null"
    exit 1
  elif $(grep -q "error" <<< "$1"); then
    echo "Error: $2"
    exit 1
  else
    echo "Success!!"
  fi
}

assert_equals() {
  if [ "$1" != "$2" ]; then
    echo "Error: $3"
    exit 1
  else
    echo "Success!!"
  fi
}

echo "Starting e2e test on $protocol://$host"

echo
echo "Request an admin token for Management Portal"
response=`curl -s "$protocol://$host/managementportal/oauth/token" \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-raw "client_id=ManagementPortalapp&username=$admin_username&password=$admin_password&grant_type=password"`
mpToken=`echo $response | jq -r '.access_token'`
check_success "$mpToken" "mpToken"

echo
echo "List all source types in Management Portal"
source_types_json=`curl -s "$protocol://$host/managementportal/api/source-types" \
  -H 'Accept: application/json, text/plain, */*' \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $mpToken"`
android_source_type_json=`echo $source_types_json | jq -r '.[] | select(.name == "ANDROID_PHONE")'`
prmt_source_type_json=`echo $source_types_json | jq -r '.[] | select(.name == "RADAR_pRMT")'`
armt_source_type_json=`echo $source_types_json | jq -r '.[] | select(.name == "RADAR_aRMT")'`

if [ -z "$android_source_type_json" ]; then
  echo
  echo "Register Android Source Type"
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

if [ -z "$prmt_source_type_json" ]; then
  echo
  echo "Register pRMT Source Type"
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

if [ -z "$armt_source_type_json" ]; then
  echo
  echo "Register aRMT Source Type"
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

echo
echo "Create a new organization with name $org_name"
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
check_success "$organization_json" "organization_json"

echo
echo "Create a new project with name $project_name"
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
check_success "$project_json" "project_json"

echo
echo "Create a new aRMT source for the $project_name project"
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
check_success "$armt_project_source_json" "armt_project_source_json"
armt_project_source_id=`echo $armt_project_source_json | jq -r '.sourceId'`

echo
echo "Create a new subject with external ID $subject_external_id"
subject_json=`curl -s "$protocol://$host/managementportal/api/subjects" \
  -H 'Accept: application/json, text/plain, */*' \
  -H 'Content-Type: application/json' \
  --header 'Authorization: Bearer '$mpToken \
  --data @- <<EOF
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
subject_id=`echo $subject_json | jq -r '.login'`
check_success "$subject_id" "subject_id"

echo
echo "Retrieve the metatoken for the $subject_external_id subject"
metatoken_json=`curl -s "$protocol://$host/managementportal/api/oauth-clients/pair?clientId=aRMT&login=$subject_id&persistent=true" \
  -H 'Accept: application/json, text/plain, */*' \
  -H 'Content-Type: application/json' \
  --header 'Authorization: Bearer '$mpToken`
metatoken=`echo $metatoken_json | jq -r '.tokenName'`
check_success "$metatoken" "metatoken"

echo
echo "Retrieve the refreshtoken with the metatoken"
refreshtoken_json=`curl -s --location "$protocol://$host/managementportal/api/meta-token/$metatoken" \
   --header 'Authorization: Bearer '$mpToken`
refreshtoken=`echo $refreshtoken_json | jq -r '.refreshToken'`
check_success "$refreshtoken" "refreshtoken"

echo
echo "Retrieve the accesstoken with the refreshtoken"
accesstoken_json=`curl -s -X POST --location "$protocol://$host/managementportal/oauth/token" \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode "grant_type=refresh_token" \
  --data-urlencode "refresh_token=$refreshtoken" \
  --header "Authorization: Basic $(echo -n 'aRMT:' | base64)"`
accesstoken=`echo $accesstoken_json | jq -r '.access_token'`
check_success $accesstoken "accesstoken"

echo
echo "List all topics in Kafka"
topics=`curl -s --location "$protocol://$host/kafka/topics" \
  --header "Authorization: Bearer $accesstoken"`
check_success "$topics" "topics"

avro_questionnaire_response_key_id=`curl -s --location "$protocol://$host/schema/subjects/questionnaire_response-key/versions/latest" | jq -r '.id'`
avro_questionnaire_response_value_id=`curl -s --location "$protocol://$host/schema/subjects/questionnaire_response-value/versions/latest" | jq -r '.id'`

echo
echo "Push a questionnaire_response message to the questionnaire_response topic"
questionnaire_response_response=`curl -s --location "$protocol://$host/kafka/topics/questionnaire_response" \
  --header 'accept: application/vnd.kafka.v2+json, application/vnd.kafka+json; q=0.9, application/json; q=0.8' \
  --header "Authorization: Bearer $accesstoken" \
  --header 'content-type: application/json' \
  --data @- <<EOF
{
    "key_schema_id": $avro_questionnaire_response_key_id,
    "value_schema_id": $avro_questionnaire_response_value_id,
    "records": [
        {
            "value": {
                "time": 1707657463.9095526,
                "timeCompleted": 1707657463.9095526,
                "timeNotification": 0,
                "name": "QuestionnaireName",
                "version": "1",
                "answers": [
                    {"questionId" : "1", "value": "Some Value", "startTime": 0, "endTime": 0}
                ]
            },
            "key": {
                "projectId": {"string": "$project_name"},
                "userId": "$subject_id",
                "sourceId": "$armt_project_source_id"
            }
        }
    ]
}
EOF
`
check_success "$questionnaire_response_response" "questionnaire_response_response"

# Wait for 6 seconds
sleep 3

mc alias set s3-alias http://api.s3.localhost/ $s3_access_key $s3_secret_key

echo
echo "Check if the data is written to intermediate storage"
intermediate_object_list=`mc ls --recursive s3-alias/radar-intermediate-storage`
assert_equals $(echo $intermediate_object_list | grep 'questionnaire_response' | wc -l) 1 "questionnaire_response data is not written to intermediate storage"

echo
echo "Check if the data is written to output storage"
output_object_list=`mc ls --recursive s3-alias/radar-output-storage`
assert_equals $(echo $output_object_list | grep 'questionnaire_response' | wc -l) 1 "questionnaire_response data is not written to output storage"
