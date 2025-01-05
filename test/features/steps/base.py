import base64

import yaml
from functools import reduce
import operator
import requests
import json

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

def check_armt_source_type_exists(context):
    source_types = get_source_types(context)
    armt_source_type_json = next((item for item in source_types if item["name"] == "RADAR_aRMT"), None)
    if context.config.userdata["dev_mode"] and armt_source_type_json is not None:
        context.cache["armt_source_id"] = armt_source_type_json['name']
    elif armt_source_type_json is not None:
        raise Exception('aRMT source type already exists')

def create_armt_source_type(context):
    if context.cache["armt_source_id"] is not None:
        return
    armt_source_id = "RADAR_aRMT"
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
        "name": armt_source_id
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/api/source-types', headers=headers, data=json.dumps(data))
    # 409 means that the armt source already exists.
    assert response.status_code == 200 or response.status_code == 409
    if response.status_code == 200:
        armt_source_id = response.json()['name']
        assert armt_source_id is not None
    context.cache["armt_source_id"] = armt_source_id

def check_organization_exists(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application /json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/organizations', headers=headers)
    assert response.status_code == 200
    test_organization_json = next((item for item in response.json() if item["name"] == context.config.userdata["org_name"]), None)
    if context.config.userdata["dev_mode"] and test_organization_json is not None:
        context.cache["organization_json"] = test_organization_json['name']
    elif test_organization_json is not None:
        raise Exception('Test organization already exists')

def create_organization(context):
    if context.cache["organization_json"] is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    data = {
        "name": context.config.userdata["org_name"],
        "description": "TEST",
        "location": "TEST"
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/api/organizations', headers=headers, data=json.dumps(data))
    # # 409 means that the armt source already exists.
    # assert response.status_code == 200 or response.status_code == 409
    # 409 means that the organization already exists
    assert response.status_code == 200 or response.status_code == 409
    if response.status_code == 200:
        test_organization_json = response.json()
        assert test_organization_json is not None
    context.cache["organization_json"] = test_organization_json

def check_project_exists(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/projects', headers=headers)
    assert response.status_code == 200
    test_project_json = next((item for item in response.json() if item["projectName"] == context.config.userdata["project_name"]), None)
    if context.config.userdata["dev_mode"] and test_project_json is not None:
        context.cache["project_json"] = test_project_json
    elif test_project_json is not None:
        raise Exception('Test project already exists')

def create_project(context):
    if context.cache["project_json"] is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    data = {
        "organization": context.cache["organization_json"],
        "projectName": context.config.userdata["project_name"],
        "description": "TEST",
        "location": "",
        "sourceTypes": [context.cache["armt_source_id"]],
        "startDate": None,
        "endDate": None
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/api/projects', headers=headers, data=json.dumps(data))
    assert response.status_code == 200
    test_project_json = response.json()
    assert test_project_json['projectName'] == context.config.userdata["project_name"]
    context.cache["project_json"] = test_project_json

def check_armt_project_source_exists(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/projects/{context.config.userdata["project_name"]}/sources', headers=headers)
    assert response.status_code == 200
    armt_project_source_json = next((item for item in response.json() if item["sourceType"]["name"] == "RADAR_aRMT"), None)
    if context.config.userdata["dev_mode"] and armt_project_source_json is not None:
        context.cache["armt_project_source_id"] = armt_project_source_json['sourceId']
    elif armt_project_source_json is not None:
        raise Exception('aRMT project source already exists')

def create_armt_project_source(context):
    if context.cache["armt_project_source_id"] is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    data = {
        "sourceName": "aRMT-test-source-TEST",
        "assigned": False,
        "sourceType": context.cache["armt_source_id"],
        "project": context.cache["organization_json"]
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/api/sources', headers=headers, data=json.dumps(data))
    assert response.status_code == 200
    armt_project_source_json = response.json()
    assert armt_project_source_json['sourceId'] is not None
    context.cache["armt_project_source_id"] = armt_project_source_json['sourceId']

def check_subject_exists(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    response = requests.get(f'{context.config.userdata["url"]}/managementportal/api/projects/{context.config.userdata["project_name"]}/subjects', headers=headers)
    assert response.status_code == 200
    test_subject_json = next((item for item in response.json() if item["externalId"] == context.config.userdata["subject_external_id"]), None)
    if context.config.userdata["dev_mode"] and test_subject_json is not None:
        context.cache["test_subject_id"] = test_subject_json['login']
    elif test_subject_json is not None:
        raise Exception('Test subject already exists')

def create_subject(context):
    if context.cache["test_subject_id"] is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token(context)}'
    }
    data = {
        "project": context.cache["organization_json"],
        "sources": [context.cache["armt_project_source_id"]],
        "status": 1,
        "externalId": context.config.userdata["subject_external_id"],
        "group": None
    }
    response = requests.post(f'{context.config.userdata["url"]}/managementportal/api/subjects', headers=headers, data=json.dumps(data))
    assert response.status_code == 200
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
