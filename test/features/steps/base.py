import os

class TestConfig:
    protocol = os.environ.get('TEST_PROTOCOL', 'http')
    host = os.environ.get('TEST_HOST', 'localhost')
    port = os.environ.get('TEST_PORT', 80)
    org_name = os.environ.get('TEST_ORGANIZATION', 'MAIN')
    project_name = os.environ.get('TEST_PROJECT', 'test')
    subject_external_id = os.environ.get('TEST_SUBJECT', 'test_user')

class Cache:
    management_portal_token = None
    armt_source_id = None
    organization_json = None
    project_json = None
    armt_project_source_id = None
    test_subject_id = None