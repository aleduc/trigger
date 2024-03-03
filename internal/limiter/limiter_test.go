package limiter

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"test_trigger/internal/realtime"
)

func TestNewSlidingWindow(t *testing.T) {
	ctrl := gomock.NewController(t)
	realTimeMock := realtime.NewMockTime(ctrl)
	expected := &SlidingWindow{
		RealTime:    realTimeMock,
		size:        10,
		limit:       25,
		windowStart: 1709464831,
		counter:     0,
		entries:     make([]uint64, 10),
		mu:          &sync.Mutex{},
	}
	realTimeMock.EXPECT().Now().Return(time.Unix(1709464831, 0)).Times(1)
	assert.Equal(t, expected, NewSlidingWindow(10, 25, realTimeMock))

}

func TestSlidingWindow_Allow(t *testing.T) {
	type fields struct {
		size        int64
		limit       uint64
		windowStart int64
		counter     uint64
		entries     []uint64
		mu          *sync.Mutex
	}
	type expected struct {
		allow               bool
		expectedWindowStart int64
		expectedCounter     uint64
		expectedEntries     []uint64
	}
	tests := []struct {
		name         string
		fields       fields
		expectedFunc func(time *realtime.MockTime)
		expected
	}{
		{
			name: "simple success",
			fields: fields{
				size:        10,
				limit:       10,
				windowStart: 1709464830,
				counter:     0,
				entries:     make([]uint64, 10),
				mu:          &sync.Mutex{},
			},
			expectedFunc: func(t *realtime.MockTime) {
				t.EXPECT().Now().Return(time.Unix(1709464831, 0)).Times(1)
			},
			expected: expected{
				allow:               true,
				expectedWindowStart: 1709464830,
				expectedCounter:     1,
				expectedEntries:     []uint64{0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		{
			name: "simple blocked",
			fields: fields{
				size:        10,
				limit:       10,
				windowStart: 1709464830,
				counter:     10,
				entries:     make([]uint64, 10),
				mu:          &sync.Mutex{},
			},
			expectedFunc: func(t *realtime.MockTime) {
				t.EXPECT().Now().Return(time.Unix(1709464831, 0)).Times(1)
			},
			expected: expected{
				allow:               false,
				expectedWindowStart: 1709464830,
				expectedCounter:     10,
				expectedEntries:     []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		{
			name: "success with state refresh(clean first 5 elements)",
			fields: fields{
				size:        10,
				limit:       10,
				windowStart: 1709464830,
				counter:     59,
				entries:     []uint64{10, 10, 10, 10, 10, 1, 2, 3, 3, 0},
				mu:          &sync.Mutex{},
			},
			expectedFunc: func(t *realtime.MockTime) {
				t.EXPECT().Now().Return(time.Unix(1709464844, 0)).Times(1)
			},
			expected: expected{
				allow:               true,
				expectedWindowStart: 1709464835,
				expectedCounter:     10,
				expectedEntries:     []uint64{0, 0, 0, 0, 1, 1, 2, 3, 3, 0},
			},
		},
		{
			name: "success with state refresh(clean first 6 elements), but head = 3",
			fields: fields{
				size:        10,
				limit:       10,
				windowStart: 1709464833,
				counter:     28,
				entries:     []uint64{1, 1, 1, 10, 10, 1, 1, 1, 1, 1},
				mu:          &sync.Mutex{},
			},
			expectedFunc: func(t *realtime.MockTime) {
				t.EXPECT().Now().Return(time.Unix(1709464848, 0)).Times(1)
			},
			expected: expected{
				allow:               true,
				expectedWindowStart: 1709464839,
				expectedCounter:     5,
				expectedEntries:     []uint64{1, 1, 1, 0, 0, 0, 0, 0, 1, 1},
			},
		},
		{
			name: "success with state refresh 8,9,0,1,2,3 elements",
			fields: fields{
				size:        10,
				limit:       10,
				windowStart: 1709464838,
				counter:     10,
				entries:     []uint64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				mu:          &sync.Mutex{},
			},
			expectedFunc: func(t *realtime.MockTime) {
				t.EXPECT().Now().Return(time.Unix(1709464853, 0)).Times(1)
			},
			expected: expected{
				allow:               true,
				expectedWindowStart: 1709464844,
				expectedCounter:     5,
				expectedEntries:     []uint64{0, 0, 0, 1, 1, 1, 1, 1, 0, 0},
			},
		},
		{
			name: "simple success with clean instead of buffer",
			fields: fields{
				size:        10,
				limit:       10,
				windowStart: 1709464830,
				counter:     59,
				entries:     []uint64{10, 10, 10, 10, 10, 1, 2, 3, 3, 0},
				mu:          &sync.Mutex{},
			},
			expectedFunc: func(t *realtime.MockTime) {
				t.EXPECT().Now().Return(time.Unix(1709464854, 0)).Times(1)
			},
			expected: expected{
				allow:               true,
				expectedWindowStart: 1709464854,
				expectedCounter:     1,
				expectedEntries:     []uint64{0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
			},
		},
		{
			name: "simple success with state refresh(clean first 5 elements)",
			fields: fields{
				size:        10,
				limit:       10,
				windowStart: 1709464830,
				counter:     59,
				entries:     []uint64{10, 10, 10, 10, 10, 1, 2, 3, 3, 0},
				mu:          &sync.Mutex{},
			},
			expectedFunc: func(t *realtime.MockTime) {
				t.EXPECT().Now().Return(time.Unix(1709464844, 0)).Times(1)
			},
			expected: expected{
				allow:               true,
				expectedWindowStart: 1709464835,
				expectedCounter:     10,
				expectedEntries:     []uint64{0, 0, 0, 0, 1, 1, 2, 3, 3, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockTime := realtime.NewMockTime(ctrl)
			s := &SlidingWindow{
				RealTime:    mockTime,
				size:        tt.fields.size,
				limit:       tt.fields.limit,
				windowStart: tt.fields.windowStart,
				counter:     tt.fields.counter,
				entries:     tt.fields.entries,
				mu:          tt.fields.mu,
			}
			ao := assert.New(t)
			if tt.expectedFunc != nil {
				tt.expectedFunc(mockTime)
			}

			actual := s.Allow()
			ao.Equal(tt.expected.allow, actual)
			ao.Equal(tt.expectedCounter, s.counter)
			ao.Equal(tt.expectedEntries, s.entries)
			ao.Equal(tt.expectedWindowStart, s.windowStart)
		})
	}
}

// TODO add tests.
func TestSlidingWindow_refreshCurrentState(t *testing.T) {
}
