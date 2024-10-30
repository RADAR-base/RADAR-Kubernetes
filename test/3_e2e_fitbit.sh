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

path=$(dirname $BASH_SOURCE)
. $path/util.sh

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
echo "Get the subject id"

users_json=`curl -s "http://localhost/rest-sources/backend/users?project-id=$project_name&authorized=all" \
  -H 'Accept: application/json, text/plain, */*' \
  -H "Authorization: Bearer $mpToken"`
subject_id=`echo $users_json | jq -r ".users[] | select(.externalId == \"$subject_external_id\") | .userId"`
check_success "$subject_id" "subject_id"

user_json=`echo $users_json | jq -r ".users[] | select(.externalId == \"$subject_external_id\")"`
user_id=`echo $user_json | jq -r '.id'`
source_id=`echo $user_json | jq -r '.sourceId'`

user_json=`curl -s -X POST 'http://localhost/rest-sources/backend/users' \
  -H 'Accept: application/json, text/plain, */*' \
  -H "Authorization: Bearer $mpToken" \
  -H 'Content-Type: application/json' \
  --data @- <<EOF
{
  "projectId": ["$project_name"],
  "userId": "$user_id",
  "sourceId"ff: "$source_id",
  "startDate": "2024-10-28T23:00:00.000Z",
  "endDate": "2024-10-29T23:00:00.000Z",
  "sourceType": "FitBit"
}
EOF
`
check_success "$user_json" "user_json"

registration_json=`curl -s -X POST 'http://localhost/rest-sources/backend/registrations' \
  -H 'Accept: application/json, text/plain, */*' \
  -H "Authorization: Bearer $mpToken" \
  -H 'Content-Type: application/json' \
  --data @- <<EOF
{
  "userId": "$user_id",
  "persistent": true
}
EOF
`
check_success "$registration_json" "registration_json"

curl -X GET 'https://www.fitbit.com/oauth2/authorize?response_type=code&response_type=code&client_id=change_me&state=LIRlVdJro51o&scope=activity+heartrate+sleep+profile&prompt=login&redirect_uri=http%3A%2F%2Flocalhost%2Frest-sources%2Fauthorizer%2Fusers%3Anew' \
  -H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8' \
  -H 'cookie: fct=63cf43cf8a0e4f09b24c14b12c8e670a; JSESSIONID=30D3B4400C134EB59E1D169877B3BD39.fitbit1' \
  -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36'

echo "Register the FitBit token"
curl -s -X POST 'http://localhost/rest-sources/authorizer/user/new
