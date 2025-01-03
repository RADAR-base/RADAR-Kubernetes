from behave import *

from test.features.steps.base import create_armt_source_type, create_organization, create_project, create_subject, \
    create_armt_project_source, get_armt_meta_token, get_armt_refresh_token, get_armt_access_token, \
    check_armt_source_type_exists, check_organization_exists, check_project_exists, check_subject_exists, \
    check_armt_project_source_exists


@given('creation of the aRMT source type')
def step_impl(context):
    check_armt_source_type_exists()
    create_armt_source_type()

@given('creation of the test organization')
def step_impl(context):
    check_organization_exists()
    create_organization()

@given('creation of the test project')
def step_impl(context):
    check_project_exists()
    create_project()

@given('creation of the test subject')
def step_impl(context):
    check_subject_exists()
    create_subject()

@given('creation of the aRMT project source')
def step_impl(context):
    check_armt_project_source_exists()
    create_armt_project_source()

@given('the aRMT application has retrieved an access token')
def step_impl(context):
    get_armt_meta_token()
    get_armt_refresh_token()
    get_armt_access_token()

@then('true')
def step_impl(context):
    assert True

