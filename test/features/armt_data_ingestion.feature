Feature: aRMT data ingestion and storage

  Scenario: Subject sends aRMT questionnaire response data to RADAR-base
    Given creation of an aRMT source type named "RADAR_aRMT"
    And creation of an organization named "TEST"
    And creation of a project named "test"
    And creation of an aRMT project source named "aRMT-test-source-TEST"
    And creation of a subject named "test_user"
    And the current object counts in the s3 storage for files
      | bucket                      | filename_pattern        |
      | radar-intermediate-storage  | questionnaire_response  |
      | radar-output-storage        | questionnaire_response  |
    And the aRMT application has retrieved an access token
    When the aRMT application sends questionnaire_response data
    """
    [
      {"questionId" : "1", "value": "Some Value", "startTime": 0, "endTime": 0}
    ]
    """
    Then the object counts in the s3 storage for files have increased by 1
      | bucket                      | filename_pattern        |
      | radar-intermediate-storage  | questionnaire_response  |
      | radar-output-storage        | questionnaire_response  |