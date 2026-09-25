# These scenarios need the Ory Kratos/Hydra/self-enrolment login flow, which is not deployed yet.
# They are excluded by default (see default_tags in behave.ini). Run them against an Ory-enabled
# deployment with: behave --tags=@ory
@ory
Feature: User authentication in Management Portal

  @playwright
  Scenario Outline: User signs in to Management Portal
    Given I open the Management Portal homepage
    When I click sign in button
    Then I should be redirected to the login page
    When I enter <email> and <password>
    Then I should see <message>
    Examples: Credentials
      | email          | password         | message                              |
      | $ADMIN_EMAIL   | $ADMIN_PASSWORD  | You are logged in as user            |
      | invalid_user   | invalid_password | The provided credentials are invalid |

  @playwright
  Scenario: User changes password
    Given I open the Management Portal homepage
    When I click sign in button
    Then I should be redirected to the login page
    When I click change password button
    And I specify new password
    Then password change should be requested
