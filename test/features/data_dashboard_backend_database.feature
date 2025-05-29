Feature: Data-dashboard-backend data transformation and ingestion

  Scenario: Questionnaire response data is stored in data dashboard backend database
    Given creation of an organization named "TEST"
    And creation of a project named "test"
    And creation of a subject named "test_user"
    And creation of an aRMT project source named "aRMT-test-source-TEST"
    And the aRMT application has retrieved an access token
    And these service states
      | service_name                                | state   |
      | ksql-server                                 | Running |
      | radar-jdbc-connector-data-dashboard-backend | Running |
      | data-dashboard-timescaledb-postgresql       | Running |
    And the number of rows in the database
      | service                                  | database        | table        |
      | data-dashboard-timescaledb-postgresql-0  | data-dashboard  | observation  |
    When the aRMT application sends questionnaire_response data
    """
    [
      {"questionId" : "1", "value": "Some Value", "startTime": 0, "endTime": 0}
    ]
    """
    Then the number of rows in the database changes
      | service                                  | database        | table        |
      | data-dashboard-timescaledb-postgresql-0  | data-dashboard  | observation  |