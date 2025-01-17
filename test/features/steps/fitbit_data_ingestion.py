from behave import *

from base import create_armt_source_type, create_organization, create_project, create_subject, \
    create_armt_project_source, get_armt_meta_token, get_armt_refresh_token, get_armt_access_token, \
    check_armt_source_type_exists, check_organization_exists, check_project_exists, check_subject_exists, \
    check_armt_project_source_exists, get_current_s3_object_counts, wait_s3_object_counts_increased_or_updated, \
    push_questionnaire_response_data, register_fitbit_user, check_fitbit_user_exists


@given('registration of the subject with Fitbit authorization service')
def step_impl(context):
    check_fitbit_user_exists(context)
    register_fitbit_user(context)

@then('Fitbit connector will pick start downloading the data and the object counts in the s3 storage for files have increased')
def step_impl(context):
    wait_s3_object_counts_increased_or_updated(context)


