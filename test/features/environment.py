from behave import use_fixture
from fixtures import s3

def before_all(context):
    context.cache = {
        "management_portal_token": None,
        "armt_source_id": None,
        "organization_json": None,
        "project_json": None,
        "armt_project_source_id": None,
        "test_subject_id": None,
        "secrets": {},
        "armt_meta_token": None,
        "armt_refresh_token": None,
        "armt_access_token": None
    }

def before_tag(context, tag):
    if tag == "fixture.s3":
        use_fixture(s3, context)
