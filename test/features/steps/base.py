import os
import yaml
from functools import reduce  # forward compatibility for Python 3
import operator
import requests

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
    secrets = {}


def get_secret(*path_elements):
    if Cache.secrets is None or len(Cache.secrets) == 0:
        secrets_file = 'etc/secrets.yaml'
        with open(secrets_file, 'r') as file:
            Cache.secrets = yaml.safe_load(file)
    return reduce(operator.getitem, path_elements, Cache.secrets)

def format_url(path):
    return f'{TestConfig.protocol}://{TestConfig.host}:{TestConfig.port}/{path}'

def request_mp_token():
    if Cache.management_portal_token is not None:
        Cache.management_portal_token
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