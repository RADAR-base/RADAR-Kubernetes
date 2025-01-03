Feature: Setup project and subject in management portal

  Scenario: Obtain management portal token
    Given management portal is running
     Then the management portal token can be requested

  Scenario: Create aRMT source type
    Given retrieval of management portal token
    Then the aRMT source type can be created

  Scenario: Create organization
    Given retrieval of management portal token
    Then the test organization should be created

  Scenario: Create project
    Given retrieval of management portal token
    Then the test project should be created

  Scenario: Create aRMT project source
    Given retrieval of management portal token
    Then the aRMT project source should be created

  Scenario: Create subject
    Given retrieval of management portal token
    Then the test subject should be created