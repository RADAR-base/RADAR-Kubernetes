from playwright.sync_api import Page
from pages.base_page import BasePage

class ManagementPortalPage(BasePage):
    def __init__(self, page: Page, dev_mode: bool):
        super().__init__(page, dev_mode)
        self.sing_in_button = page.locator("jhi-home").get_by_text("Sign in")
        self.header = page.get_by_role("heading", name="Management Portal")

    def click_sign_in(self):
        self.sing_in_button.click()

    def get_header_locator(self):
        return self.header
