import base64
import random

import yaml
from functools import reduce
import operator
import requests
import json
from minio import Minio
import time
import datetime
import os

minio_client = None

def wait_for_management_portal(context):
    time = 10
    while True:
        if time < 0:
            raise Exception('Management portal did not start in time')
        response = requests.get(f'{context.config.userdata["url"]}/managementportal')
        if response.status_code == 200:
            break
        time.sleep(1)
        time -= 1

def get_secret(*path_elements, context):
    if context.cache["secrets"] is None or len(context.cache["secrets"].keys()) == 0:
        secrets_file = context.config.userdata["secrets_file"]
        print("reading secrets from file...", secrets_file)
        with open(secrets_file, 'r') as file:
            context.cache["secrets"] = yaml.safe_load(file)
    return reduce(operator.getitem, path_elements, context.cache["secrets"])

def format_url(path, context):
    return f'{context.config.userdata["url"]}/{path}'

def get_mp_token(context):
    if context.cache["management_portal_token"] is not None:
        return context.cache["management_portal_token"]
    mp_admin_password = get_secret('management_portal', 'managementportal', 'common_admin_password', context=context)
    headers = {
        'Accept': 'application/json',
        'Content-Type': 'application/x-www-form-urlencoded'
    }
    data={
        'client_id': 'ManagementPortalapp',
        'username': 'admin',
        'password': mp_admin_password,
        'grant_type': 'password',
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/oauth/token', headers=headers, data=data)
    assert response.status_code == 200
    management_portal_json = response.json()
    token = management_portal_json['access_token']
    assert token is not None
    context.cache["management_portal_token"] = token
    return token

def get_source_types(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/source-types', headers=headers)
    assert response.status_code == 200
    return response.json()

def check_armt_source_type_exists(context, source_type):
    source_types = get_source_types(context)
    armt_source_type_json = next((item for item in source_types if item["name"] == source_type), None)
    if not context.config.userdata["dev_mode"] and armt_source_type_json is not None:
        raise Exception('aRMT source type already exists')
    context.cache["armt_source_type_json"] = armt_source_type_json

def create_armt_source_type(context, source_type):
    if context.cache["armt_source_type_json"] is not None:
        return
    armt_source_type_id = source_type
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    data = {
        "appProvider": "org.radarcns.application.ApplicationServiceProvider",
        "producer": "RADAR",
        "model": "aRMT",
        "catalogVersion": "1.5.0",
        "sourceTypeScope": "ACTIVE",
        "canRegisterDynamically": True,
        "name": armt_source_type_id
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/api/source-types', headers=headers, data=json.dumps(data))
    # 409 means that the armt source already exists.
    assert response.status_code == 201 or response.status_code == 409
    if response.status_code == 201:
        armt_source_type_json = response.json()
        assert armt_source_type_json is not None
        context.cache["armt_source_type_json"] = armt_source_type_json

def check_organization_exists(context, organization):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application /json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/organizations', headers=headers)
    assert response.status_code == 200
    test_organization_json = next((item for item in response.json() if item["name"] == organization), None)
    if not context.config.userdata["dev_mode"] and test_organization_json is not None:
        raise Exception('Test organization already exists')
    context.cache["organization_json"] = test_organization_json

def create_organization(context, organization):
    if context.cache["organization_json"] is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    data = {
        "name": organization,
        "description": "TEST",
        "location": "TEST"
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/api/organizations', headers=headers, data=json.dumps(data))
    # 409 means that the organization already exists
    assert response.status_code == 201 or response.status_code == 409
    if response.status_code == 201:
        test_organization_json = response.json()
        assert test_organization_json is not None
    context.cache["organization_json"] = test_organization_json

def check_project_exists(context, project):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/projects', headers=headers)
    assert response.status_code == 200
    test_project_json = next((item for item in response.json() if item["projectName"] == project), None)
    if context.config.userdata["dev_mode"] and test_project_json is not None:
        context.cache["project_json"] = test_project_json
    elif test_project_json is not None:
        raise Exception('Test project already exists')

def create_project(context, project):
    if context.cache["project_json"] is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    data = {
        "organization": context.cache["organization_json"],
        "projectName": project,
        "description": "TEST",
        "location": "",
        "sourceTypes": [context.cache["armt_source_type_json"]],
        "startDate": None,
        "endDate": None
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/api/projects', headers=headers, data=json.dumps(data))
    assert response.status_code == 200 or response.status_code == 201
    if response.status_code == 201:
        test_project_json = response.json()
        assert test_project_json is not None
        assert test_project_json['projectName'] == project
    context.cache["project_json"] = test_project_json

def check_armt_project_source_exists(context, source):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/projects/{context.config.userdata["project_name"]}/sources', headers=headers)
    assert response.status_code == 200
    armt_project_source_json = next((item for item in response.json() if item["sourceType"]["name"] == "RADAR_aRMT"), None)
    if context.config.userdata["dev_mode"] and armt_project_source_json is not None:
        context.cache["armt_project_source_json"] = armt_project_source_json
    elif armt_project_source_json is not None:
        raise Exception('aRMT project source already exists')

def create_armt_project_source(context, source):
    if context.cache["armt_project_source_json"] is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    data = {
        "sourceName": source,
        "assigned": False,
        "sourceType": context.cache["armt_source_type_json"],
        "project": context.cache["project_json"]
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/api/sources', headers=headers, data=json.dumps(data))
    assert response.status_code == 200 or response.status_code == 201
    armt_project_source_json = response.json()
    assert armt_project_source_json['sourceId'] is not None
    context.cache["armt_project_source_json"] = armt_project_source_json

def check_subject_exists(context, subject):
    project = context.cache["project_json"]["projectName"]
    assert project is not None
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/projects/{project}/subjects', headers=headers)
    assert response.status_code == 200
    test_subject_json = next((item for item in response.json() if item["externalId"] == subject), None)
    if context.config.userdata["dev_mode"] and test_subject_json is not None:
        context.cache["test_subject_id"] = test_subject_json['login']
    elif test_subject_json is not None:
        raise Exception('Test subject already exists')

def create_subject(context, subject):
    if context.cache["test_subject_id"] is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    data = {
        "project": context.cache["project_json"],
        "sources": [context.cache["armt_project_source_json"]],
        "status": 1,
        "externalId": subject,
        "group": None
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/api/subjects', headers=headers, data=json.dumps(data))
    assert response.status_code == 200 or response.status_code == 201
    test_subject_json = response.json()
    assert test_subject_json['login'] is not None
    context.cache["test_subject_id"] = test_subject_json['login']

def get_armt_meta_token(context):
    if context.cache["armt_meta_token"] is not None:
        return context.cache["armt_meta_token"]
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    params = {
        'clientId': 'aRMT',
        'login': context.cache["test_subject_id"],
        'persistent': True
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/oauth-clients/pair', headers=headers, params=params)
    assert response.status_code == 200
    meta_token = response.json()['tokenName']
    assert meta_token is not None
    context.cache["armt_meta_token"] = meta_token
    return meta_token

def get_armt_refresh_token(context):
    if context.cache["armt_refresh_token"] is not None:
        return context.cache["armt_refresh_token"]
    headers = {
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/meta-token/{get_armt_meta_token(context)}', headers=headers)
    assert response.status_code == 200
    refresh_token = response.json()['refreshToken']
    assert refresh_token is not None
    context.cache["armt_refresh_token"] = refresh_token
    return refresh_token

def get_armt_access_token(context):
    if context.cache["armt_access_token"] is not None:
        return context.cache["armt_access_token"]
    headers = {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Authorization': f'Basic {base64.b64encode(b"aRMT:").decode("utf-8")}'
    }
    token = get_armt_refresh_token(context)
    data = {
        'grant_type': 'refresh_token',
        'refresh_token': token
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/oauth/token', headers=headers, data=data)
    assert response.status_code == 200
    access_token = response.json()['access_token']
    assert access_token is not None
    context.cache["armt_access_token"] = access_token
    return access_token

def get_current_s3_object_state(context):
    minio_client = _get_minio_client(context)
    for (bucket, filename_pattern, change_type) in context.table:
        state_dict = _get_bucket_object_state(minio_client, bucket, filename_pattern)
        _cache_bucket_object_state(context, bucket, filename_pattern, state_dict)

def _get_current_s3_objects(minio_client, bucket, pattern):
    objects = minio_client.list_objects(
        bucket_name=bucket,
        recursive=True
    )
    filtered_objects = [ object for object in objects if pattern in object.object_name ]
    return filtered_objects

def _get_bucket_object_state(minio_client, bucket, pattern):
    objects = _get_current_s3_objects(minio_client, bucket, pattern)
    return { object.object_name: object.last_modified for object in objects }

# number = 0 means any increase is good.
def wait_s3_object_counts_state_changed(context):
    minio_client = _get_minio_client(context)
    for (bucket, filename_pattern, change_type) in context.table:
        timeout = int(context.config.userdata["timeout_s"])
        found = False
        while not found:
            if timeout < 0:
                raise Exception(f'{bucket} s3 object count did not increase in time')
            current_timestamp_dict = _get_bucket_object_state(minio_client, bucket, filename_pattern)
            stored_timestamp_dict = _cache_bucket_object_state(context, bucket, filename_pattern)
            difference = len(current_timestamp_dict) - len(stored_timestamp_dict)
            if change_type == "count":
               if difference != 0:
                   found = True
                   print(f'{bucket} s3 object count increased by {difference}')
            elif change_type == "timestamp":
                for object_name, stored_timestamp in stored_timestamp_dict.items():
                    current_timestamp = current_timestamp_dict.get(object_name)
                    if current_timestamp > stored_timestamp:
                        # If any timestamp is updated, we consider it as an update.
                        found = True
                        print(f'Object {object_name} timestamp updated to {current_timestamp}')
                        break
            else:
                raise Exception(f'Invalid change type {change_type}')
            time.sleep(1)
            timeout -= 1
        _cache_bucket_object_state(context, bucket, filename_pattern, current_timestamp_dict)

def _cache_bucket_object_state(context, bucket, pattern, state_dict: dict=None):
    key = f'{bucket}_{pattern}'
    if state_dict is not None:
        context.state["storage"][key] = state_dict
    else:
        return context.state["storage"][key]

def _get_minio_client(context):
    global minio_client
    if minio_client is not None:
        return minio_client
    minio_client = Minio(context.config.userdata["s3_connection_url"],
        access_key=get_secret("s3_access_key", context = context),
        secret_key=get_secret("s3_secret_key", context = context),
        secure=False
    )
    return minio_client

def push_questionnaire_response_data(context):
    key_id, value_id  = _get_kafka_topic_key_value_ids(context, "questionnaire_response")
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_armt_access_token(context)}',
        'accept': 'application/vnd.kafka.v2+json, application/vnd.kafka+json; q=0.9, application/json; q=0.8'
    }
    project_name = context.cache["project_json"]["projectName"]
    user_id = context.cache["test_subject_id"]
    source_id = context.cache["armt_project_source_json"]["sourceId"]
    assert context.text is not None
    answers = json.loads(context.text)
    # Send data with random timestamp so that different files will be created
    # on the output storage (not removed by data deduplication).
    timestamp = _random_timestamp()
    data = {
        "key_schema_id": key_id,
        "value_schema_id": value_id,
        "records": [
            {
                "value": {
                    "time":  timestamp,
                    "timeCompleted": timestamp,
                    "timeNotification": 0,
                    "name": "QuestionnaireName",
                    "version": "1",
                    "answers": answers
                },
                "key": {
                    "projectId": {"string": project_name},
                    "userId": user_id,
                    "sourceId": source_id
                }
            }
        ]
    }
    response = requests.post(f'{context.config.userdata["url"]}/kafka/topics/questionnaire_response', headers=headers, data=json.dumps(data))
    assert response.status_code == 200

def _random_timestamp():
    days = random.randint(1, 100000)
    now = datetime.datetime.now()
    random_date = now - datetime.timedelta(days=days)
    return random_date.timestamp()

def _get_kafka_topic_key_value_ids(context, topic_name):
    response = requests.get(f'{context.config.userdata["url"]}/schema/subjects/{topic_name}-key/versions/latest')
    assert response.status_code == 200
    key_id = response.json()['id']
    response = requests.get(f'{context.config.userdata["url"]}/schema/subjects/{topic_name}-value/versions/latest')
    assert response.status_code == 200
    value_id = response.json()['id']
    return key_id, value_id

def register_fitbit_user(context):
    # Here the study manager enrolls a subject for Fitbit data collection.
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    data = {
        "userId": context.cache["fitbit_user_json"]["id"],
        "persistent": True
    }
    response = requests.post(f'{context.config.userdata["url"]}/rest-sources/backend/registrations', headers=headers, data=json.dumps(data))
    assert response.status_code == 200 or response.status_code == 201
    registration_json = response.json()
    assert registration_json is not None
    registration_token = registration_json['token']

    # At this moment the user is redirected to Fitbi7t authorization page (or clicks this link sent via email).
    # The user is redirected to RADAR-base frontend URL callback with the 'code' (of authorization_code grant).
    # The frontend redirects to the backend with the 'code' and the 'state' (registration_token).
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json'
    }
    data = {
        "code": "mycode"
    }
    response = requests.post(f'{context.config.userdata["url"]}/rest-sources/backend/registrations/{registration_token}/authorize', headers=headers, data=json.dumps(data))
    assert response.status_code == 200 or response.status_code == 201

    # The rest source backend will exchange the code for an access token and refresh token.
    # This happens in the back chanel and is not part of the e2e test.
    # the call will be to $fitbit_url/oauth2/token?client_id=$client;client_secret=some_secret;redirect_uri=$protocol://$host/rest-sources/backend/users:new;grant_type=authorization_code
    # And body: '{"code":"'$code'"}'
    # The response is provided by the mock server.

def check_fitbit_user_exists(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    project_name = context.cache["project_json"]["projectName"]
    user_id = context.cache["test_subject_id"]
    response = requests.get(f'{context.config.userdata["url"]}/rest-sources/backend/users?project-id={project_name}&authorized=all', headers=headers)
    assert response.status_code == 200
    users_json = response.json()
    user_json = next((item for item in users_json['users'] if item["userId"] == user_id), None)
    if not context.config.userdata["dev_mode"] and user_json is not None:
        raise Exception('Fitbit user already exists')
    elif user_json is not None:
        context.cache["fitbit_user_json"] = user_json
    else:
        _create_fitbit_user(context)

def _create_fitbit_user(context):
    if context.cache["fitbit_user_json"] is not None:
        return
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    project_name = context.cache["project_json"]["projectName"]
    user_id = context.cache["test_subject_id"]
    data = {
        "projectId": project_name,
        "userId": user_id,
        "startDate": "2020-11-18T23:00:00.000Z",
        "endDate": "3000-11-29T23:00:00.000Z",
        "sourceType": "FitBit"
    }
    response = requests.post(f'{context.config.userdata["url"]}/rest-sources/backend/users', headers=headers, data=json.dumps(data))
    assert response.status_code == 200 or response.status_code == 201
    user_json = response.json()
    assert user_json is not None
    context.cache["fitbit_user_json"] = user_json

def check_service_states(context):
    for (service_name, state) in context.table:
        stream = os.popen(f'kubectl get pods | grep {service_name} | grep {state}')
        output = stream.read()
        if service_name not in output:
            raise Exception(f'{service_name} is not running')

def count_rows_in_postgresql(context):
    for (service, database, table) in context.table:
        count = _get_postgres_table_state(context, service, database, table)
        _cache_database_table_state(context, service, database, table, count)

def _get_postgres_table_state(context, service, database, table):
    password = get_secret('data_dashboard_db_password', context = context)
    stream = os.popen(f'kubectl exec {service} -c postgresql -- psql postgresql://postgres:{password}@localhost/{database} -t -c "SELECT COUNT(*) FROM {table}"')
    try:
        count = int(stream.read().strip())
    except ValueError:
        count = 0
    return count

def wait_for_postgresql_table_state_change(context):
    for (service, database, table) in context.table:
        timeout = int(context.config.userdata["timeout_s"])
        while True:
            if timeout < 0:
                raise Exception(f'{service} {database} {table} table state did not change in time')
            current_count = _get_postgres_table_state(context, service, database, table)
            stored_count = _cache_database_table_state(context, service, database, table)
            if stored_count is not None and current_count > stored_count:
                break
            time.sleep(1)
            timeout -= 1
        _cache_database_table_state(context, service, database, table, current_count)

def _cache_database_table_state(context, service, database, table, count: int=None):
    key = f'{service}_{database}_{table}'
    if count is not None:
        context.state["database"][key] = count
    else:
        return context.state["database"][key]

def wait_for_postgresql_record_count(context):
    for (service, database, table, count) in context.table:
        timeout = int(context.config.userdata["timeout_s"])
        while True:
            if timeout < 0:
                raise Exception(f'{service} {database} {table} row counts did not change in time')
            current_count = _get_postgres_table_state(context, service, database, table)
            if current_count == int(count):
                break
            time.sleep(1)
            timeout -= 1
        current_count = _get_postgres_table_state(context, service, database, table)
        assert current_count == int(count)