from playwright.sync_api import Page
from pages.base_page import BasePage

class LoginPage(BasePage):
    def __init__(self, page: Page, dev_mode: bool):
        super().__init__(page, dev_mode)
        self.username_input = page.locator("#identifier")
        self.password_input = page.locator("#password")
        self.email_input = page.get_by_role("textbox", name="Email")
        self.login_button = page.locator("button", has_text="Login")
        self.change_password_button = page.get_by_role("link", name="Forgot password?")
        self.submit_button = page.locator("button", has_text="Submit")
        self.signin_button = page.get_by_role("button", name="Sign in")
        self.header = page.get_by_role("heading", name="Sign In")
        self.password_change_confirmation = page.get_by_text("Recovery email has been sent!")

    def login(self, username, password):
        self.username_input.fill(username)
        self.password_input.fill(password)
        self.login_button.click()

    def get_header_locator(self):
        return self.header

    def click_change_password(self):
        self.change_password_button.click()

    def change_password(self):
        self.email_input.fill("test@test.com")
        self.submit_button.click()

    def get_confirm_password_change_locator(self):
        return self.password_change_confirmation

    def click_sign_in(self):
        self.signin_button.click()