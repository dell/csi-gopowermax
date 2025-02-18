Feature: PMAX Client library
  As a developer of the PMAX CSI driver that uses the PMAX REST client library
  I want to test the client library methods
  So they are known to work and have a high test coverage.

  Scenario Outline: Authenticate test cases
    When I induce error <induced> 
    And I call authenticate with endpoint <endpoint> credentials <credentials> apiversion <apiversion>
    Then the error message contains <errormsg>

    Examples:
    | endpoint    | credentials    | apiversion     |induced          | errormsg                    |
    | "mockurl"   | "good"         |     "90"       | "none"          | "none"                      |
    | "mockurl"   | "bad"          |     "90"       | "none"          | "Unauthorized"              |
    | "badurl"    | "good"         |     "90"       | "none"          | "connect"                   | 
    | "nilurl"    | "good"         |     "90"       | "none"          | "Endpoint must be supplied" |
    | "mockurl"   | "good"         |     "90"       | "httpStatus500" | "Internal Error"            |
    | "mockurl"   | "good"         |     "91"       | "none"          | "none"                      |
    | "mockurl"   | "bad"          |     "91"       | "none"          | "Unauthorized"              |
    | "badurl"    | "good"         |     "91"       | "none"          | "connect"                   | 
    | "nilurl"    | "good"         |     "91"       | "none"          | "Endpoint must be supplied" |
    | "mockurl"   | "good"         |     "91"       | "httpStatus500" | "Internal Error"            |
  
  Scenario Outline: TestCases for GetSymmetrixIDList
    Given a valid connection
    And I induce error <induced>
    When I call GetSymmetrixIDList
    Then the error message contains <errormsg>
    And I get a valid Symmetrix ID List if no error

    Examples:
    | induced               | errormsg                      |
    | "none"                | "none"                        |
    | "GetSymmetrixError"   | "induced error"               |

  Scenario Outline: Get Symmetrix System
    Given a valid connection
    And I induce error <induced>
    When I call GetSymmetrixByID <id>
    Then the error message contains <errormsg>
    And I get a valid Symmetrix Object if no error

    Examples:
    | id              | induced               | errormsg                    |
    | "000197900046"  | "none"                | "none"                      |
    | "000000000000"  | "none"                | "not found"                 |
    | "000197900046"  | "GetSymmetrixError"   | "induced error"             |
    | "000197900046"  | "httpStatus500"       | "Internal Error"            |
    | "000197900046"  | "InvalidJSON"         | "invalid character"         |

  Scenario Outline: Test cases for GetVolumeIDList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have <nvols> volumes
    And I induce error <induced>
    When I call GetVolumeIDList <volume_identifier> 
    Then the error message contains <errormsg>
    And I get a valid VolumeIDList with <vols> if no error
    
    Examples:                # volumes are numbered 1...n  Vol00001, Vol00002, ...
    | nvols      | vols  | volume_identifier | induced                    | errormsg                      | arrays    |
    | 7          | 7     | ""                | "none"                     | "none"                        | ""        |
    | 10         | 10    | ""                | "none"                     | "none"                        | ""        |
    | 11         | 11    | ""                | "none"                     | "none"                        | ""        |
    | 23         | 23    | ""                | "none"                     | "none"                        | ""        |
    | 23         | 23    | ""                | "GetVolumeIteratorError"   | "induced error"               | ""        |
    | 23         | 23    | ""                | "httpStatus500"            | "Internal Error"              | ""        |
    | 23         | 23    | ""                | "InvalidJSON"              | "invalid character"           | ""        |
    | 23         | 1     | "Vol00005"        | "none"                     | "none"                        | ""        |
    | 23         | 1     | "Vol00015"        | "none"                     | "none"                        | ""        |
    | 23         | 1     | "Vol00015"        | "none"                     | "none"                        | ""        |
    | 23         | 0     | "ABCDEFGH"        | "none"                     | "none"                        | ""        |
    | 23         | 9     | "<like>Vol0000"   | "none"                     | "none"                        | ""        |
    | 23         | 10    | "<like>Vol0001"   | "none"                     | "none"                        | ""        |
    | 23         | 4     | "<like>Vol0002"   | "none"                     | "none"                        | ""        |
    | 5          | 5     | ""                | "none"                     | "ignored as it is not managed"| "ignore"  |

  Scenario Outline: Test cases for GetVolumeByID
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have 5 volumes
    And I induce error <induced>
    When I call GetVolumeByID <id>
    Then the error message contains <errormsg>
    And I get a valid Volume Object <id> if no error

    Examples:
    | id              | induced               | errormsg                      | arrays    |
    | "00001"         | "none"                | "none"                        | ""        |
    | "00003"         | "none"                | "none"                        | ""        |
    | "00010"         | "none"                | "Could not find"              | ""        |
    | "00001"         | "GetVolumeError"      | "induced error"               | ""        |
    | "00001"         | "httpStatus500"       | "Internal Error"              | ""        |
    | "00001"         | "InvalidJSON"         | "invalid character"           | ""        |
    | "00001"         | "none"                | "ignored as it is not managed"| "ignored" |

  Scenario Outline: Test cases for volume expand
    Given a valid connection
    And I have 2 volumes
    And I induce error <induced>
    Then I expand volume <id> to <size> in GB
    And the error message contains <errormsg>
    And I validate that volume <id> has has size <size> in GB

    Examples:
      | id      | size | induced             | errormsg        |
      | "00001" | "10" | "none"              | "none"          |
      | "00002" | "10" | "GetVolumeError"    | "induced error" |
      | "00001" | "10" | "ExpandVolumeError" | "induced error" |

  Scenario Outline: Test cases for volume expand with units
    Given a valid connection
    And I have 2 volumes
    And I induce error <induced>
    Then I expand volume <id> to <size> in <units>
    And the error message contains <errormsg>
    And I validate that volume <id> has has size <size> in GB

    Examples:
      | id      | size | units |   induced             | errormsg        |
      | "00001" | "10" | "GB"  |   "none"              | "none"          |
      | "00002" | "10" | "GB"  |   "GetVolumeError"    | "induced error" |
      | "00001" | "10" | "GB"  |   "ExpandVolumeError" | "induced error" |

  Scenario Outline: Test cases for GetStorageGroupIDList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call GetStorageGroupIDList
    Then the error message contains <errormsg>
    And I get a valid StorageGroupIDList if no errors

    Examples:
    | induced               | errormsg                      | arrays    |
    | "none"                | "none"                        | ""        |
    | "GetStorageGroupError"| "induced error"               | ""        |
    | "httpStatus500"       | "Internal Error"              | ""        |
    | "InvalidJSON"         | "invalid character"           | ""        |
    | "none"                | "ignored as it is not managed"| "ignored" |
    | "InvalidResponse"     | "EOF"                         | ""        |

  Scenario Outline: Test cases for GetStorageGroup
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call GetStorageGroup <name>
    Then the error message contains <errormsg>
    And I get a valid StorageGroup if no errors

    Examples:
    | name               | induced               | errormsg                      | arrays    |
    | "CSI-Test-SG-1"    | "none"                | "none"                        | ""        |
    | "CSI-Test-SG-1"    | "GetStorageGroupError"| "induced error"               | ""        |
    | "CSI-Test-SG-1"    | "httpStatus500"       | "Internal Error"              | ""        |
    | "CSI-Test-SG-1"    | "InvalidJSON"         | "invalid character"           | ""        |
    | "CSI-Test-SG-1"    | "none"                | "ignored as it is not managed"| "ignored" |
    | "CSI-Test-SG-1"    | "InvalidResponse"     | "EOF"                         | ""        |

  Scenario Outline: Test cases for GetStoragePool
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call GetStoragePool <name>
    Then the error message contains <errormsg>
    And I get a valid GetStoragePool if no errors

    Examples:
    | name     | induced               | errormsg                      | arrays    |
    | "SRP_1"  | "none"                | "none"                        | ""        |
    | "SRP_1"  | "GetStoragePoolError" | "induced error"               | ""        |
    | "SRP_1"  | "httpStatus500"       | "Internal Error"              | ""        |
    | "SRP_1"  | "InvalidJSON"         | "invalid character"           | ""        |
    | "SRP_1"  | "none"                | "ignored as it is not managed"| "ignored" |

  Scenario Outline: Test cases for GetJobIDList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    And I have <njobs> jobs
    When I call GetJobIDList with <status>
    Then the error message contains <errormsg>
    And I get a valid JobsIDList with <njobs> if no errors

    Examples:
    | njobs         | induced                     | status       | errormsg                        | arrays    |
    | 1             | "none"                      | ""           | "none"                          | ""        |
    | 0             | "none"                      | ""           | "none"                          | ""        |
    | 20            | "none"                      | ""           | "none"                          | ""        |
    | 20            | "none"                      | "SCHEDULED"  | "none"                          | ""        |
    | 1             | "GetJobError"               | ""           | "induced error"                 | ""        |
    | 20            | "httpStatus500"             | ""           | "Internal Error"                | ""        |
    | 20            | "InvalidJSON"               | ""           | "invalid character"             | ""        |
    | 1             | "none"                      | ""           | "ignored as it is not managed"  | "ignored" |

  Scenario Outline: Test cases for GetJobByID
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    And I create a job with initial state <initial> and final state <final>
    When I call GetJobByID
    And the error message contains <errormsg>
    And I get a valid Job with state <initial> if no error
    Then I call GetJobByID
    And I get a valid Job with state <final> if no error

    Examples:
    | initial        | final            | induced                       | errormsg                       | arrays    |
    | "RUNNING"      | "SUCCEEDED"      | "none"                        | "none"                         | ""        |
    | "RUNNING"      | "FAILED"         | "none"                        | "none"                         | ""        |
    | "RUNNING"      | "RUNNING"        | "none"                        | "none"                         | ""        |
    | "RUNNING"      | "SUCCEEDED"      | "GetJobError"                 | "induced error"                | ""        |
    | "RUNNING"      | "SUCCEEDED"      | "httpStatus500"               | "Internal Error"               | ""        |
    | "RUNNING"      | "SUCCEEDED"      | "InvalidJSON"                 | "invalid character"            | ""        |
    | "RUNNING"      | "SUCCEEDED"      | "GetJobCannotFindRoleForUser" | "none"                         | ""        |
    | "RUNNING"      | "SUCCEEDED"      | "none"                        | "ignored as it is not managed" | "ignored" |

  Scenario Outline: Test cases WaitOnJobCompletion
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    And I create a job with initial state <initial> and final state <final>
    When I call WaitOnJobCompletion
    Then the error message contains <errormsg>
    And I get a valid Job with state <final> if no error

    Examples:
    | initial        | final            | induced          | errormsg                       | arrays    |
    | "RUNNING"      | "SUCCEEDED"      | "none"           | "none"                         | ""        |
    | "RUNNING"      | "FAILED"         | "none"           | "none"                         | ""        |
    | "RUNNING"      | "RUNNING"        | "none"           | "timed out after"              | ""        |
    | "RUNNING"      | "SUCCEEDED"      | "GetJobError"    | "induced error"                | ""        |
    | "RUNNING"      | "SUCCEEDED"      | "none"           | "ignored as it is not managed" | "ignored" |

  Scenario Outline: Test cases for CreateVolumeInStorageGroup for v90
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateVolumeInStorageGroup with name <volname> and size <size>
    Then the error message contains <errormsg>
    And I get a valid Volume with name <volname> if no error

    Examples:
    | volname                                                                        | size     | induced                   | errormsg                                               | arrays    |
    | "IntgA"                                                                        | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgB"                                                                        | 5        | "none"                    | "none"                                                 | ""        |
    | "IntgC"                                                                        | 1        | "UpdateStorageGroupError" | "A job was not returned from UpdateStorageGroup"       | ""        |
    | "IntgD"                                                                        | 1        | "httpStatus500"           | "A job was not returned from UpdateStorageGroup"       | ""        |
    | "IntgE"                                                                        | 1        | "GetJobError"             | "induced error"                                        | ""        |
    | "IntgF"                                                                        | 1        | "JobFailedError"          | "the UpdateStorageGroup job failed"                    | ""        |
    | "IntgG"                                                                        | 1        | "GetVolumeError"          | "failed to find newly created volume with name: IntgG" | ""        |
    | "IntgH"                                                                        | 1        | "VolumeNotCreatedError"   | "failed to find newly created volume with name: IntgH" | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy"| 1        | "none"                    | "Length of volumeName exceeds max limit"               | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk"              | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgA"                                                                        | 1        | "none"                    | "ignored as it is not managed"                         | "ignored" |


