from playwright.sync_api import sync_playwright, expect
from pages.management_portal_page import ManagementPortalPage
from pages.login_page import LoginPage
from pages.mp_login_page import MpLoginPage

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
    # Created here, at the root of the context, so that a browser started in a scenario hook
    # outlives that scenario (behave removes attributes set during a scenario when it ends).
    context.browser_state = {"playwright": None, "browser": None}
    context.page = None

def get_browser(context):
    state = context.browser_state
    if state["browser"] is None:
        # Playwright supports two variations of the API: synchronous and asynchronous.
        # Currently using the synchronous one for simplicity, which could be changed to `async_playwright` if needed
        state["playwright"] = sync_playwright().start()
        dev_mode = context.config.userdata.get("dev_mode", "").lower() == "true"
        state["browser"] = state["playwright"].chromium.launch(headless=(not dev_mode))
        # Redirects through the login pages can take a while on a busy test cluster
        expect.set_options(timeout=int(context.config.userdata.get("timeout_s", 10)) * 1000)
    return state["browser"]

def before_scenario(context, scenario):
    if "playwright" in scenario.effective_tags:
        # Start the browser only when a @playwright scenario actually runs, so that test runs
        # without UI scenarios don't need a browser installed.
        dev_mode = context.config.userdata.get("dev_mode", "").lower() == "true"
        context.page = get_browser(context).new_page()
        context.login_page = LoginPage(context.page, dev_mode)
        context.management_portal_page = ManagementPortalPage(context.page, dev_mode)
        context.mp_login_page = MpLoginPage(context.page, dev_mode)

def after_scenario(context, scenario):
    if context.page:
        context.page.close()
        context.page = None

def after_all(context):
    state = context.browser_state
    if state["browser"]:
        state["browser"].close()
    if state["playwright"]:
        state["playwright"].stop()
