from behave import *

from base import wait_for_postgresql_record_count


@then('Fitbit records are present in the database')
def step_impl(context):
    wait_for_postgresql_record_count(context)
