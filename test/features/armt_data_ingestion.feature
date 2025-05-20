Feature: aRMT data ingestion and storage

  Scenario: Subject sends aRMT questionnaire response data to RADAR-base
    Given creation of an aRMT source type named "RADAR_aRMT"
    And creation of an organization named "TEST"
    And creation of a project named "test"
    And creation of an aRMT project source named "aRMT-test-source-TEST"
    And creation of a subject named "test_user"
    And the state of objects in the s3 storage
      | bucket                      | filename_pattern        | change_type |
      | radar-intermediate-storage  | questionnaire_response  | count       |
      | radar-output-storage        | questionnaire_response  | count       |
    And the aRMT application has retrieved an access token
    When the aRMT application sends questionnaire_response data
    """
    [
      {"questionId" : "1", "value": "Some Value", "startTime": 0, "endTime": 0}
    ]
    """
    Then the state of objects in the s3 storage changes
      | bucket                      | filename_pattern        | change_type |
      | radar-intermediate-storage  | questionnaire_response  | count       |
      | radar-output-storage        | questionnaire_response  | count       |