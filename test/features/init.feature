Feature: Create project and subject in management portal

  Scenario: Obtain management portal token
    Given management portal is running
     Then the management portal token can be requested

  Scenario: Create aRMT source type
    Given retrieval of management portal token
    And the aRMT source type does not exist
    Then the aRMT source type can be created

  Scenario: Create organization
    Given retrieval of management portal token
    And the test organization does not exist
    Then the test organization should be created

  Scenario: Create project
    Given retrieval of management portal token
    And the test project does not exist
    Then the test project should be created

  Scenario: Create aRMT project source
    Given retrieval of management portal token
    And the aRMT project source does not exist
    Then the aRMT project source should be created

  Scenario: Create subject
    Given retrieval of management portal token
    And the test subject does not exist
    Then the test subject should be created