package call

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	expected := &Storage{
		toProcess: make([]Meta, 0),
		statuses:  make(map[ID]int, 0),
		mu:        &sync.Mutex{},
	}
	assert.Equal(t, expected, NewStorage())
}

func TestStorage_AddToQueueBack(t *testing.T) {
	type fields struct {
		toProcess []Meta
		statuses  map[ID]int
		mu        *sync.Mutex
	}
	type args struct {
		in0  context.Context
		meta Meta
	}
	type expectedValues struct {
		toProcess []Meta
		err       error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expectedValues
	}{
		{
			name: "success, empty",
			fields: fields{
				toProcess: make([]Meta, 0),
				statuses:  nil,
				mu:        &sync.Mutex{},
			},
			args: args{
				in0: nil,
				meta: Meta{
					PhoneNumber:    "777-777-777",
					VirtualAgentID: "sdas2-sdsada-dsad-a-sdasd",
					ID:             "3",
				},
			},
			expectedValues: expectedValues{
				toProcess: []Meta{
					{
						PhoneNumber:    "777-777-777",
						VirtualAgentID: "sdas2-sdsada-dsad-a-sdasd",
						ID:             "3",
					},
				},
				err: nil,
			},
		},
		{
			name: "success, not empty",
			fields: fields{
				toProcess: []Meta{
					{
						PhoneNumber:    "888-888-888",
						VirtualAgentID: "aaaa-aaaa-aaa-3-ddd",
						ID:             "2",
					},
				},
				statuses: nil,
				mu:       &sync.Mutex{},
			},
			args: args{
				in0: nil,
				meta: Meta{
					PhoneNumber:    "777-777-777",
					VirtualAgentID: "sdas2-sdsada-dsad-a-sdasd",
					ID:             "3",
				},
			},
			expectedValues: expectedValues{
				toProcess: []Meta{
					{
						PhoneNumber:    "888-888-888",
						VirtualAgentID: "aaaa-aaaa-aaa-3-ddd",
						ID:             "2",
					},
					{
						PhoneNumber:    "777-777-777",
						VirtualAgentID: "sdas2-sdsada-dsad-a-sdasd",
						ID:             "3",
					},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				toProcess: tt.fields.toProcess,
				statuses:  tt.fields.statuses,
				mu:        tt.fields.mu,
			}
			ao := assert.New(t)
			actualErr := s.AddToQueueBack(tt.args.in0, tt.args.meta)
			ao.Equal(tt.expectedValues.err, actualErr)
			ao.Equal(tt.expectedValues.toProcess, s.toProcess)

		})
	}
}

func TestStorage_AddToQueueFront(t *testing.T) {
	type fields struct {
		toProcess []Meta
		statuses  map[ID]int
		mu        *sync.Mutex
	}
	type args struct {
		in0  context.Context
		meta Meta
	}
	type expectedValues struct {
		toProcess []Meta
		err       error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expectedValues
	}{
		{
			name: "success, empty",
			fields: fields{
				toProcess: make([]Meta, 0),
				statuses:  nil,
				mu:        &sync.Mutex{},
			},
			args: args{
				in0: nil,
				meta: Meta{
					PhoneNumber:    "777-777-777",
					VirtualAgentID: "sdas2-sdsada-dsad-a-sdasd",
					ID:             "3",
				},
			},
			expectedValues: expectedValues{
				toProcess: []Meta{
					{
						PhoneNumber:    "777-777-777",
						VirtualAgentID: "sdas2-sdsada-dsad-a-sdasd",
						ID:             "3",
					},
				},
				err: nil,
			},
		},
		{
			name: "success, not empty",
			fields: fields{
				toProcess: []Meta{
					{
						PhoneNumber:    "888-888-888",
						VirtualAgentID: "aaaa-aaaa-aaa-3-ddd",
						ID:             "2",
					},
				},
				statuses: nil,
				mu:       &sync.Mutex{},
			},
			args: args{
				in0: nil,
				meta: Meta{
					PhoneNumber:    "777-777-777",
					VirtualAgentID: "sdas2-sdsada-dsad-a-sdasd",
					ID:             "3",
				},
			},
			expectedValues: expectedValues{
				toProcess: []Meta{
					{
						PhoneNumber:    "777-777-777",
						VirtualAgentID: "sdas2-sdsada-dsad-a-sdasd",
						ID:             "3",
					},
					{
						PhoneNumber:    "888-888-888",
						VirtualAgentID: "aaaa-aaaa-aaa-3-ddd",
						ID:             "2",
					},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				toProcess: tt.fields.toProcess,
				statuses:  tt.fields.statuses,
				mu:        tt.fields.mu,
			}
			ao := assert.New(t)
			actualErr := s.AddToQueueFront(tt.args.in0, tt.args.meta)
			ao.Equal(tt.expectedValues.err, actualErr)
			ao.Equal(tt.expectedValues.toProcess, s.toProcess)

		})
	}
}

