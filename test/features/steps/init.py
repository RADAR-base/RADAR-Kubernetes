from behave import *
from test.features.steps.base import check_armt_source_type_exists, create_armt_source_type, check_organization_exists, \
    create_organization, get_mp_token, check_project_exists, create_project, check_armt_project_source_exists, create_armt_project_source, \
    check_subject_exists, create_subject, wait_for_management_portal


@given('management portal is running')
def step_impl(context):
    wait_for_management_portal()

@then('the management portal token can be requested')
def step_impl(context):
    get_mp_token()

@given('retrieval of management portal token')
def step_impl(context):
    get_mp_token()

@given('the aRMT source type does not exist')
def step_impl(context):
    check_armt_source_type_exists()

@then('the aRMT source type can be created')
def step_impl(context):
    create_armt_source_type()

@given('the test organization does not exist')
def step_impl(context):
    check_organization_exists()

@then('the test organization should be created')
def step_impl(context):
    create_organization()

@given('the test project does not exist')
def step_impl(context):
    check_project_exists()

@then('the test project should be created')
def step_impl(context):
    create_project()

@given('the aRMT project source does not exist')
def step_impl(context):
    check_armt_project_source_exists()

@then('the aRMT project source should be created')
def step_impl(context):
    create_armt_project_source()

@given('the test subject does not exist')
def step_impl(context):
    check_subject_exists()

@then('the test subject should be created')
def step_impl(context):
    create_subject()
