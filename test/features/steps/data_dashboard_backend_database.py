from behave import *

from base import check_service_states, count_rows_in_postgresql, wait_for_postgresql_table_state_change


@given('these service states')
def step_impl(context):
    check_service_states(context)

@given('the number of rows in the database')
def step_impl(context):
    count_rows_in_postgresql(context)

@then('the number of rows in the database changes')
def step_impl(context):
    wait_for_postgresql_table_state_change(context)
