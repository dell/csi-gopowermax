/*
Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package pmax

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	types "github.com/dell/gopowermax/v2/types/v100"
)

const (
	urlPrefix = "/univmax/restapi/100/"
)

func TestModifyMigrationSession(t *testing.T) {
	type testCase struct {
		server         *httptest.Server
		localSymID     string
		storageGroupID string
		expectedErr    error
	}

	cases := map[string]testCase{
		"get one device success": {
			localSymID:     "mock-local-sym-id",
			storageGroupID: "mock-storage-group-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				url := fmt.Sprintf("%s%s%s%s%s/%s", urlPrefix, XMigration, SymmetrixX, "mock-local-sym-id", XStorageGroup, "mock-storage-group-id")
				switch req.RequestURI {
				case url:
					resp.WriteHeader(http.StatusOK)
					response := &types.ModifyMigrationSessionRequest{
						Action:          "mock-action",
						ExecutionOption: types.ExecutionOptionSynchronous,
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					_, err = resp.Write(content)
					if err != nil {
						t.Fatal(err)
					}
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			localSymID:     "mock-local-sym-id",
			storageGroupID: "mock-storage-group-id-2",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
		"invalid array": {
			localSymID: "invalid-array-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("the requested array (invalid-array-id) is ignored as it is not managed"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", true, true, "")
		if err != nil {
			t.Fatal(err)
		}
		client.SetAllowedArrays([]string{"mock-local-sym-id"})
		err = client.ModifyMigrationSession(context.TODO(), tc.localSymID, "mock-action", tc.storageGroupID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
		tc.server.Close()
	}
}

func TestCreateMigrationEnvironment(t *testing.T) {
	type testCase struct {
		server         *httptest.Server
		localSymID     string
		storageGroupID string
		expectedErr    error
	}

	cases := map[string]testCase{
		"get one device success": {
			localSymID: "mock-local-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				url := fmt.Sprintf("%s%s%s%s", urlPrefix, XMigration, SymmetrixX, "mock-local-sym-id")
				switch req.RequestURI {
				case url:
					resp.WriteHeader(http.StatusOK)
					response := &types.CreateMigrationEnv{
						OtherArrayID:    "mock-storage-group-id",
						ExecutionOption: types.ExecutionOptionSynchronous,
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					_, err = resp.Write(content)
					if err != nil {
						t.Fatal(err)
					}
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			localSymID: "mock-local-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
		"invalid array": {
			localSymID: "invalid-array-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("the requested array (invalid-array-id) is ignored as it is not managed"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", true, true, "")
		if err != nil {
			t.Fatal(err)
		}
		client.SetAllowedArrays([]string{"mock-local-sym-id"})
		_, err = client.CreateMigrationEnvironment(context.TODO(), tc.localSymID, tc.storageGroupID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
		tc.server.Close()
	}
}

func TestDeleteMigrationEnvironment(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		localSymID  string
		remoteSymID string
		expectedErr error
	}

	cases := map[string]testCase{
		"get one device success": {
			localSymID:  "mock-local-sym-id",
			remoteSymID: "mock-remote-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				url := fmt.Sprintf("%s%s%s%s", urlPrefix, XMigration, SymmetrixX, "mock-local-sym-id")
				switch req.RequestURI {
				case url:
					resp.WriteHeader(http.StatusOK)
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			localSymID: "mock-local-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
		"invalid array": {
			localSymID: "invalid-array-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("the requested array (invalid-array-id) is ignored as it is not managed"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", true, true, "")
		if err != nil {
			t.Fatal(err)
		}
		client.SetAllowedArrays([]string{"mock-local-sym-id"})
		err = client.DeleteMigrationEnvironment(context.TODO(), tc.localSymID, tc.remoteSymID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
		tc.server.Close()
	}
}

func TestCreateSGMigrationByID(t *testing.T) {
	type testCase struct {
		server         *httptest.Server
		localSymID     string
		remoteSymID    string
		storageGroupID string
		expectedErr    error
	}

	cases := map[string]testCase{
		"get one device success": {
			localSymID:     "mock-local-sym-id",
			remoteSymID:    "mock-remote-sym-id",
			storageGroupID: "mock-storage-group-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				url := fmt.Sprintf("%s%s%s%s%s/%s", urlPrefix, XMigration, SymmetrixX, "mock-local-sym-id", XStorageGroup, "mock-storage-group-id")
				switch req.RequestURI {
				case url:
					resp.WriteHeader(http.StatusOK)
					response := types.CreateMigrationEnv{
						OtherArrayID:    "mock-remote-sym-id",
						ExecutionOption: types.ExecutionOptionSynchronous,
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					_, err = resp.Write(content)
					if err != nil {
						t.Fatal(err)
					}
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			localSymID: "mock-local-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
		"invalid array": {
			localSymID: "invalid-array-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("the requested array (invalid-array-id) is ignored as it is not managed"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", true, true, "")
		if err != nil {
			t.Fatal(err)
		}
		client.SetAllowedArrays([]string{"mock-local-sym-id"})
		_, err = client.CreateSGMigration(context.TODO(), tc.localSymID, tc.remoteSymID, tc.storageGroupID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
		tc.server.Close()
	}
}

func TestMigrateStorageGroup(t *testing.T) {
	type testCase struct {
		server         *httptest.Server
		localSymID     string
		storageGroupID string
		srpID          string
		serviceLevel   string
		thickVolumes   bool
		expectedErr    error
	}

	cases := map[string]testCase{
		"get one device success": {
			localSymID:     "mock-local-sym-id",
			storageGroupID: "mock-storage-group-id",
			srpID:          "mock_SRP_1",
			serviceLevel:   "mock-service-level",
			thickVolumes:   true,
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				url := fmt.Sprintf("%s%s%s%s%s", urlPrefix, XMigration, SymmetrixX, "mock-local-sym-id", XStorageGroup)
				switch req.RequestURI {
				case url:
					resp.WriteHeader(http.StatusOK)

					sloParams := []types.SLOBasedStorageGroupParam{}
					var snapshotPolicies []string

					response := &types.CreateStorageGroupParam{
						StorageGroupID:            "mock-storage-group-id",
						SRPID:                     "mock_SRP_1",
						Emulation:                 "mock-emulation",
						ExecutionOption:           types.ExecutionOptionSynchronous,
						SLOBasedStorageGroupParam: sloParams,
						SnapshotPolicies:          snapshotPolicies,
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					_, err = resp.Write(content)
					if err != nil {
						t.Fatal(err)
					}
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			localSymID: "mock-local-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
		"invalid array": {
			localSymID: "invalid-array-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("the requested array (invalid-array-id) is ignored as it is not managed"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", true, true, "")
		if err != nil {
			t.Fatal(err)
		}
		client.SetAllowedArrays([]string{"mock-local-sym-id"})
		_, err = client.MigrateStorageGroup(context.TODO(), tc.localSymID, tc.storageGroupID, tc.srpID, tc.serviceLevel, tc.thickVolumes)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
		tc.server.Close()
	}
}

func TestGetStorageGroupMigrationByID(t *testing.T) {
	type testCase struct {
		server         *httptest.Server
		localSymID     string
		storageGroupID string
		expectedErr    error
	}

	cases := map[string]testCase{
		"get one device success": {
			localSymID:     "mock-local-sym-id",
			storageGroupID: "mock-storage-group-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				url := fmt.Sprintf("%s%s%s%s%s/%s", urlPrefix, XMigration, SymmetrixX, "mock-local-sym-id", XStorageGroup, "mock-storage-group-id")
				switch req.RequestURI {
				case url:
					resp.WriteHeader(http.StatusOK)
					response := &types.MigrationSession{
						SourceArray:  "mock-local-sym-id",
						TargetArray:  "mock-target-sym-id",
						StorageGroup: "mock-storage-group-id",
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					_, err = resp.Write(content)
					if err != nil {
						t.Fatal(err)
					}
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			localSymID: "mock-local-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
		"invalid array": {
			localSymID: "invalid-array-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("the requested array (invalid-array-id) is ignored as it is not managed"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", true, true, "")
		if err != nil {
			t.Fatal(err)
		}
		client.SetAllowedArrays([]string{"mock-local-sym-id"})
		_, err = client.GetStorageGroupMigrationByID(context.TODO(), tc.localSymID, tc.storageGroupID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
		tc.server.Close()
	}
}

func TestGetStorageGroupMigration(t *testing.T) {
	type testCase struct {
		server      *httptest.Server
		localSymID  string
		expectedErr error
	}

	cases := map[string]testCase{
		"get one device success": {
			localSymID: "mock-local-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				queryParam := fmt.Sprintf("%s=%s", IncludeMigrations, "true")
				url := fmt.Sprintf("%s%s%s%s%s%s", urlPrefix, XMigration, SymmetrixX, "mock-local-sym-id", XStorageGroup, queryParam)
				switch req.RequestURI {
				case url:
					resp.WriteHeader(http.StatusOK)
					response := &types.MigrationStorageGroups{
						StorageGroupIDList: []string{"mock-storage-group-id"},
						MigratingNameList:  []string{"mock-migrating-name"},
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					_, err = resp.Write(content)
					if err != nil {
						t.Fatal(err)
					}
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			localSymID: "mock-local-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
		"invalid array": {
			localSymID: "invalid-array-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("the requested array (invalid-array-id) is ignored as it is not managed"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", true, true, "")
		if err != nil {
			t.Fatal(err)
		}
		client.SetAllowedArrays([]string{"mock-local-sym-id"})
		_, err = client.GetStorageGroupMigration(context.TODO(), tc.localSymID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
		tc.server.Close()
	}
}

func TestGetMigrationEnvironment(t *testing.T) {
	type testCase struct {
		server         *httptest.Server
		localSymID     string
		remoteSystemID string
		expectedErr    error
	}

	cases := map[string]testCase{
		"get one device success": {
			localSymID:     "mock-local-sym-id",
			remoteSystemID: "mock-remote-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				url := fmt.Sprintf("%s%s%s%s%s%s", urlPrefix, XMigration, SymmetrixX, "mock-local-sym-id", XEnvironment, "mock-remote-sym-id")
				switch req.RequestURI {
				case url:
					resp.WriteHeader(http.StatusOK)
					response := &types.MigrationEnv{
						ArrayID:               "mock-local-sym-id",
						StorageGroupCount:     2,
						MigrationSessionCount: 2,
						Local:                 true,
					}

					content, err := json.Marshal(response)
					if err != nil {
						t.Fatal(err)
					}

					_, err = resp.Write(content)
					if err != nil {
						t.Fatal(err)
					}
				default:
					resp.WriteHeader(http.StatusNoContent)
				}
			})),
			expectedErr: nil,
		},
		"bad request": {
			localSymID: "mock-local-sym-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("bad request"),
		},
		"invalid array": {
			localSymID: "invalid-array-id",
			server: httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write([]byte(`{"message":"bad request","httpStatusCode":400,"errorCode":0}`))
			})),
			expectedErr: errors.New("the requested array (invalid-array-id) is ignored as it is not managed"),
		},
	}

	for _, tc := range cases {
		client, err := NewClientWithArgs(tc.server.URL, "", true, true, "")
		if err != nil {
			t.Fatal(err)
		}

		client.SetAllowedArrays([]string{"mock-local-sym-id"})
		_, err = client.GetMigrationEnvironment(context.TODO(), tc.localSymID, tc.remoteSystemID)
		if err != nil {
			if tc.expectedErr.Error() != err.Error() {
				t.Fatal(err)
			}
		}
		tc.server.Close()
	}
}
