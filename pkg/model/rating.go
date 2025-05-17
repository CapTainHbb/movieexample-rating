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