func TestStorage_Next(t *testing.T) {
	type fields struct {
		toProcess []Meta
		statuses  map[ID]int
		mu        *sync.Mutex
	}
	type args struct {
		in0 context.Context
	}
	type expectedValues struct {
		toProcess []Meta
		value     Meta
		exists    bool
		err       error
	}
	var tests = []struct {
		name   string
		fields fields
		args   args
		expectedValues
	}{
		{
			name: "empty",
			fields: fields{
				toProcess: make([]Meta, 0),
				mu:        &sync.Mutex{},
			},
			args: args{
				nil,
			},
			expectedValues: expectedValues{
				toProcess: make([]Meta, 0),
				value:     Meta{},
				exists:    false,
				err:       nil,
			},
		},
		{
			name: "success",
			fields: fields{
				toProcess: []Meta{
					{
						PhoneNumber:    "888-888-888",
						VirtualAgentID: "aaaa-aaaa-aaa-3-ddd",
						ID:             "2",
					},
					{
						PhoneNumber:    "777-777-777",
						VirtualAgentID: "sdas2-sdsada-dsad-a-sdasd",
						ID:             "3",
					},
				},
				mu: &sync.Mutex{},
			},
			args: args{
				nil,
			},
			expectedValues: expectedValues{
				toProcess: []Meta{
					{
						PhoneNumber:    "888-888-888",
						VirtualAgentID: "aaaa-aaaa-aaa-3-ddd",
						ID:             "2",
					},
					{
						PhoneNumber:    "777-777-777",
						VirtualAgentID: "sdas2-sdsada-dsad-a-sdasd",
						ID:             "3",
					},
				},
				value: Meta{
					PhoneNumber:    "888-888-888",
					VirtualAgentID: "aaaa-aaaa-aaa-3-ddd",
					ID:             "2",
				},
				exists: true,
				err:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				toProcess: tt.fields.toProcess,
				statuses:  tt.fields.statuses,
				mu:        tt.fields.mu,
			}
			ao := assert.New(t)
			actualMeta, actualExists, actualErr := s.Next(tt.args.in0)
			ao.Equal(tt.expectedValues.value, actualMeta)
			ao.Equal(tt.expectedValues.exists, actualExists)
			ao.Equal(tt.expectedValues.err, actualErr)
		})
	}
}

func TestStorage_QueueLength(t *testing.T) {
	type fields struct {
		toProcess []Meta
		statuses  map[ID]int
		mu        *sync.Mutex
	}
	type args struct {
		in0 context.Context
	}
	type expectedValues struct {
		value int
		err   error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expectedValues
	}{
		{
			name: "success",
			fields: fields{
				toProcess: make([]Meta, 10),
				statuses:  nil,
				mu:        &sync.Mutex{},
			},
			args: args{nil},
			expectedValues: expectedValues{
				err:   nil,
				value: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				toProcess: tt.fields.toProcess,
				statuses:  tt.fields.statuses,
				mu:        tt.fields.mu,
			}
			ao := assert.New(t)
			actualValue, actualErr := s.QueueLength(tt.args.in0)
			ao.Equal(tt.expectedValues.value, actualValue)
			ao.Equal(tt.expectedValues.err, actualErr)
		})
	}
}

func TestStorage_SaveStatus(t *testing.T) {
	type fields struct {
		toProcess []Meta
		statuses  map[ID]int
		mu        *sync.Mutex
	}
	type args struct {
		in0    context.Context
		status int
		meta   Meta
	}
	type expectedValues struct {
		statuses map[ID]int
		err      error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expectedValues
	}{
		{
			name: "simple success",
			fields: fields{
				statuses: make(map[ID]int),
				mu:       &sync.Mutex{},
			},
			args: args{
				in0:    nil,
				status: 200,
				meta: Meta{
					PhoneNumber:    "777-777-77",
					VirtualAgentID: "aaaa-bbbb-cccc-dddd",
					ID:             "1",
				},
			},
			expectedValues: expectedValues{
				statuses: map[ID]int{
					"1": 200,
				},
				err: nil,
			},
		},
		{
			name: "success, rewrite value",
			fields: fields{
				statuses: map[ID]int{
					"1": 200,
				},
				mu: &sync.Mutex{},
			},
			args: args{
				in0:    nil,
				status: 499,
				meta: Meta{
					PhoneNumber:    "777-777-77",
					VirtualAgentID: "aaaa-bbbb-cccc-dddd",
					ID:             "1",
				},
			},
			expectedValues: expectedValues{
				statuses: map[ID]int{
					"1": 499,
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				toProcess: tt.fields.toProcess,
				statuses:  tt.fields.statuses,
				mu:        tt.fields.mu,
			}
			ao := assert.New(t)
			actualErr := s.SaveStatus(tt.args.in0, tt.args.status, tt.args.meta)
			ao.Equal(tt.expectedValues.err, actualErr)
			ao.Equal(tt.expectedValues.statuses, s.statuses)
		})
	}
}
