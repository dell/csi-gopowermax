package v100

import "testing"

func TestGetJobResource(t *testing.T) {
	tests := []struct {
		name                 string
		job                  Job
		expectedSymmetrixID  string
		expectedResourceType string
		expectedResourceID   string
	}{
		{
			name: "valid resource link",
			job: Job{
				ResourceLink: "provisioning/system/SYMMETRIX-1234/volume/1234",
			},
			expectedSymmetrixID:  "SYMMETRIX-1234",
			expectedResourceType: "volume",
			expectedResourceID:   "1234",
		},
		{
			name: "Empty resource link",
			job: Job{
				ResourceLink: "",
			},
			expectedSymmetrixID:  "",
			expectedResourceType: "",
			expectedResourceID:   "",
		},
		{
			name: "invalid resource link",
			job: Job{
				ResourceLink: "system/SYMMETRIX-1234",
			},
			expectedSymmetrixID:  "",
			expectedResourceType: "",
			expectedResourceID:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			symID, resourceType, id := tt.job.GetJobResource()

			if symID != tt.expectedSymmetrixID {
				t.Errorf("expected %s, got %s", tt.expectedSymmetrixID, symID)
			}

			if resourceType != tt.expectedResourceType {
				t.Errorf("expected %s, got %s", tt.expectedResourceType, resourceType)
			}

			if id != tt.expectedResourceID {
				t.Errorf("expected %s, got %s", tt.expectedResourceID, id)
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
	}{
		{
			name: "valid error",
			err: &Error{
				Message: "test-error",
			},
		},
		{
			name: "empty error",
			err: &Error{
				Message: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.err.Message {
				t.Errorf("expected %s, got %s", tt.err.Message, tt.err.Error())
			}
		})
	}
}
