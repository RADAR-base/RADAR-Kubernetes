from behave import *
from base import check_armt_source_type_exists, create_armt_source_type, check_organization_exists, \
    create_organization, get_mp_token, check_project_exists, create_project, check_armt_project_source_exists, create_armt_project_source, \
    check_subject_exists, create_subject, wait_for_management_portal


@given('management portal is running')
def step_impl(context):
    wait_for_management_portal(context)

@then('the management portal token can be requested')
def step_impl(context):
    get_mp_token(context)

@given('retrieval of management portal token')
def step_impl(context):
    get_mp_token(context)

# @given('the aRMT source type does not exist')
# def step_impl(context):
#     check_armt_source_type_exists(context)

@then('the aRMT source type can be created')
def step_impl(context):
    check_armt_source_type_exists(context)
    create_armt_source_type(context)

# @given('the test organization does not exist')
# def step_impl(context):
#     check_organization_exists(context)

@then('the test organization should be created')
def step_impl(context):
    check_organization_exists(context)
    create_organization(context)

# @given('the test project does not exist')
# def step_impl(context):
#     check_project_exists(context)

@then('the test project should be created')
def step_impl(context):
    check_project_exists(context)
    create_project(context)

# @given('the aRMT project source does not exist')
# def step_impl(context):
#     check_armt_project_source_exists(context)

@then('the aRMT project source should be created')
def step_impl(context):
    check_armt_project_source_exists(context)
    create_armt_project_source(context)

# @given('the test subject does not exist')
# def step_impl(context):
#     check_subject_exists(context)

@then('the test subject should be created')
def step_impl(context):
    check_subject_exists(context)
    create_subject(context)
