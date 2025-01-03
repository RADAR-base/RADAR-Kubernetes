import requests

from behave import *
from base import TestConfig
from base import Cache
from base import get_secret
from base import format_url
from base import request_mp_token

dev = TestConfig.dev_mode
mp_admin_password = get_secret('management_portal', 'managementportal', 'common_admin_password')

@given('management portal is running')
def step_impl(context):
    time = 10
    while True:
        if time < 0:
            raise Exception('Management portal did not start in time')
        response = requests.get(format_url('managementportal'))
        if response.status_code == 200:
            break
        time.sleep(1)
        time -= 1

@then('the management portal token can be requested')
def step_impl(context):
    request_mp_token()

@given('retrieval of management portal token')
def step_impl(context):
    request_mp_token()

@given('the aRMT source type does not exist')
def step_impl(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {Cache.management_portal_token}'
    }
    response = requests.get(format_url('managementportal/api/source-types'), headers=headers)
    assert response.status_code == 200
    armt_source_type_json = next((item for item in response.json() if item["name"] == "RADAR_aRMT"), None)
    if dev and armt_source_type_json is not None:
        Cache.armt_source_id = armt_source_type_json['name']
    else:
        raise Exception('aRMT source type already exists')

@then('the aRMT source type can be created')
def step_impl(context):
    if dev and Cache.armt_source_id is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {Cache.management_portal_token}'
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
    android_source_type_json = response.json()
    assert android_source_type_json['name'] is not None
    Cache.armt_source_id = android_source_type_json['name']

@given('the test organization does not exist')
def step_impl(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application /json',
        'Authorization': f'Bearer {Cache.management_portal_token}'
    }
    response = requests.get(format_url('managementportal/api/organizations'), headers=headers)
    assert response.status_code == 200
    test_organization_json = next((item for item in response.json() if item["name"] == TestConfig.org_name), None)
    # This should never happen during tests on a fresh deployment
    if dev and test_organization_json is not None:
        Cache.organization_json = test_organization_json['name']
    else:
        raise Exception('Test organization already exists')

@then('the test organization should be created')
def step_impl(context):
    if dev and Cache.organization_json is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {Cache.management_portal_token}'
    }
    data = {
        "name": TestConfig.org_name,
        "description": "TEST",
        "location": "TEST"
    }
    response = requests.post(format_url('managementportal/api/organizations'), headers=headers, data=json.dumps(data))
    assert response.status_code == 200
    test_organization_json = response.json()
    assert test_organization_json['name'] is not None
    Cache.organization_json = test_organization_json

@given('the test project does not exist')
def step_impl(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {Cache.management_portal_token}'
    }
    response = requests.get(format_url('managementportal/api/projects'), headers=headers)
    assert response.status_code == 200
    test_project_json = next((item for item in response.json() if item["projectName"] == TestConfig.project_name), None)
    if dev and test_project_json is not None:
        Cache.project_json = test_project_json
    else:
        raise Exception('Test project already exists')

@then('the test project should be created')
def step_impl(context):
    if dev and Cache.project_json is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {Cache.management_portal_token}'
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

@given('the aRMT project source does not exist')
def step_impl(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {Cache.management_portal_token}'
    }
    response = requests.get(format_url(f'managementportal/api/projects/{TestConfig.project_name}/sources'), headers=headers)
    assert response.status_code == 200
    armt_project_source_json = next((item for item in response.json() if item["sourceType"]["name"] == "RADAR_aRMT"), None)
    if dev and armt_project_source_json is not None:
        Cache.armt_project_source_id = armt_project_source_json['sourceId']
    else:
        raise Exception('aRMT project source already exists')

@then('the aRMT project source should be created')
def step_impl(context):
    if dev and Cache.armt_project_source_id is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {Cache.management_portal_token}'
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

@given('the test subject does not exist')
def step_impl(context):
    headers = {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {Cache.management_portal_token}'
    }
    response = requests.get(format_url(f'managementportal/api/projects/{TestConfig.project_name}/subjects'), headers=headers)
    assert response.status_code == 200
    test_subject_json = next((item for item in response.json() if item["externalId"] == TestConfig.subject_external_id), None)
    if dev and test_subject_json is not None:
        Cache.test_subject_id = test_subject_json['login']
    else:
        raise Exception('Test subject already exists')

@then('the test subject should be created')
def step_impl(context):
    if dev and Cache.test_subject_id is not None:
        return
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {Cache.management_portal_token}'
    }
    data = {
        "project": Cache.organization_json,
        "sources": [Cache.armt_project_source_id],
        "status": 1,
        "externalId": "TEST",
        "group": None
    }
    response = requests.post(format_url('managementportal/api/subjects'), headers=headers, data=json.dumps(data))
    assert response.status_code == 200
    test_subject_json = response.json()
    assert test_subject_json['login'] is not None
