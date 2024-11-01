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
fitbit_url=http://mockserver:8080

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


# ------ SUBJECT -------

echo
echo "Get the subject id"
users_json=`curl -s "$protocol://$host/rest-sources/backend/projects/$project_name/users" \
  -H 'Accept: application/json, text/plain, */*' \
  -H "Authorization: Bearer $mpToken"`
user_json=`echo $users_json | jq -r ".users[] | select(.externalId == \"$subject_external_id\")"`
check_success "$user_json" "user_json"
subject_id=`echo $user_json | jq -r '.id'`
check_success "$subject_id" "subject_id"


# ------ AUTHORIZER USER -------

echo
echo "Get the user id"
users_json=`curl -s "$protocol://$host/rest-sources/backend/users?project-id=$project_name&authorized=all" \
  -H 'Accept: application/json, text/plain, */*' \
  -H "Authorization: Bearer $mpToken"`
user_json=`echo $users_json | jq -r ".users[] | select(.externalId == \"$subject_external_id\")"`
if [ -z "$user_json" ]; then
  echo "Create a new user in rest sources authorizer backend"
  user_json=`curl -s -X "$protocol://$host/rest-sources/backend/users" \
      -H 'Accept: application/json, text/plain, */*' \
      -H 'Content-Type: application/json' \
      --data @- <<EOF
{
  "projectId": "$project_name",
  "userId": "$subject_id",
  "startDate": "2024-11-18T23:00:00.000Z",
  "endDate": "2024-11-29T23:00:00.000Z",
  "sourceType": "FitBit"
}
EOF
  `
fi
user_id=`echo $user_json | jq -r '.id'`
check_success "$user_id" "user_id"


# ------ FITBIT REGISTRATION -------

echo
echo "Register a new Fitbit user"
registration_json=`curl -s -X POST "$protocol://$host/rest-sources/backend/registrations" \
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
registration_token=`echo $registration_json | jq -r '.token'`
registration_secret=`echo $registration_json | jq -r '.secret'`

echo
echo "Registration token: $registration_token"
echo "Registration secret: $registration_secret"

# At this moment the user is redirected to Fitbit authorization page.
# The user is redirected to RADAR-base frontend URL callback with the 'code' (of authorization_code grant).
# The frontend redirects to the backend with the 'code' and the 'state' (registration_token).
echo
echo "Return the code that the user obtained from Fitbit authorization page to backend"
code=mycode
authorization_json=`curl -s "$protocol://$host/rest-sources/backend/registrations/$registration_token/authorize" \
  -H 'accept: application/json, text/plain, */*' \
  -H 'content-type: application/json' \
  --data-raw '{"code":"'$code'"}'`
check_success "$authorization_json" "authorization_json"

# The rest source backend will exchange the code for an access token and refresh token.
# This happens in the back chanel and is not part of the e2e test.
# the call will be to $fitbit_url/oauth2/token?client_id=$client;client_secret=some_secret;redirect_uri=$protocol://$host/rest-sources/backend/users:new;grant_type=authorization_code
# And body: '{"code":"'$code'"}'
# The response is provided by the mock server.
