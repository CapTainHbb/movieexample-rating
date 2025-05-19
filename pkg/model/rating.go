package model

import "gorm.io/gorm"

// RecordID defines a record id. Together with RecordType
// identifies unique records across all types.
type RecordID string

// RecordType defines a record type. Together with RecordID
// identifies unique records across all types.
type RecordType string

const RecordTypeMovie = RecordType("movie")

// UserID defines a user id.
type UserID string

// RatingValue defines a value of a rating record.
type RatingValue int

type Rating struct {
	gorm.Model
	RecordID   string      `json:"recordId"`
	RecordType string      `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}

type RatingEvent struct {
	gorm.Model
	UserID     UserID          `json:"userId"`
	RecordID   RecordID        `json:"recordId"`
	RecordType string          `json:"recordType"`
	Value      RatingValue     `json:"value"`
	EventType  RatingEventType `json:"eventType"`
}

type RatingEventType string

const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)
