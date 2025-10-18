package helpers

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// very niche helper
func ptrInt(v int32) *int32 {
	return &v
}

// Convert protobuf Timestamp to Go time.Time
func timestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
