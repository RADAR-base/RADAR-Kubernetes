from behave import *

from base import create_armt_source_type, create_organization, create_project, create_subject, \
    create_armt_project_source, get_armt_meta_token, get_armt_refresh_token, get_armt_access_token, \
    check_armt_source_type_exists, check_organization_exists, check_project_exists, check_subject_exists, \
    check_armt_project_source_exists, get_current_s3_object_state, wait_s3_object_counts_state_changed, \
    push_questionnaire_response_data


@given('creation of an aRMT source type named "{source_type}"')
def step_impl(context, source_type):
    check_armt_source_type_exists(context, source_type)
    create_armt_source_type(context, source_type)

@given('creation of an organization named "{organization}"')
def step_impl(context, organization):
    check_organization_exists(context, organization)
    create_organization(context, organization)

@given('creation of a project named "{project}"')
def step_impl(context, project):
    check_project_exists(context, project)
    create_project(context, project)

@given('creation of a subject named "{subject}"')
def step_impl(context, subject):
    check_subject_exists(context, subject)
    create_subject(context, subject)

@given('creation of an aRMT project source named "{source}"')
def step_impl(context, source):
    check_armt_project_source_exists(context, source)
    create_armt_project_source(context, source)

@given('the aRMT application has retrieved an access token')
def step_impl(context):
    get_armt_meta_token(context)
    get_armt_refresh_token(context)
    get_armt_access_token(context)

@given('the state of objects in the s3 storage')
def step_impl(context):
    get_current_s3_object_state(context)

@when('the aRMT application sends questionnaire_response data')
def step_impl(context):
    push_questionnaire_response_data(context)

@then('the state of objects in the s3 storage changes')
def step_impl(context):
    wait_s3_object_counts_state_changed(context)


