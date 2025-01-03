import base64
import os
import yaml
from functools import reduce  # forward compatibility for Python 3
import operator
import requests
import json

class TestConfig:
    protocol = os.environ.get('TEST_PROTOCOL', 'http')
    host = os.environ.get('TEST_HOST', 'localhost')
    port = os.environ.get('TEST_PORT', 80)
    org_name = os.environ.get('TEST_ORGANIZATION', 'MAIN')
    project_name = os.environ.get('TEST_PROJECT', 'test')
    subject_external_id = os.environ.get('TEST_SUBJECT', 'test_user')
    # TODO
    # dev_mode = os.environ.get('DEV_MODE', False)
    dev_mode = os.environ.get('DEV_MODE', True)

class Cache:
    management_portal_token = None
    armt_source_id = None
    organization_json = None
    project_json = None
    armt_project_source_id = None
    test_subject_id = None
    secrets = dict(),
    armt_meta_token = None,
    armt_refresh_token = None,
    armt_access_token = None


def wait_for_management_portal():
    time = 10
    while True:
        if time < 0:
            raise Exception('Management portal did not start in time')
        response = requests.get(format_url('managementportal'))
        if response.status_code == 200:
            break
        time.sleep(1)
        time -= 1

def get_secret(*path_elements):
    if Cache.secrets is None or len(Cache.secrets[0].keys()) == 0:
        secrets_file = 'etc/secrets.yaml'
        with open(secrets_file, 'r') as file:
            Cache.secrets = yaml.safe_load(file)
    return reduce(operator.getitem, path_elements, Cache.secrets)

def format_url(path):
    return f'{TestConfig.protocol}://{TestConfig.host}:{TestConfig.port}/{path}'

def get_mp_token():
    if Cache.management_portal_token is not None:
        return Cache.management_portal_token
    mp_admin_password = get_secret('management_portal', 'managementportal', 'common_admin_password')
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
    response = requests.post(url=format_url('managementportal/oauth/token'), headers=headers, data=data)
    assert response.status_code == 200
    management_portal_json = response.json()
    token = management_portal_json['access_token']
    assert token is not None
    Cache.management_portal_token = token
    return token

def get_source_types():
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    response = requests.get(format_url('managementportal/api/source-types'), headers=headers)
    assert response.status_code == 200
    return response.json()

def check_armt_source_type_exists():
    source_types = get_source_types()
    armt_source_type_json = next((item for item in source_types if item["name"] == "RADAR_aRMT"), None)
    if TestConfig.dev_mode and armt_source_type_json is not None:
        Cache.armt_source_id = armt_source_type_json['name']
    elif armt_source_type_json is not None:
        raise Exception('aRMT source type already exists')

def create_armt_source_type():
    if Cache.armt_source_id is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    data = {
        "appProvider": "org.radarcns.application.ApplicationServiceProvider",
        "producer": "RADAR",
        "model": "aRMT",
        "catalogVersion": "1.5.0",
        "sourceTypeScope": "ACTIVE",
        "canRegisterDynamically": True,
        "name": "RADAR_aRMT"
    }
    response = requests.post(format_url('managementportal/api/source-types'), headers=headers, data=json.dumps(data))
    assert response.status_code == 200
    armt_source_id = response.json()['name']
    assert armt_source_id is not None
    Cache.armt_source_id = armt_source_id

def check_organization_exists():
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application /json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    response = requests.get(format_url('managementportal/api/organizations'), headers=headers)
    assert response.status_code == 200
    test_organization_json = next((item for item in response.json() if item["name"] == TestConfig.org_name), None)
    if TestConfig.dev_mode and test_organization_json is not None:
        Cache.organization_json = test_organization_json['name']
    elif test_organization_json is not None:
        raise Exception('Test organization already exists')

def create_organization():
    if Cache.organization_json is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    data = {
        "name": TestConfig.org_name,
        "description": "TEST",
        "location": "TEST"
    }
    response = requests.post(format_url('managementportal/api/organizations'), headers=headers, data=json.dumps(data))
    assert response.status_code == 200
    test_organization_json = response.json()
    assert test_organization_json is not None
    Cache.organization_json = test_organization_json

def check_project_exists():
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    response = requests.get(format_url('managementportal/api/projects'), headers=headers)
    assert response.status_code == 200
    test_project_json = next((item for item in response.json() if item["projectName"] == TestConfig.project_name), None)
    if TestConfig.dev_mode and test_project_json is not None:
        Cache.project_json = test_project_json
    elif test_project_json is not None:
        raise Exception('Test project already exists')