Scenario Outline: Test cases for modifyMobility for volume
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have 2 volumes
    And I induce error <induced>
    Then I call ModifyMobility for Volume with id <id> to mobility <mobilityenabled>
    And the error message contains <errormsg>
    And I validate that volume has mobility modified to <mobilityenabled>
    

    Examples:
    | id                                                                             | induced                   | errormsg                               | mobilityenabled  | arrays    |
    | "00001"                                                                        | "none"                    | "none"                                 | "true"           | ""        |
    | "00002"                                                                        | "ModifyMobilityError"     | "Error modifying mobility for volume:" | "false"          | ""        |
    | "00002"                                                                        | "none"                    | "ignored as it is not managed"         | "false"          | "ignored" |

  Scenario Outline: Test cases for GetStorageGroupSnapshotPolicy
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call GetStorageGroupSnapshotPolicy with <symID> <snapshotPolicyID> <storageGroupID>
    Then the error message contains <errormsg>
    And I get a valid StorageGroupSnapshotPolicy Object if no error

    Examples:
    | symID              | snapshotPolicyID     | storageGroupID   | induced                                        | errormsg                                                                | arrays    |
    | "000000000001"     | "IntSPA"             | "IntSGA"         | "none"                                         | "none"                                                                  | ""        |
    | "000000000002"     | "IntSPB"             | "IntSGB"         | "GetStorageGroupSnapshotPolicyError"           | "Error retrieving storage group snapshot policy: induced error"    | ""        |


  Scenario Outline: Test cases for CreateVolumeInStorageGroup for v90 with capacity unit
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateVolumeInStorageGroup with name <volname> and size <size> and unit <capUnit>
    Then the error message contains <errormsg>
    And I get a valid Volume with name <volname> if no error

    Examples:
    | volname                                                                        | size     | capUnit   | induced                   | errormsg                                               | arrays    |
    | "IntgA"                                                                        | 1        | "CYL"     | "none"                    | "none"                                                 | ""        |
    | "IntgB"                                                                        | 5        | "CYL"     | "none"                    | "none"                                                 | ""        |
    | "IntgC"                                                                        | 1        | "GB"      | "UpdateStorageGroupError" | "A job was not returned from UpdateStorageGroup"       | ""        |
    | "IntgD"                                                                        | 1        | "GB"      | "httpStatus500"           | "A job was not returned from UpdateStorageGroup"       | ""        |
    | "IntgE"                                                                        | 1        | "GB"      | "GetJobError"             | "induced error"                                        | ""        |
    | "IntgF"                                                                        | 1        | "GB"      | "JobFailedError"          | "the UpdateStorageGroup job failed"                    | ""        |
    | "IntgG"                                                                        | 1        | "GB"      | "GetVolumeError"          | "failed to find newly created volume with name: IntgG" | ""        |
    | "IntgH"                                                                        | 1        | "GB"      | "VolumeNotCreatedError"   | "failed to find newly created volume with name: IntgH" | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy"| 1        | "GB"      | "none"                    | "Length of volumeName exceeds max limit"               | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk"              | 1        | "CYL"     | "none"                    | "none"                                                 | ""        |
    | "IntgA"                                                                        | 1        | "GB"      | "none"                    | "ignored as it is not managed"                         | "ignored" |
    | "IntgI"                                                                        | 2        | "GB"      | "none"                    | "none"                                                 | ""        |           


