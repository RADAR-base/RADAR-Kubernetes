#!/usr/bin/env bash

# TODO remove this when minio mc is provided in the PATH
export PATH=$PATH:$HOME/minio-binaries/

host=${HOST_NAME:-localhost}
protocol=${PROTOCOL:-http}

# Get secrets from environment. Default to secrets from etc/secrets.yaml.
admin_username=admin
admin_password=${MP_ADMIN_PASSWORD:-`grep "common_admin_password" etc/secrets.yaml | awk '{print $NF}'`}
s3_access_key=${S3_ACCESS_KEY:-`grep "s3_access_key" etc/secrets.yaml | awk '{print $NF}'`}
s3_secret_key=${S3_SECRET_KEY:-`grep "s3_secret_key" etc/secrets.yaml | awk '{print $NF}'`}

org_name=${ORGANIZATION_NAME:-MAIN}
project_name=${PROJECT_NAME:-test}
subject_external_id=${SUBJECT_EXTERNAL_ID:-test_user}

# TEST LOGIC
test_s3_storage=${TEST_S3_STORAGE:-true}
s3_storage_timeout=${S3_STORAGE_TIMEOUT:-3}

if ! command -v mc 2>&1 >/dev/null && [ $test_s3_storage = "true" ]
then
    echo "The mc could not be found. Please in"
    exit 1
fi

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

# ------ SUBJECT -------
echo
echo "Get a new subject with external ID $subject_external_id"
subjects_json=`curl -s -X GET "http://localhost/managementportal/api/projects/$project_name/subjects" \
  --header 'Accept: application/json, text/plain, */*' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer '$mpToken`
subject_id=`echo $subjects_json | jq -r '.[] | select(.externalId == "'$subject_external_id'") | .login'`
check_success "$subject_id" "subject_id"
echo "Subject ID: $subject_id"

# ------ SOURCE -------
echo
echo "Get aRMT source for the $project_name project"
project_sources_json=`curl -s -X GET "http://localhost/managementportal/api/projects/$project_name/sources" \
  --header 'Accept: application/json, text/plain, */*' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer '$mpToken`
armt_project_source_id=`echo $project_sources_json | jq -r '.[] | select(.sourceType.name == "RADAR_aRMT") | .sourceId'`

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

if [ $test_s3_storage = "true" ]
then
  mc alias set s3-alias http://api.s3.localhost/ $s3_access_key $s3_secret_key
  object_count_intermediate_storage=`mc ls --recursive s3-alias/radar-intermediate-storage | wc -l`
  object_count_output_storage=`mc ls --recursive s3-alias/radar-output-storage | wc -l`
  echo "Intermediate storage object count: $object_count_intermediate_storage"
  echo "Output storage object count: $object_count_output_storage"
fi

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
echo $questionnaire_response_response
check_success "$questionnaire_response_response" "questionnaire_response_response"

if [ "$test_s3_storage" = "false" ]; then
  echo
  echo "Skipping S3 storage test"
  exit 0
fi

# Wait for 6 seconds
sleep $s3_storage_timeout

echo "Waiting for the data to be written to intermediate storage"
while [ $object_count_intermediate_storage -eq `mc ls --recursive s3-alias/radar-intermediate-storage | wc -l` ]; do
  echo "Waiting for the data to be written to intermediate storage"
  sleep 1
done
echo "Success!!"

echo "Waiting for the data to be written to output storage"
while [ $object_count_output_storage -eq `mc ls --recursive s3-alias/radar-output-storage | wc -l` ]; do
  echo "Waiting for the data to be written to output storage"
  sleep 1
done
echo "Success!!"
