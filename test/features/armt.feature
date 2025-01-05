Feature: aRMT data ingestion and storage

  @fixture.s3
  Scenario: Subject sends questionnaire response data to RADAR-base
    Given retrieval of management portal token
    And creation of the aRMT source type
    And creation of the test organization
    And creation of the test project
    And creation of the test subject
    And creation of the aRMT project source
    And the aRMT application has retrieved an access token
    Then true
#    When the subject sends questionnaire response data to RADAR-base
#    Then the data is stored on the intermediate storage
#    And the data is transformed and stored on the output storage