Scenario Outline: Test cases for Synchronous CreateVolumeInStorageGroup for v90
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateVolumeInStorageGroupS with name <volname> and size <size>
    Then the error message contains <errormsg>
    And I get a valid Volume with name <volname> if no error

    Examples:
    | volname                                                                        | size     | induced                   | errormsg                                               | arrays    |
    | "IntgA"                                                                        | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgB"                                                                        | 5        | "none"                    | "none"                                                 | ""        |
    | "IntgG"                                                                        | 1        | "GetVolumeError"          | "failed to find newly created volume with name: IntgG" | ""        |
    | "IntgH"                                                                        | 1        | "VolumeNotCreatedError"   | "failed to find newly created volume with name: IntgH" | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy"| 1        | "none"                    | "Length of volumeName exceeds max limit"               | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk"              | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgA"                                                                        | 1        | "none"                    | "ignored as it is not managed"                         | "ignored" |

Scenario Outline: Test cases for Synchronous CreateVolumeInStorageGroup for v90 with capacity unit
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateVolumeInStorageGroupS with name <volname> and size <size> and unit <capUnit>
    Then the error message contains <errormsg>
    And I get a valid Volume with name <volname> if no error

    Examples:
    | volname                                                                        | size     |capUnit  | induced                   | errormsg                                               | arrays    |
    | "IntgA"                                                                        | 1        | "CYL"   | "none"                    | "none"                                                 | ""        |
    | "IntgB"                                                                        | 5        | "CYL"   | "none"                    | "none"                                                 | ""        |
    | "IntgG"                                                                        | 1        | "CYL"   | "GetVolumeError"          | "failed to find newly created volume with name: IntgG" | ""        |
    | "IntgH"                                                                        | 1        | "CYL"   | "VolumeNotCreatedError"   | "failed to find newly created volume with name: IntgH" | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy"| 1        | "CYL"   | "none"                    | "Length of volumeName exceeds max limit"               | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk"              | 1        | "CYL"   | "none"                    | "none"                                                 | ""        |
    | "IntgA"                                                                        | 1        | "CYL"   | "none"                    | "ignored as it is not managed"                         | "ignored" |

