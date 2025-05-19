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
      | grafana-metrics-timescaledb-postgresql      | Running |
    And registration of the subject with Fitbit authorization service
    Then Fitbit records are present in the database
      # 'count' matches the number of unique records returned by mockserver
      | service                                  | database         | table                         | count |
      | grafana-metrics-timescaledb-postgresql-0 | grafana-metrics  | connect_fitbit_intraday_steps | 4     |