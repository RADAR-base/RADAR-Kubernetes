Feature: Fitbit registration, data ingestion and storage
e3)
  Scenario: Subject sends aRMT questionnaire response data to RADAR-base
    Given creation of an organization named "TEST"
    And creation of a project named "test"
    And creation of a subject named "test_user"
    And the current object counts in the s3 storage for files
      | bucket                      | filename_pattern        |
      | radar-intermediate-storage  | fitbit                  |
      | radar-output-storage        | fitbit                  |
    And registration of the subject with Fitbit authorization service
    Then Fitbit connector will pick start downloading the data and the object counts in the s3 storage for files have increased
      | bucket                      | filename_pattern        |
      | radar-intermediate-storage  | fitbit                  |
      | radar-output-storage        | fitbit                  |