Scenario Outline: Test cases for Synchronous CreateVolumeInStorageGroup with metadata headers for v90
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateVolumeInStorageGroupSWithMetaDataHeaders with name <volname> and size <size>
    Then the error message contains <errormsg>
    And I get a valid Volume with name <volname> if no error

    Examples:
    | volname                                                                        | size     | induced                   | errormsg                                               | arrays    |
    | "IntgA"                                                                        | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgB"                                                                        | 5        | "none"                    | "none"                                                 | ""        |
    | "IntgG"                                                                        | 1        | "GetVolumeError"          | "failed to find newly created volume with name: IntgG" | ""        |
    | "IntgH"                                                                        | 1        | "VolumeNotCreatedError"   | "failed to find newly created volume with name: IntgH" | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy"| 1        | "none"                    | "Length of volumeName exceeds max limit"               | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk"              | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgA"                                                                        | 1        | "none"                    | "ignored as it is not managed"                         | "ignored" |
  
  Scenario Outline: Test cases for CreateVolumeInStorageGroup for v91
    Given a valid v91 connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateVolumeInStorageGroup with name <volname> and size <size>
    Then the error message contains <errormsg>
    And I get a valid Volume with name <volname> if no error

    Examples:
    | volname                                                                        | size     | induced                   | errormsg                                               | arrays    |
    | "IntgA"                                                                        | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgB"                                                                        | 5        | "none"                    | "none"                                                 | ""        |
    | "IntgC"                                                                        | 1        | "UpdateStorageGroupError" | "A job was not returned from UpdateStorageGroup"       | ""        |
    | "IntgD"                                                                        | 1        | "httpStatus500"           | "A job was not returned from UpdateStorageGroup"       | ""        |
    | "IntgE"                                                                        | 1        | "GetJobError"             | "induced error"                                        | ""        |
    | "IntgF"                                                                        | 1        | "JobFailedError"          | "the UpdateStorageGroup job failed"                    | ""        |
    | "IntgG"                                                                        | 1        | "GetVolumeError"          | "failed to find newly created volume with name: IntgG" | ""        |
    | "IntgH"                                                                        | 1        | "VolumeNotCreatedError"   | "failed to find newly created volume with name: IntgH" | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy"| 1        | "none"                    | "Length of volumeName exceeds max limit"               | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk"              | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgA"                                                                        | 1        | "none"                    | "ignored as it is not managed"                         | "ignored" |

  Scenario Outline: Test cases for Synchronous CreateVolumeInStorageGroup for v91
    Given a valid v91 connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateVolumeInStorageGroup with name <volname> and size <size>
    Then the error message contains <errormsg>
    And I get a valid Volume with name <volname> if no error

    Examples:
    | volname                                                                        | size     | induced                   | errormsg                                               | arrays    |
    | "IntgA"                                                                        | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgB"                                                                        | 5        | "none"                    | "none"                                                 | ""        |
    | "IntgG"                                                                        | 1        | "GetVolumeError"          | "failed to find newly created volume with name: IntgG" | ""        |
    | "IntgH"                                                                        | 1        | "VolumeNotCreatedError"   | "failed to find newly created volume with name: IntgH" | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy"| 1        | "none"                    | "Length of volumeName exceeds max limit"               | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk"              | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgA"                                                                        | 1        | "none"                    | "ignored as it is not managed"                         | "ignored" |

  Scenario Outline: Test cases for Synchronous CreateVolumeInStorageGroup with metadata headers for v91
    Given a valid v91 connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateVolumeInStorageGroupSWithMetaDataHeaders with name <volname> and size <size>
    Then the error message contains <errormsg>
    And I get a valid Volume with name <volname> if no error

    Examples:
    | volname                                                                        | size     | induced                   | errormsg                                               | arrays    |
    | "IntgA"                                                                        | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgB"                                                                        | 5        | "none"                    | "none"                                                 | ""        |
    | "IntgG"                                                                        | 1        | "GetVolumeError"          | "failed to find newly created volume with name: IntgG" | ""        |
    | "IntgH"                                                                        | 1        | "VolumeNotCreatedError"   | "failed to find newly created volume with name: IntgH" | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy"| 1        | "none"                    | "Length of volumeName exceeds max limit"               | ""        |
    | "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk"              | 1        | "none"                    | "none"                                                 | ""        |
    | "IntgA"                                                                        | 1        | "none"                    | "ignored as it is not managed"                         | "ignored" |

  Scenario Outline: Test cases for Remove Volume From Storage Group
    Given a valid connection
    And I call CreateVolumeInStorageGroup with name "IntM" and size 1
    And I induce error <induced>
    And I have an allowed list of <arrays>
    When I call RemoveVolumeFromStorageGroup
    Then the error message contains <errormsg>
    And the volume is no longer a member of the Storage Group if no error

    Examples:
    | induced                   | errormsg                                         | arrays    |
    | "none"                    | "none"                                           | ""        |
    | "UpdateStorageGroupError" | "induced error"                                  | ""        |
    | "none"                    | "ignored as it is not managed"                   | "ignored" |

    Scenario Outline: Test cases for Remove Volume From Storage Group for v91
    Given a valid v91 connection
    And I call CreateVolumeInStorageGroup with name "IntM" and size 1
    And I induce error <induced>
    And I have an allowed list of <arrays>
    When I call RemoveVolumeFromStorageGroup
    Then the error message contains <errormsg>
    And the volume is no longer a member of the Storage Group if no error

    Examples:
    | induced                   | errormsg                                         | arrays    |
    | "none"                    | "none"                                           | ""        |
    | "UpdateStorageGroupError" | "induced error"                                  | ""        |
    | "none"                    | "ignored as it is not managed"                   | "ignored" |

  Scenario Outline: Test cases for Rename Volume
    Given a valid connection
    And I call CreateVolumeInStorageGroup with name "IntN" and size 1
    And I induce error <induced>
    And I have an allowed list of <arrays>
    When I call RenameVolume with <newname>
    Then the error message contains <errormsg>
    And I get a valid Volume with name <newname> if no error

    Examples:
    | newname              | induced                   | errormsg                                         | arrays    |
    | "Renamed"            | "none"                    | "none"                                           | ""        |               
    | "Renamed"            | "UpdateVolumeError"       | "induced error"                                  | ""        |
    | "Renamed"            | "none"                    | "ignored as it is not managed"                   | "ignored" |

    Scenario Outline: Test cases for Initiate Deallocation of Tracks
    Given a valid connection
    And I call CreateVolumeInStorageGroup with name "IntO" and size 1
    And I induce error <induced>
    And I have an allowed list of <arrays>
    When I call InitiateDeallocationOfTracksFromVolume
    Then the error message contains <errormsg>
    And I get a valid Job with state "RUNNING" if no error

    Examples:
    | induced                   | errormsg                                         | arrays    |
    | "none"                    | "none"                                           | ""        |
    | "UpdateVolumeError"       | "induced error"                                  | ""        |
    | "none"                    | "ignored as it is not managed"                   | "ignored" |

  Scenario Outline: Test cases for Delete Volume with errors
    Given a valid connection
    And I call CreateVolumeInStorageGroup with name "IntP" and size 1
    And I induce error <induced>
    And I have an allowed list of <arrays>
    When I call DeleteVolume
    Then the error message contains <errormsg>

    Examples:
    | induced                   | errormsg                                         | arrays    |
    | "none"                    | "present in storage group"                       | ""        |
    | "DeleteVolumeError"       | "induced error"                                  | ""        |
    | "none"                    | "ignored as it is not managed"                   | "ignored" |

  Scenario: Test cases for Delete Volume
    Given a valid connection
    And I call CreateVolumeInStorageGroup with name "DEL" and size 1
    When I call RemoveVolumeFromStorageGroup
    Then the error message contains "none"
    And the volume is no longer a member of the Storage Group if no error
    When I call DeleteVolume
    Then the error message contains "none"

  Scenario Outline: Test cases for CreateStorageGroup for v90
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateStorageGroup with name <sgname> and srp <srp> and sl <sl>
    Then the error message contains <errormsg>
    And I get a valid StorageGroup with name <sgname> if no error

    Examples:
    | sgname               | srp      | sl           | induced                    | errormsg                                              | arrays    |
    | "CSI-Test-New-SG1"   | "SRP_1"  | "Diamond"    | "none"                     | "none"                                                | ""        |
    | "CSI-Test-New-SG1"   | "None"   | "Diamond"    | "none"                     | "none"                                                | ""        |
    | "CSI-Test-New-SG2"   | "SRP_1"  | "Optimized"  | "none"                     | "none"                                                | ""        |
    | "CSI-Test-New-SG2"   | "SRP_1"  | "Optimized"  | "StorageGroupAlreadyExists"| "The requested storage group resource already exists" | ""        |
    | "CSI-Test-New-SG3"   | "SRP_1"  | "Diamond"    | "CreateStorageGroupError"  | "induced error"                                       | ""        |
    | "CSI-Test-New-SG4"   | "SRP_1"  | "Diamond"    | "httpStatus500"            | "Internal Error"                                      | ""        |
    | "CSI-Test-New-SG1"   | "SRP_1"  | "Diamond"    | "none"                     | "ignored as it is not managed"                        | "ignored" |
    | "CSI-Test-New-SG1"   | "SRP_1"  | "Diamond"    | "InvalidResponse"          | "EOF"                                                 | ""        |

  Scenario Outline: Test cases for CreateStorageGroup for v91
    Given a valid v91 connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateStorageGroup with name <sgname> and srp <srp> and sl <sl>
    Then the error message contains <errormsg>
    And I get a valid StorageGroup with name <sgname> if no error

    Examples:
    | sgname               | srp      | sl           | induced                    | errormsg                                              | arrays    |
    | "CSI-Test-New-SG1"   | "SRP_1"  | "Diamond"    | "none"                     | "none"                                                | ""        |
    | "CSI-Test-New-SG1"   | "None"   | "Diamond"    | "none"                     | "none"                                                | ""        |
    | "CSI-Test-New-SG2"   | "SRP_1"  | "Optimized"  | "none"                     | "none"                                                | ""        |
    | "CSI-Test-New-SG2"   | "SRP_1"  | "Optimized"  | "StorageGroupAlreadyExists"| "The requested storage group resource already exists" | ""        |
    | "CSI-Test-New-SG3"   | "SRP_1"  | "Diamond"    | "CreateStorageGroupError"  | "induced error"                                       | ""        |
    | "CSI-Test-New-SG4"   | "SRP_1"  | "Diamond"    | "httpStatus500"            | "Internal Error"                                      | ""        |
    | "CSI-Test-New-SG1"   | "SRP_1"  | "Diamond"    | "none"                     | "ignored as it is not managed"                        | "ignored" |
    | "CSI-Test-New-SG1"   | "SRP_1"  | "Diamond"    | "InvalidResponse"          | "EOF"                                                 | ""        |

  Scenario Outline: Test cases for CreateStorageGroupWithHostLimits for v91
    Given a valid v91 connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateStorageGroup with name <sgname> and srp <srp> and sl <sl> and hostlimits <hl>
    Then the error message contains <errormsg>
    And I get a valid StorageGroup with name <sgname> if no error

    Examples:
    | sgname               | srp      | sl           | hl             |   induced                    | errormsg                                              | arrays    |
    | "CSI-Test-New-SG1"   | "SRP_1"  | "Diamond"    | "1:100:Never"  |   "none"                     | "none"                                                | ""        |
    | "CSI-Test-New-SG1"   | "None"   | "Diamond"    | "1:100:Never"  |   "none"                     | "none"                                                | ""        |
    | "CSI-Test-New-SG2"   | "SRP_1"  | "Optimized"  | "1:100:Never"  |   "none"                     | "none"                                                | ""        |
    | "CSI-Test-New-SG2"   | "SRP_1"  | "Optimized"  | "1:100:Never"  |   "StorageGroupAlreadyExists"| "The requested storage group resource already exists" | ""        |
    | "CSI-Test-New-SG3"   | "SRP_1"  | "Diamond"    | "1:100:Never"  |   "CreateStorageGroupError"  | "induced error"                                       | ""        |
    | "CSI-Test-New-SG4"   | "SRP_1"  | "Diamond"    | "1:100:Never"  |   "httpStatus500"            | "Internal Error"                                      | ""        |
    | "CSI-Test-New-SG1"   | "SRP_1"  | "Diamond"    | "1:100:Never"  |   "none"                     | "ignored as it is not managed"                        | "ignored" |
    | "CSI-Test-New-SG1"   | "SRP_1"  | "Diamond"    | "1:100:Never"  |   "InvalidResponse"          | "EOF"                                                 | ""        |

  Scenario Outline: Test DeleteStorageGroup
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call DeleteStorageGroup <name>
    Then the error message contains <errormsg>

    Examples:
    | induced                        | name                  | errormsg                           | arrays    |
    | "none"                         | "CSI-Test-SG-2"       | "none"                             | ""        |
    | "DeleteStorageGroupError"      | "CSI-Test-SG-3"       | "induced error"                    | ""        |
    | "none"                         | "CSI-Test-SG-3"       |"ignored as it is not managed"      | "ignored" |

  Scenario Outline: Test GetStoragePoolList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call GetStoragePoolList
    Then the error message contains <errormsg>
    And I get a valid StoragePoolList if no error

    Examples:
    | induced                        | errormsg                                              | arrays    |
    | "none"                         | "none"                                                | ""        |
    | "GetStoragePoolListError"      | "induced error"                                       | ""        |
    | "none"                         | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test GetMaskingViewList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a MaskingView <mvname>
    And I induce error <induced>
    When I call GetMaskingViewList
    Then the error message contains <errormsg>
    And I get a valid MaskingViewList if no error

    Examples:
    | induced                        | errormsg                              | mvname            | arrays    |
    | "none"                         | "none"                                | "CSI-Test-MV"     | ""        |
    | "GetMaskingViewError"          | "induced error"                       | "CSI-Test-MV"     | ""        |
    | "none"                         | "ignored as it is not managed"        | "CSI-Test-MV"     | "ignored" |

  Scenario Outline: Test GetMaskingViewByID
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a MaskingView <mvname>
    And I induce error <induced>
    When I call GetMaskingViewByID <mvname>
    Then the error message contains <errormsg>
    And I get a valid MaskingView if no error

    Examples:
    | mvname                | induced                        | errormsg                                              | arrays    |
    | "Test-MV"             | "none"                         | "none"                                                | ""        |
    | "Test-MV"             | "GetMaskingViewError"          | "induced error"                                       | ""        |
    | "Test-MV"             | "none"                         | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test Rename Masking View
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a MaskingView <mvname>
    And I induce error <induced>
    When I call RenameMaskingView with <newname>
    Then the error message contains <errormsg>
    And I get a valid MaskingView if no error

    Examples:
    | mvname            | newname              | induced                    | errormsg                             | arrays    |
    | "CSI-Test-MV"     | "Renamed"            | "none"                     | "none"                               | ""        |
    | "CSI-Test-MV"     | "Renamed"            | "UpdateMaskingViewError"   | "induced error"                      | ""        |
    | "CSI-Test-MV"     | "Renamed"            | "none"                     | "ignored as it is not managed"       | "ignored" |

  Scenario Outline: Test DeleteMaskingView
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a MaskingView <mvname>
    And I induce error <induced>
    When I call DeleteMaskingView
    Then the error message contains <errormsg>

    Examples:
    | induced                        | errormsg                          | mvname                 | arrays    |
    | "none"                         | "none"                            | "CSI-Test-MV"          | ""        |
    | "DeleteMaskingViewError"       | "induced error"                   | "CSI-Test-MV"          | ""        |
    | "none"                         | "ignored as it is not managed"    | "CSI-Test-MV"          | "ignored" |

  Scenario Outline: Test GetPortGroupList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a PortGroup
    And I induce error <induced>
    When I call GetPortGroupList
    Then the error message contains <errormsg>
    And I get a valid PortGroupList if no error

    Examples:
    | induced                        | errormsg                                              | arrays    |
    | "none"                         | "none"                                                | ""        |
    | "GetPortGroupError"            | "induced error"                                       | ""        |
    | "none"                         | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test GetPortGroupByID
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a PortGroup
    And I induce error <induced>
    When I call GetPortGroupByID
    Then the error message contains <errormsg>
    And I get a valid PortGroup if no error

    Examples:
    | induced                        | errormsg                                              | arrays    |
    | "none"                         | "none"                                                | ""        |
    | "GetPortGroupError"            | "induced error"                                       | ""        |
    | "none"                         | "ignored as it is not managed"                        | "ignored" |

