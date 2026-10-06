# Signs in through the Management Portal's own login page (internal IDP and authserver).
# For the login through Ory Kratos/Hydra, see user_authentication.feature.
@playwright
Feature: Management Portal login with the internal authserver

  Scenario: Admin signs in to the Management Portal
    Given I open the Management Portal homepage
    When I click sign in button
    Then I should be redirected to the Management Portal login page
    When I sign in to the Management Portal as admin with password $ADMIN_PASSWORD
    Then I should be logged in to the Management Portal as admin

  Scenario: Sign in with wrong credentials is rejected
    Given I open the Management Portal homepage
    When I click sign in button
    Then I should be redirected to the Management Portal login page
    When I sign in to the Management Portal as invalid_user with password invalid_password
    Then the Management Portal login should be rejected
