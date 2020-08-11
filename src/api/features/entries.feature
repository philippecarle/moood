Feature: Create entries
  In order to have my texts analyzed
  As an app user,
  I can post entries

  Scenario: Post an entry
    When I post an entry "Today I feel good. Just enjoying being myself, doing my best, writing an app. The weather is great and the weekend is finally here."
    Then the response header "content-type" should be "application/json"
    And the response status should be "Created"
    And the response JSON node "hash" should exist
    And a message should have been sent in the bus
    And this message MIME content type should be "application/json"
    And this message JSON node "hash" should exist
    And this message JSON node "createdAt" value should be now
    And this message JSON node "content" should be "Today I feel good. Just enjoying being myself, doing my best, writing an app. The weather is great and the weekend is finally here."