Scenario Outline: Test CreatePortGroup
  Given a valid connection
  And I induce error <induced>
  And I call CreatePortGroup <groupname> with ports <initialPorts>
  Then the error message contains <errormsg>
  And I get PortGroup <groupname> if no error
  Then I expect PortGroup to have these ports <finalPorts>

  Examples:
    | groupname             | initialPorts          | finalPorts            | induced                | errormsg        |
    | "Test-CreatePG"       | "SE-1E:000,SE-2E:001" | "SE-1E:000,SE-2E:001" | "none"                 | "none"          |
    | "Test-CreatePG-error" | "SE-1E:000,SE-2E:001" | ""                    | "CreatePortGroupError" | "induced error" |


Scenario Outline: Test UpdatePortGroup
  Given a valid connection
  And I induce error <induced>
  And I call CreatePortGroup <groupname> with ports <initialPorts>
  And I get PortGroup <groupname> if no error
  When I call UpdatePortGroup <groupname> with ports <updatedPorts>
  Then the error message contains <errormsg>
  Then I expect PortGroup to have these ports <finalPorts>

  Examples:
    | groupname             | initialPorts                              | updatedPorts                              | finalPorts                                | induced                | errormsg        |
    | "Test-UpdatePG1"      | "SE-1E:000,SE-2E:001"                     | "SE-1E:000,SE-2A:002"                     | "SE-1E:000,SE-2A:002"                     | "none"                 | "none"          |
    | "Test-UpdatePG2"      | "SE-1E:000,SE-2E:001"                     | "SE-1E:000,SE-2E:001"                     | "SE-1E:000,SE-2E:001"                     | "none"                 | "none"          |
    | "Test-UpdatePG3"      | "SE-1E:000,SE-2E:001"                     | "SE-4E:000,SE-3E:000"                     | "SE-4E:000,SE-3E:000"                     | "none"                 | "none"          |
    | "Test-UpdatePG4"      | "SE-1E:000,SE-2E:001,SE-4E:000,SE-3E:000" | "SE-4E:000,SE-3E:000"                     | "SE-4E:000,SE-3E:000"                     | "none"                 | "none"          |
    | "Test-UpdatePG5"      | "SE-1E:000,SE-2E:001"                     | "SE-1E:000,SE-2E:001,SE-4E:000,SE-3E:000" | "SE-1E:000,SE-2E:001,SE-4E:000,SE-3E:000" | "none"                 | "none"          |
    | "Test-UpdatePG-error" | "SE-1E:000,SE-2E:001"                     | "SE-1E:000,SE-2A:002"                     | ""                                        | "UpdatePortGroupError" | "induced error" |

Scenario Outline: Test RenamePortGroup
    Given a valid connection
    And I induce error <induced>
    And I call CreatePortGroup <groupname> with ports <initialPorts>
    And I get PortGroup <groupname> if no error
    And I call RenamePortGroup with <newname>
    Then the error message contains <errormsg>
    And I get PortGroup <newname> if no error    

    Examples:
    | groupname                | initialPorts            | newname                       | induced                    | errormsg                 |           
    | "Test-UpdatePG11"        | "SE-1E:000,SE-2E:001"   | "Renamed1"                    | "none"                     | "none"                   |
    | "Test-UpdatePG21-error"  | "SE-1E:000,SE-2E:001"   | "Renamed1"                    | "UpdatePortGroupError"     | "induced error"          |
                               
Scenario Outline: Test DeletePortGroup
  Given a valid connection
  And I induce error <induced>
  And I call CreatePortGroup <groupname> with ports <initialPorts>
  And I get PortGroup <groupname> if no error
  Then I call DeletePortGroup <groupname>
  Then the error message contains <errormsg>
  And the PortGroup <groupname> should not exist

  Examples:
    | groupname             | initialPorts          | induced                | errormsg        |
    | "Test-DeletePG"       | "SE-1E:000,SE-2E:001" | "none"                 | "none"          |
    | "Test-DeletePG-error" | "SE-1E:000,SE-2E:001" | "DeletePortGroupError" | "induced error" |