def create_project():
    if Cache.project_json is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    data = {
        "organization": Cache.organization_json,
        "projectName": TestConfig.project_name,
        "description": "TEST",
        "location": "",
        "sourceTypes": [Cache.armt_source_id],
        "startDate": None,
        "endDate": None
    }
    response = requests.post(format_url('managementportal/api/projects'), headers=headers, data=json.dumps(data))
    assert response.status_code == 200
    test_project_json = response.json()
    assert test_project_json['projectName'] == TestConfig.project_name
    Cache.project_json = test_project_json

def check_armt_project_source_exists():
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    response = requests.get(format_url(f'managementportal/api/projects/{TestConfig.project_name}/sources'), headers=headers)
    assert response.status_code == 200
    armt_project_source_json = next((item for item in response.json() if item["sourceType"]["name"] == "RADAR_aRMT"), None)
    if TestConfig.dev_mode and armt_project_source_json is not None:
        Cache.armt_project_source_id = armt_project_source_json['sourceId']
    elif armt_project_source_json is not None:
        raise Exception('aRMT project source already exists')

def create_armt_project_source():
    if Cache.armt_project_source_id is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    data = {
        "sourceName": "aRMT-test-source-TEST",
        "assigned": False,
        "sourceType": Cache.armt_source_id,
        "project": Cache.organization_json
    }
    response = requests.post(format_url('managementportal/api/sources'), headers=headers, data=json.dumps(data))
    assert response.status_code == 200
    armt_project_source_json = response.json()
    assert armt_project_source_json['sourceId'] is not None
    Cache.armt_project_source_id = armt_project_source_json['sourceId']

def check_subject_exists():
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    response = requests.get(format_url(f'managementportal/api/projects/{TestConfig.project_name}/subjects'), headers=headers)
    assert response.status_code == 200
    test_subject_json = next((item for item in response.json() if item["externalId"] == TestConfig.subject_external_id), None)
    if TestConfig.dev_mode and test_subject_json is not None:
        Cache.test_subject_id = test_subject_json['login']
    elif test_subject_json is not None:
        raise Exception('Test subject already exists')

def create_subject():
    if Cache.test_subject_id is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    data = {
        "project": Cache.organization_json,
        "sources": [Cache.armt_project_source_id],
        "status": 1,
        "externalId": TestConfig.subject_external_id,
        "group": None
    }
    response = requests.post(format_url('managementportal/api/subjects'), headers=headers, data=json.dumps(data))
    assert response.status_code == 200
    test_subject_json = response.json()
    assert test_subject_json['login'] is not None
    Cache.test_subject_id = test_subject_json['login']

def get_armt_meta_token():
    if Cache.armt_meta_token is not None:
        return Cache.armt_meta_token
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {get_mp_token()}'
    }
    data = {
        'clientId': 'aRMT',
        'login': Cache.test_subject_id,
        'persistent': True
    }
    response = requests.post(format_url('managementportal/api/oauth-clients/pair'), headers=headers, data=json.dumps(data))
    assert response.status_code == 200
    meta_token = response.json()['tokenName']
    assert meta_token is not None
    Cache.armt_meta_token = meta_token
    return meta_token


def get_armt_refresh_token():
    if Cache.armt_refresh_token is not None:
        return Cache.armt_refresh_token
    headers = {
        'Authorization': f'Bearer {get_mp_token()}'
    }
    response = requests.get(format_url(f'managementportal/api/meta-token/{get_armt_meta_token()}'), headers=headers)
    assert response.status_code == 200
    refresh_token = response.json()['refreshToken']
    assert refresh_token is not None
    Cache.armt_refresh_token = refresh_token
    return refresh_token


def get_armt_access_token():
    if Cache.armt_access_token is not None:
        return Cache.armt_access_token
    headers = {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Authorization': f'Basic {base64.b64encode(b"aRMT:").decode("utf-8")}'
    }
    data = {
        'grant_type': 'refresh_token',
        'refresh_token': get_armt_refresh_token()
    }
    response = requests.post(format_url('managementportal/oauth/token'), headers=headers, data=data)
    assert response.status_code == 200
    access_token = response.json()['access_token']
    assert access_token is not None
    Cache.armt_access_token = access_token
    return access_token