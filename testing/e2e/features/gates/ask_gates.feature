Feature: Ask Gate time

  As a User
  I should be able to ask time

  Scenario: Ask question "What time is it?"
    Given a valid Question
    And Question asked
    When Use sent question
    Then We get currect time

  Scenario: Cannot  get time
    Given ask wrong question
    When I ask question " time is it?"
    Then I receive a "Error wrong question" error
