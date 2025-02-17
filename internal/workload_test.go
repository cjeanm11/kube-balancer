package internal

import (
	"testing"
)

func TestWorkload_Validate(t *testing.T) {
	tests := []struct {
		name    string
		workload Workload
		wantErr bool
	}{
		{
			name: "Valid workload",
			workload: Workload{
				Name:     "ValidTask",
				CPU:      10,
				Memory:   512,
				Priority: 1,
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			workload: Workload{
				Name:     "",
				CPU:      10,
				Memory:   512,
				Priority: 1,
			},
			wantErr: true,
		},
		{
			name: "Zero CPU",
			workload: Workload{
				Name:     "Task",
				CPU:      0,
				Memory:   512,
				Priority: 1,
			},
			wantErr: true,
		},
		{
			name: "Zero Memory",
			workload: Workload{
				Name:     "Task",
				CPU:      10,
				Memory:   0,
				Priority: 1,
			},
			wantErr: true,
		},
		{
			name: "Zero Priority",
			workload: Workload{
				Name:     "Task",
				CPU:      10,
				Memory:   512,
				Priority: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.workload.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Workload.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkload_String(t *testing.T) {
	workload := Workload{
		Name:     "test-workload",
		CPU:      2,
		Memory:   2048,
		Priority: 3,
	}

	expected := "Workload{Name: test-workload, CPU: 2, Memory: 2048MB, Priority: 3}"
	if got := workload.String(); got != expected {
		t.Errorf("Workload.String() = %v, want %v", got, expected)
	}
} 