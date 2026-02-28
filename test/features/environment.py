from playwright.sync_api import sync_playwright
from pages.management_portal_page import ManagementPortalPage
from pages.login_page import LoginPage

def before_all(context):
    context.cache = {
        "management_portal_token": None,
        "armt_source_type_json": None,
        "organization_json": None,
        "project_json": None,
        "armt_project_source_json": None,
        "test_subject_id": None,
        "secrets": None,
        "armt_meta_token": None,
        "armt_refresh_token": None,
        "armt_access_token": None,
        "rest_auth_registration_json": None,
        "fitbit_user_json": None,
    }
    context.state = {
        "database": {},
        "storage": {},
    }
    context.playwright = None
    context.browser = None
    context.page = None

    # Set up playwright if ANY scenario is tagged with @playwright.
    needs_playwright = any(
        scenario for feature in context._runner.features
        for scenario in feature.scenarios
        if "playwright" in scenario.tags
    )
    if needs_playwright:
        # Playwright supports two variations of the API: synchronous and asynchronous.
        # Currently using the synchronous one for simplicity, which could be changed to `async_playwright` if needed
        context.playwright = sync_playwright().start()
        dev_mode = context.config.userdata.get("dev_mode", "").lower() == "true"
        context.browser = context.playwright.chromium.launch(headless=(not dev_mode))

def before_scenario(context, scenario):
    if "playwright" in scenario.effective_tags:
        dev_mode = context.config.userdata.get("dev_mode", "").lower() == "true"
        context.page = context.browser.new_page()
        context.login_page = LoginPage(context.page, dev_mode)
        context.management_portal_page = ManagementPortalPage(context.page, dev_mode)

def after_scenario(context, scenario):
    if context.page:
        context.page.close()

def after_all(context):
    if context.browser:
        context.browser.close()
    if context.playwright:
        context.playwright.stop()
