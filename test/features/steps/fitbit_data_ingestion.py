from behave import *

from base import register_fitbit_user, check_fitbit_user_exists, wait_s3_object_counts_state_changed


@given('registration of the subject with Fitbit authorization service')
def step_impl(context):
    check_fitbit_user_exists(context)
    register_fitbit_user(context)

@then('Fitbit connector will download data and the state of objects in the s3 storage changes')
def step_impl(context):
    wait_s3_object_counts_state_changed(context)
