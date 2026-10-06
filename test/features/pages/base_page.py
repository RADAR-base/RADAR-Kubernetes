from playwright.sync_api import Page

class BasePage:
    def __init__(self, page: Page, dev_mode: bool):
        self.page = page
        self.dev_mode = dev_mode

    def go_to(self, url):
        self.page.goto(url)

    def take_screenshot(self, step_name):
        if self.dev_mode:
            screenshot_path = f"reports/screenshots/{step_name}.png"
            self.page.screenshot(path=screenshot_path)
