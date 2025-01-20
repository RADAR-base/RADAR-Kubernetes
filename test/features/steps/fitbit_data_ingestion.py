from behave import *

from base import register_fitbit_user, check_fitbit_user_exists, wait_s3_object_counts_increased_or_updated


@given('registration of the subject with Fitbit authorization service')
def step_impl(context):
    check_fitbit_user_exists(context)
    register_fitbit_user(context)

@then('Fitbit connector will pick start downloading the data and the object counts in the s3 storage for files have increased')
def step_impl(context):
    wait_s3_object_counts_increased_or_updated(context)


