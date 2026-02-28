from behave import given, when, then, fixture, use_fixture
from playwright.sync_api import sync_playwright, expect
from base import get_credential_value, get_secret
import re

@given('I open the Management Portal homepage')
def step_impl(context):
    base_url = context.config.userdata.get("url")
    if not base_url:
        raise ValueError("'url' not set in userdata")
    context.management_portal_page.go_to(f'{base_url}/managementportal')

@when('I click sign in button')
def step_impl(context):
    context.management_portal_page.click_sign_in()

@then('I should be redirected to the login page')
def step_impl(context):
    pattern = re.compile(r".*/study/auth/oauth-login.*login_challenge.*")
    expect(context.page).to_have_url(pattern)


@when('I enter {email} and {password}')
def step_impl(context, email, password):
    try:
        final_email = get_credential_value(context, email)
        final_password = get_secret('management_portal', 'managementportal', 'common_admin_password', context=context) \
            if password.startswith("$") \
            else password
    except ValueError as e:
        raise
    context.login_page.login(final_email, final_password)
    context.login_page.take_screenshot("after_login")

@then('I should see {message}')
def step_then_see_message(context, message):
    if "invalid" in message.lower():
        expect(context.login_page.get_header_locator()).to_be_visible()
        target_locator = context.login_page.page.get_by_text(message)
        context.management_portal_page.take_screenshot("invalid_login_message")
    else:
        context.login_page.click_sign_in()
        expect(context.page).to_have_url(re.compile(r".*/managementportal.*"))
        expect(context.management_portal_page.get_header_locator()).to_be_visible()
        target_locator = context.management_portal_page.page.get_by_text(message)
        context.management_portal_page.take_screenshot("successful_login_message")
    expect(target_locator).to_be_visible()

@when('I click change password button')
def step_impl(context):
    context.login_page.click_change_password()

@when('I specify new password')
def step_impl(context):
    context.login_page.change_password()

@then('password change should be requested')
async def step_impl(context):
    response = await context.login_page.wait_for_response(
        lambda r: (
                "study/api/ory/recovery" in r.url
                and r.request.method == "POST"
        )
    )
    confirmation_locator = context.login_page.get_confirm_password_change_locator()
    expect(confirmation_locator).to_be_visible()
    assert response.ok, f"Expected 200 OK but got {response.status} for {response.url}"
    context.login_page.take_screenshot("after_password_change")
