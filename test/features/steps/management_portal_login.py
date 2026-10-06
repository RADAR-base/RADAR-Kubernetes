from behave import then, when
from playwright.sync_api import expect
from base import get_secret
import re

@then('I should be redirected to the Management Portal login page')
def step_impl(context):
    expect(context.page).to_have_url(re.compile(r".*/managementportal/login.*"))
    expect(context.mp_login_page.get_header_locator()).to_be_visible()

@when('I sign in to the Management Portal as {username} with password {password}')
def step_impl(context, username, password):
    if password == "$ADMIN_PASSWORD":
        password = get_secret('management_portal', 'managementportal', 'common_admin_password', context=context)
    context.mp_login_page.login(username, password)
    context.mp_login_page.take_screenshot("after_mp_login")

@then('I should be logged in to the Management Portal as {username}')
def step_impl(context, username):
    expect(context.page).to_have_url(re.compile(r".*/managementportal/?(#.*)?$"))
    expect(context.management_portal_page.get_header_locator()).to_be_visible()
    expect(context.page.get_by_text(f'You are logged in as user "{username}"')).to_be_visible()

@then('the Management Portal login should be rejected')
def step_impl(context):
    expect(context.mp_login_page.get_login_error_locator()).to_be_visible()
    expect(context.page).to_have_url(re.compile(r".*/managementportal/login.*"))