Scenario Outline: Test GetHostList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a FC Host <fchostname>
    And I have a ISCSI Host <hostname>
    And I have a NVMeTCP Host <nvmehostname>
    And I induce error <induced>
    When I call GetHostList
    Then the error message contains <errormsg>
    And I get a valid HostList if no error

    Examples:
    | fchostname     | hostname     | nvmehostname     | induced                        | errormsg                                              | arrays    |
    | "Test-Host-FC" | "Test-Host"  | "Test-Host-NVME" | "none"                         | "none"                                                | ""        |
    | "Test-Host-FC" | "Test-Host"  | "Test-Host-NVME" | "GetHostError"                 | "induced error"                                       | ""        |
    | "Test-Host-FC" | "Test-Host"  | "Test-Host-NVME" | "none"                         | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test GetHostByID
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a ISCSI Host <hostname>
    And I induce error <induced>
    When I call GetHostByID <hostname>
    Then the error message contains <errormsg>
    And I get a valid Host if no error

    Examples:
    | hostname     | induced                        | errormsg                                              | arrays    |
    | "Test-Host"  | "none"                         | "none"                                                | ""        |
    | "Test-Host"  | "GetHostError"                 | "induced error"                                       | ""        |
    | "Test-Host"  | "none"                         | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test CreateHost
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateHost <hostname>
    Then the error message contains <errormsg>
    And I get a valid Host if no error

    Examples:
    | hostname       | induced                        | errormsg                                              | arrays    |
    | "Test-Host"    | "none"                         | "none"                                                | ""        |
    | "Test-Host"    | "CreateHostError"              | "induced error"                                       | ""        |
    | "Test-Host"    | "none"                         | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test UpdateHost
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateHost <hostname>
    And I get a valid Host if no error
    When I call UpdateHost
    Then the error message contains <errormsg>
    And I get a valid Host if no error

    Examples:
    | hostname       | induced                        | errormsg                                              | arrays    |
    | "Test-Host"    | "none"                         | "none"                                                | ""        |
    | "Test-Host"    | "UpdateHostError"              | "induced error"                                       | ""        |
    | "Test-Host"    | "none"                         | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test UpdateHostFlags
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateHost <hostname>
    And I get a valid Host if no error
    When I call UpdateHostFlags
    Then the error message contains <errormsg>
    And I get a valid Host if no error

    Examples:
    | hostname       | induced                        | errormsg                                              | arrays    |
    | "Test-Host"    | "none"                         | "none"                                                | ""        |
    | "Test-Host"    | "UpdateHostError"              | "induced error"                                       | ""        |

  Scenario Outline: Test DeleteHost
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateHost <hostname>
    And I get a valid Host if no error
    When I call DeleteHost <hostname>
    Then the error message contains <errormsg>

    Examples:
    | hostname            | induced                        | errormsg                                              | arrays    |
    | "Test-Host"         | "none"                         | "none"                                                | ""        |
    | "Test-Host"         | "DeleteHostError"              | "induced error"                                       | ""        |
    | "Test-Host"         | "none"                         | "ignored as it is not managed"                        | "ignored" |

    Scenario Outline: Test GetInitiatorList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a Initiator
    And I induce error <induced>
    When I call GetInitiatorList
    Then the error message contains <errormsg>
    And I get a valid InitiatorList if no error

    Examples:
    | induced                        | errormsg                                              | arrays    |
    | "none"                         | "none"                                                | ""        |
    | "GetInitiatorError"            | "induced error"                                       | ""        |
    | "none"                         | "ignored as it is not managed"                        | "ignored" |

    Scenario Outline: Test GetInitiatorList with filters
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a Initiator
    When I call GetInitiatorList with filters
    Then the error message contains <errormsg>

    Examples:
    |errormsg          | arrays    |
    | "none"           | ""        |

  Scenario Outline: Test GetInitiatorByID
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a Initiator
    And I induce error <induced>
    When I call GetInitiatorByID
    Then the error message contains <errormsg>
    And I get a valid Initiator if no error

    Examples:
    | induced                        | errormsg                                              | arrays    |
    | "none"                         | "none"                                                | ""        |
    | "GetInitiatorError"            | "induced error"                                       | ""        |
    | "none"                         | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test cases for CreateMaskingViewWithHost
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a ISCSI Host <hostname>
    And I have a PortGroup
    And I have a StorageGroup <sgname>
    And I induce error <induced>
    When I call CreateMaskingViewWithHost <mvname>
    Then the error message contains <errormsg>
    And I get a valid MaskingView if no error

    Examples:
    | hostname     | sgname      | mvname         | induced                      | errormsg                                              | arrays    |
    | "TestHost"   | "TestSG"    | "TestMV"       | "none"                       | "none"                                                | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "CreateMaskingViewError"     | "Failed to create masking view"                       | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "MaskingViewAlreadyExists"   | "The requested masking view resource already exists"  | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "PortGroupNotFoundError"     | "Port Group on Symmetrix cannot be found"             | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "InitiatorGroupNotFoundError"| "Initiator Group on Symmetrix cannot be found"        | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "StorageGroupNotFoundError"  | "Storage Group on Symmetrix cannot be found"          | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "none"                       | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test cases for CreateMaskingViewWithHost
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a NVMeTCP Host <hostname>
    And I have a PortGroup
    And I have a StorageGroup <sgname>
    And I induce error <induced>
    When I call CreateMaskingViewWithHost <mvname>
    Then the error message contains <errormsg>
    And I get a valid MaskingView if no error

    Examples:
    | hostname     | sgname      | mvname         | induced                      | errormsg                                              | arrays    |
    | "TestHost"   | "TestSG"    | "TestMV"       | "none"                       | "none"                                                | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "CreateMaskingViewError"     | "Failed to create masking view"                       | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "MaskingViewAlreadyExists"   | "The requested masking view resource already exists"  | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "PortGroupNotFoundError"     | "Port Group on Symmetrix cannot be found"             | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "InitiatorGroupNotFoundError"| "Initiator Group on Symmetrix cannot be found"        | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "StorageGroupNotFoundError"  | "Storage Group on Symmetrix cannot be found"          | ""        |
    | "TestHost"   | "TestSG"    | "TestMV"       | "none"                       | "ignored as it is not managed"                        | "ignored" |


  Scenario Outline: Test cases for CreateMaskingViewWithHostGroup
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a HostGroup <hostname>
    And I have a PortGroup
    And I have a StorageGroup <sgname>
    And I induce error <induced>
    When I call CreateMaskingViewWithHostGroup <mvname>
    Then the error message contains <errormsg>
    And I get a valid MaskingView if no error
    Examples:
    | hostname     | sgname      | mvname         | induced                      | errormsg                                              | arrays    |
    | "TestHostGrp"| "TestSG"    | "TestMV"       | "none"                       | "none"                                                | ""        |
    | "TestHostGrp"| "TestSG"    | "TestMV"       | "InitiatorGroupNotFoundError"| "Initiator Group on Symmetrix cannot be found"        | ""        |
    | "TestHostGrp"| "TestSG"    | "TestMV"       | "none"                       | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test cases for Asynchronous AddVolumesToStorageGroup
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a StorageGroup <sgname>
    And I have <nvols> volumes
    And I induce error <induced>
    When I call AddVolumesToStorageGroup <sgname>
    Then the error message contains <errormsg>
    And then the Volumes are part of StorageGroup if no error
    Examples:
    | nvols | sgname        |induced                   | errormsg                                                 | arrays    |
    | 5     | "TestSG"      |"none"                    | "none"                                                   | ""        |
    | 1     | "TestSG"      |"none"                    | "none"                                                   | ""        |
    | 0     | "TestSG"      |"none"                    | "At least one volume id has to be specified"             | ""        |
    | 5     | "TestSG"      |"VolumeNotAddedError"     | "A job was not returned from UpdateStorageGroup"         | ""        |
    | 3     | "TestSG"      |"UpdateStorageGroupError" | "A job was not returned from UpdateStorageGroup"         | ""        |
    | 1     | "TestSG"      |"JobFailedError"          | "UpdateStorageGroup job failed"                          | ""        |
    | 1     | "TestSG"      |"GetJobError"             | "induced error"                                          | ""        |
    | 1     | "TestSG"      |"none"                    | "ignored as it is not managed"                           | "ignored" |

  Scenario Outline: Test cases for Synchronous AddVolumesToStorageGroup
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a StorageGroup <sgname>
    And I have <nvols> volumes
    And I induce error <induced>
    When I call AddVolumesToStorageGroupS <sgname>
    Then the error message contains <errormsg>
    And then the Volumes are part of StorageGroup if no error
    Examples:
    | nvols | sgname        |induced                   | errormsg                                          | arrays    |
    | 5     | "TestSG"      |"none"                    | "none"                                            | ""        |
    | 1     | "TestSG"      |"none"                    | "none"                                            | ""        |
    | 0     | "TestSG"      |"none"                    | "at least one volume id has to be specified"      | ""        |
    | 5     | "TestSG"      |"VolumeNotAddedError"     | "Error adding volume to the SG"                   | ""        |
    | 3     | "TestSG"      |"UpdateStorageGroupError" | "Error updating Storage Group: induced error"     | ""        |
    | 1     | "TestSG"      |"none"                    | "ignored as it is not managed"                    | "ignored" |

  Scenario Outline: Test cases for Asynchronous AddVolumesToStorageGroup for v91
    Given a valid v91 connection
    And I have an allowed list of <arrays>
    And I have a StorageGroup <sgname>
    And I have <nvols> volumes
    And I induce error <induced>
    When I call AddVolumesToStorageGroup <sgname>
    Then the error message contains <errormsg>
    And then the Volumes are part of StorageGroup if no error
    Examples:
    | nvols | sgname        |induced                   | errormsg                                                 | arrays    |
    | 5     | "TestSG"      |"none"                    | "none"                                                   | ""        |
    | 1     | "TestSG"      |"none"                    | "none"                                                   | ""        |
    | 0     | "TestSG"      |"none"                    | "At least one volume id has to be specified"             | ""        |
    | 5     | "TestSG"      |"VolumeNotAddedError"     | "A job was not returned from UpdateStorageGroup"         | ""        |
    | 3     | "TestSG"      |"UpdateStorageGroupError" | "A job was not returned from UpdateStorageGroup"         | ""        |
    | 1     | "TestSG"      |"JobFailedError"          | "UpdateStorageGroup job failed"                          | ""        |
    | 1     | "TestSG"      |"GetJobError"             | "induced error"                                          | ""        |
    | 1     | "TestSG"      |"none"                    | "ignored as it is not managed"                           | "ignored" |

  Scenario Outline: Test cases for Synchronous AddVolumesToStorageGroup for v91
    Given a valid v91 connection
    And I have an allowed list of <arrays>
    And I have a StorageGroup <sgname>
    And I have <nvols> volumes
    And I induce error <induced>
    When I call AddVolumesToStorageGroupS <sgname>
    Then the error message contains <errormsg>
    And then the Volumes are part of StorageGroup if no error
    Examples:
    | nvols | sgname        |induced                   | errormsg                                          | arrays    |
    | 5     | "TestSG"      |"none"                    | "none"                                            | ""        |
    | 1     | "TestSG"      |"none"                    | "none"                                            | ""        |
    | 0     | "TestSG"      |"none"                    | "at least one volume id has to be specified"      | ""        |
    | 5     | "TestSG"      |"VolumeNotAddedError"     | "Error adding volume to the SG"                   | ""        |
    | 3     | "TestSG"      |"UpdateStorageGroupError" | "Error updating Storage Group: induced error"     | ""        |
    | 1     | "TestSG"      |"none"                    | "ignored as it is not managed"                    | "ignored" |

  Scenario Outline: Test case for retriving list of target IP addresses
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call GetListOfTargetAddresses
    Then the error message contains <errormsg>
    And I recieve <count> IP addresses
    Examples:
    | count | induced                   | errormsg                                                 | arrays    |
    | 8     | "none"                    | "none"                                                   | ""        |
    | 0     | "GetPortError"            | "none"                                                   | ""        |
    | 0     | "GetDirectorError"        | "Error retrieving Director"                              | ""        |
    | 0     | "none"                    | "ignored as it is not managed"                           | "ignored" |

  Scenario Outline: Test Array allowed lists
    Given a valid connection
    And I have an allowed list of <arrays>
    And it contains <count> arrays
    And should include <included>
    And should not include <excluded>
    Examples:
    | arrays        | count | included       | excluded     |
    | ""            | 0     | "1,2,3"        | ""           |
    | "1"           | 1     | "1"            | "2"          |
    | "1,2,3,4"     | 4     | "1,2,3,4"      | "8,9"        |

  Scenario Outline: TestCases for GetSymmetrixIDList with an allowed list of arrays
    Given a valid connection
    And I have an allowed list of <arrays>
    When I call GetSymmetrixIDList
    Then I get a valid Symmetrix ID List that contains <included> and does not contains <excluded>
    Examples:
    | arrays                        | included                       | excluded         | explanation                                      |
    | ""                            | "000197802104, 000197900046"   | ""               | an empty allowed list will allow any array       |
    | "000197900046"                | "000197900046"                 | "000197802104"   | including one specific array will exclude others |
    | "000197802104, 999999999999"  | "000197802104"                 | "999999999999"   | make sure that non existent arrays are not found |

  Scenario Outline: Get Symmetrix System with an allowed list of arrays
    Given a valid connection
    And I have an allowed list of <arrays>
    When I call GetSymmetrixByID <id>
    Then the error message contains <errormsg>
    And I get a valid Symmetrix Object if no error
    Examples:
    | id              | arrays                | errormsg                         |
    | "000197900046"  | ""                    | "none"                           |
    | "000000000000"  | "none"                | "ignored as it is not managed"   |
    | "000197900046"  | "000197900046"        | "none"                           |
    | "000197900046"  | "000197802104"        | "ignored as it is not managed"   |

  Scenario Outline: Refresh Symmetrix system
    Given a valid connection
    And I have an allowed list of <allowedarrays>
    When I call RefreshSymmetrix <array>
    Then the error message contains <errormsg>
    Examples:
    | allowedarrays   | array               | errormsg                            |
    | "000197900046"  |  "000197900046"     | "none"                              | 
    | "000197900046"  |  "000011112222"     |"is ignored as it is not managed"    |


  Scenario Outline: Get ISCSI targets
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call GetISCSITargets
    Then the error message contains <errormsg>
    And I recieve <count> targets
    Examples:
    | arrays           | induced                   | errormsg                         | count |
    | "000000000000"   | "none"                    | "ignored as it is not managed"   | 0     |
    | "000197900046"   | "GetDirectorError"        | "Error retrieving Director"      | 0     |
    | "000197900046"   | "GetPortGigEError"        | "none"                           | 0     |
    | "000197900046"   | "GetPortISCSITargetError" | "Error retrieving ISCSI targets" | 0     |
    | "000197900046"   | "GetSpecificPortError"    | "none"                           | 0     |
    | "000197900046"   | "none"                    | "none"                           | 8     |

  Scenario Outline: Get NVMeTCP targets
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call GetNVMeTCPTargets
    Then the error message contains <errormsg>
    And I recieve <count> targets
    Examples:
    | arrays           | induced                     | errormsg                           | count |
    | "000197900046"   | "none"                      | "none"                             | 0     |
    | "000000000000"   | "none"                      | "ignored as it is not managed"     | 0     |
    | "000197900046"   | "GetDirectorError"          | "Error retrieving Director"        | 0     |
    | "000197900046"   | "GetPortGigEError"          | "none"                             | 0     |
    | "000197900046"   | "GetSpecificPortError"      | "none"                             | 0     |
    | "000197900046"   | "GetPortNVMeTCPTargetError" | "Error retrieving NVMeTCP targets" | 0     |
   

  Scenario Outline: Test UpdateHostName
      Given a valid connection
      And I have an allowed list of <arrays>
      And I induce error <induced>
      When I call CreateHost <hostname>
      And I get a valid Host if no error
      When I call UpdateHostName <newname>
      Then the error message contains <errormsg>
      And I get a valid Host if no error

      Examples:
      | hostname       | newname      | induced                        | errormsg                                              | arrays    |
      | "Test-Host"    | "Test-Host"  | "none"                         | "none"                                                | ""        |
      | "Test-Host"    | "Test-Host"  | "UpdateHostError"              | "induced error"                                       | ""        |
      | "Test-Host"    | "Test-Host"  | "none"                         | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test CreateHostGroup
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateHostGroup <hostGroupname> with flags <setHostFlags>
    Then the error message contains <errormsg>
    And I get a valid HostGroup if no error

    Examples:
    | hostGroupname       | induced                        | errormsg                                   | arrays    | setHostFlags |
    | "Test-HostGroup"    | "none"                         | "none"                                     | ""        | "false"      |
    | "Test-HostGroup"    | "none"                         | "none"                                     | ""        | "true"       |
    | "Test-HostGroup"    | "CreateHostGroupError"         | "induced error"                            | ""        | "false"      |
    | "Test-HostGroup"    | "none"                         | "ignored as it is not managed"             | "ignored" | "false"      |

  Scenario Outline: Test GetHostGroupByID
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a valid HostGroup <hostGroupname>
    And I induce error <induced>
    When I call GetHostGroupByID <hostGroupname>
    Then the error message contains <errormsg>
    And I get a valid HostGroup if no error

    Examples:
    | hostGroupname     | induced                        | errormsg                                     | arrays    |
    | "Test-HostGroup"  | "none"                         | "none"                                       | ""        |
    | "Test-HostGroup"  | "GetHostGroupError"            | "induced error"                              | ""        |
    | "Test-HostGroup"  | "none"                         | "ignored as it is not managed"               | "ignored" |

  Scenario Outline: Test UpdateHostGroupHosts
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateHostGroup <hostGroupname> with flags <setHostFlags>
    And I get a valid HostGroup if no error
    When I call UpdateHostGroupHosts <hostGroupname> with hosts <hostsIDs>
    Then the error message contains <errormsg>
    And I get a valid HostGroup if no error

    Examples:
    | hostGroupname       | induced                        | errormsg                              | arrays    | hostsIDs     | setHostFlags |
    | "Test-HostGroup"    | "none"                         | "none"                                | ""        | "testHostID" | "true"       |
    | "Test-HostGroup"    | "UpdateHostGroupError"         | "induced error"                       | ""        | ""           | "false"      |
    | "Test-HostGroup"    | "none"                         | "ignored as it is not managed"        | "ignored" | "testHostID" | "false"      |

  Scenario Outline: Test UpdateHostGroupFlags
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateHostGroup <hostGroupname> with flags <setHostFlags>
    And I get a valid HostGroup if no error
    When I call UpdateHostGroupFlags <hostGroupname> with flags <updateHostFlags>
    Then the error message contains <errormsg>
    And I get a valid HostGroup if no error

    Examples:
    | hostGroupname       | induced                        | errormsg                  | arrays    | setHostFlags | updateHostFlags |
    | "Test-HostGroup"    | "none"                         | "none"                    | ""        | "false"      | "true"          |
    | "Test-HostGroup"    | "UpdateHostGroupError"         | "induced error"           | ""        | "false"      | "false"         |

  Scenario Outline: Test UpdateHostGroupName
      Given a valid connection
      And I have an allowed list of <arrays>
      And I induce error <induced>
      When I call CreateHostGroup <hostGroupname> with flags <setHostFlags>
      And I get a valid HostGroup if no error
      When I call UpdateHostGroupName <newhostGroupName>
      Then the error message contains <errormsg>
      And I get a valid HostGroup if no error

      Examples:
      | hostGroupname       | newhostGroupName     | induced                        | errormsg                            | arrays    | setHostFlags |
      | "Test-HostGroup"    | "Test-HostGroup123"  | "none"                         | "none"                              | ""        | "true"       |
      | "Test-HostGroup"    | "Test-HostGroup123"  | "UpdateHostGroupError"         | "induced error"                     | ""        | "false"      |
      | "Test-HostGroup"    | "Test-HostGroup123"  | "none"                         | "ignored as it is not managed"      | "ignored" | "true"       |

  Scenario Outline: Test DeleteHostGroup
    Given a valid connection
    And I have an allowed list of <arrays>
    And I induce error <induced>
    When I call CreateHostGroup <hostGroupname> with flags <setHostFlags>
    And I get a valid HostGroup if no error
    When I call DeleteHostGroup <hostGroupname>
    Then the error message contains <errormsg>

    Examples:
    | hostGroupname           | induced                 | errormsg                              | arrays    | setHostFlags |
    | "Test-HostGroup"        | "none"                  | "none"                                | ""        | "true"       |
    | "Test-HostGroup"        | "DeleteHostGroupError"  | "induced error"                       | ""        | "false"      |
    | "Test-HostGroup"        | "none"                  | "ignored as it is not managed"        | "ignored" | "false"      |

  Scenario Outline: Test GetHostGroupList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a valid HostGroup <hostGroupname>
    And I induce error <induced>
    When I call GetHostGroupList
    Then the error message contains <errormsg>
    And I get a valid HostGroupList if no error

    Examples:
    | hostGroupname         | induced                            | errormsg                                              | arrays    |
    | "Test-HostGroup-1"    | "none"                             | "none"                                                | ""        |
    | "Test-HostGroup-1"    | "GetHostGroupListError"            | "induced error"                                       | ""        |
    | "Test-HostGroup-1"    | "none"                             | "ignored as it is not managed"                        | "ignored" |

  Scenario Outline: Test cases for GetVolumeIDList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have <nvols> volumes
    And I induce error <induced>
    When I call GetVolumeIDListWithParams
    Then the error message contains <errormsg>
    And I get a valid VolumeIDList with <vols> if no error

    Examples:
    | nvols      | vols  | induced                    | errormsg                      | arrays    |
    | 7          | 7     | "none"                     | "none"                        | ""        |
    | 5          | 5     | "none"                     | "ignored as it is not managed"| "ignore"  |

  Scenario Outline: Test cases for GetStorageGroupVolumeIDList
    Given a valid connection
    And I have an allowed list of <arrays>
    And I have a StorageGroup <sgname>
    And I have <nvols> volumes
    And I induce error <induced>
    When I call GetVolumeIDListInStorageGroup <sgname>
    Then the error message contains <errormsg>
    And I get a valid VolumeIDList with <vols> if no error

    Examples:
    | nvols      | vols  | induced  | errormsg                      | arrays          | sgname      | 
    | 7          | 7     | "none"   | "none"                        | "000197900046"  | "TestSG"    | 
    | 5          | 5     | "none"   | "storageGroupID is empty"     | "000197900046"  | ""          |
