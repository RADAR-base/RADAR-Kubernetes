Feature: Grafana metrics dashboard data ingestion
e3)
  Scenario: Fitbit data is stored in the grafana metrics database
    Given creation of an organization named "TEST"
    And creation of a project named "test"
    And creation of a subject named "test_user"
    And creation of an aRMT project source named "aRMT-test-source-TEST"
    And these service states
      | service_name                                | state   |
      | ksql-server                                 | Running |
      | radar-jdbc-connector-grafana                | Running |
      | grafana-timescaledb                         | Running |
    And registration of the subject with Fitbit authorization service
    Then Fitbit records are present in the database
      # 'count' matches the number of unique records returned by mockserver
      | service                     | database     | table                         | count |
      | grafana-timescaledb-1       | grafana      | connect_fitbit_intraday_steps | 4     |
