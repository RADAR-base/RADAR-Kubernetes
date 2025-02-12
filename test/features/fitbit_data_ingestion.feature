Feature: Fitbit registration, data ingestion and storage

  Scenario: Subject sends aRMT questionnaire response data to RADAR-base
    Given creation of an organization named "TEST"
    And creation of a project named "test"
    And creation of a subject named "test_user"
    And these service states
      | service_name            | state   |
      | fitbit-connector        | Running |
      | rest-sources-backend    | Running |
    And the state of objects in the s3 storage
      | bucket                      | filename_pattern        | change_type |
      | radar-intermediate-storage  | fitbit                  | count       |
      | radar-output-storage        | fitbit                  | timestamp   |
    And registration of the subject with Fitbit authorization service
    Then Fitbit connector will download data and the state of objects in the s3 storage changes
      | bucket                      | filename_pattern        | change_type |
      | radar-intermediate-storage  | fitbit                  | count       |
      | radar-output-storage        | fitbit                  | timestamp   |