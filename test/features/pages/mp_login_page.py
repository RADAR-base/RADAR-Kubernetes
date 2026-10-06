from playwright.sync_api import Page
from pages.base_page import BasePage

class MpLoginPage(BasePage):
    """Login form of the Management Portal's internal authserver (templates/login.html in ManagementPortal)."""
    def __init__(self, page: Page, dev_mode: bool):
        super().__init__(page, dev_mode)
        self.username_input = page.locator("#username")
        self.password_input = page.locator("#password")
        self.signin_button = page.get_by_role("button", name="Sign in")
        self.header = page.get_by_role("heading", name="Sign In")
        self.login_error = page.get_by_text("Wrong username or password")

    def login(self, username, password):
        self.username_input.fill(username)
        self.password_input.fill(password)
        self.signin_button.click()

    def get_header_locator(self):
        return self.header

    def get_login_error_locator(self):
        return self.login_error
