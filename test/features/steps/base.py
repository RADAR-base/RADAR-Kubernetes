import base64

import yaml
from functools import reduce
import operator
import requests
import json
from minio import Minio
import time

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
    assert response.status_code == 200 or response.status_code == 409
    if response.status_code == 200:
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
    assert response.status_code == 200 or response.status_code == 409
    if response.status_code == 200:
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
    # 200 means that the project already exists
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

def get_current_s3_object_counts(context):
    minio_client = _get_minio_client(context)
    for (bucket, filename_pattern) in context.table:
        count = get_current_s3_object_count(minio_client, bucket, filename_pattern)
        _get_or_set_count(context, bucket, filename_pattern, count)

def get_current_s3_object_count(minio_client, bucket, pattern):
    objects = minio_client.list_objects(
        bucket_name=bucket,
        recursive=True
    )
    filtered_objects = [ object for object in objects if pattern in object.object_name ]
    return len(filtered_objects)

def wait_s3_object_counts_increased(context, number):
    minio_client = _get_minio_client(context)
    desired_change = int(number)
    for (bucket, filename_pattern) in context.table:
        timeout = int(context.config.userdata["timeout_s"])
        while True:
            if timeout < 0:
                raise Exception(f'{bucket} s3 object count did not increase in time')
            current_count = get_current_s3_object_count(minio_client, bucket, filename_pattern)
            stored_count = _get_or_set_count(context, bucket, filename_pattern)
            difference = current_count - stored_count
            if difference == desired_change:
                break
            if difference != 0 and (difference < desired_change or difference > desired_change):
                raise Exception(f'S3 object count changed but not by desired number ({desired_change})')
            time.sleep(1)
            timeout -= 1
        _get_or_set_count(context, bucket, filename_pattern, current_count)

def _get_minio_client(context):
    global minio_client
    if minio_client is not None:
        return minio_client
    minio_client = Minio(context.config.userdata["s3_connection_url"],
        access_key=context.config.userdata["s3_key"],
        secret_key=context.config.userdata["s3_secret"],
        secure=False
    )
    return minio_client

def _get_or_set_count(context, bucket, pattern, count: int=None):
    key = f'{bucket}{pattern}'
    if count is not None:
        context.counts[key] = count
    else:
        return context.counts[key]

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
    data = {
        "key_schema_id": key_id,
        "value_schema_id": value_id,
        "records": [
            {
                "value": {
                    "time": 1707657463.9095526,
                    "timeCompleted": 1707657463.9095526,
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

def _get_kafka_topic_key_value_ids(context, topic_name):
    response = requests.get(f'{context.config.userdata["url"]}/schema/subjects/{topic_name}-key/versions/latest')
    assert response.status_code == 200
    key_id = response.json()['id']
    response = requests.get(f'{context.config.userdata["url"]}/schema/subjects/{topic_name}-value/versions/latest')
    assert response.status_code == 200
    value_id = response.json()['id']
    return key_id, value_id
