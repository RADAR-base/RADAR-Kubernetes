from behave import *

from base import create_armt_source_type, create_organization, create_project, create_subject, \
    create_armt_project_source, get_armt_meta_token, get_armt_refresh_token, get_armt_access_token, \
    check_armt_source_type_exists, check_organization_exists, check_project_exists, check_subject_exists, \
    check_armt_project_source_exists


@given('creation of the aRMT source type')
def step_impl(context):
    check_armt_source_type_exists(context)
    create_armt_source_type(context)

@given('creation of the test organization')
def step_impl(context):
    check_organization_exists(context)
    create_organization(context)

@given('creation of the test project')
def step_impl(context):
    check_project_exists(context)
    create_project(context)

@given('creation of the test subject')
def step_impl(context):
    check_subject_exists(context)
    create_subject(context)

@given('creation of the aRMT project source')
def step_impl(context):
    check_armt_project_source_exists(context)
    create_armt_project_source(context)

@given('the aRMT application has retrieved an access token')
def step_impl(context):
    get_armt_meta_token(context)
    get_armt_refresh_token(context)
    get_armt_access_token(context)

@then('true')
def step_impl(context):
    assert True

