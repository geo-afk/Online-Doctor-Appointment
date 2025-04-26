package models

import (
	"time"
)

type model interface {
	ToBdType()
}

//
// type Doctor struct {
// 	ID                int32
// 	Specialty         string
// 	Experience        int3`2
// 	Note              string
// 	AppointmentStatus NullApStatus
// 	IsAvailable       boolean
// 	AppointmentCount  pgtype.Int4
// 	CreatedAt         pgtype.Timestamp
// }

// type MedicalHistory struct {
// 	ID            int32
// 	PatientID     pgtype.Int4
// 	Description   string
// 	DiagnosedDate pgtype.Timestamp
// 	Status        NullMedicalStatus
// 	Medication    pgtype.Int4
// 	Allergies     pgtype.Int4
// 	Surgeries     pgtype.Int4
// }

type Doctor struct {
	ID                int32
	Specialty         string
	Experience        int32
	Note              string
	AppointmentStatus string
	IsAvailable       bool
	AppointmentCount  int
	CreatedAt         time.Time
}

// type MedicalHistory struct {
// 	ID            int32
// 	PatientID     pgtype.Int4
// 	Description   string
// 	DiagnosedDate pgtype.Timestamp
// 	Status        NullMedicalStatus
// 	Medication    pgtype.Int4
// 	Allergies     pgtype.Int4
// 	Surgeries     pgtype.Int4
// }

// type MedicationHistory struct {
// 	ID          int32
// 	Description string
// 	Dosage      string
// 	Frequency   pgtype.Int4
// 	StartDate   pgtype.Timestamp
// 	EndDate     pgtype.Timestamp
// }

type Patient struct {
	ID               int32
	AppointmentID    int
	CreatedAt        time.Time
	MedicalHistoryID int
}

// type Payment struct {
// 	ID        int32
// 	PatientID int32
// 	DoctorID  int32
// 	Amount    float32
// }
//
// type Surgery struct {
// 	ID          int32
// 	Description string
// 	DateDone    pgtype.Timestamp
// 	Hospital    string
// }

type ContactDetail struct {
	Id              int32  `json:"id,omitempty"`
	PrimaryNumber   string `json:"primary_contact"`
	SecondaryNumber string `json:"secondary_number,omitempty"`
	Email           string `json:"email,omitempty"`
}
type Auth struct {
	Id       string `json:"id,omitempty"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type User struct {
	ID               int32          `json:"id,omitempty"`
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	Location         string         `json:"location"`
	UserRole         string         `json:"user_role"`
	Contact          *ContactDetail `json:"contact,omitempty"`
	EmergencyContact *ContactDetail `json:"emergency_contact,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	Auth             `json:"auth"`
}
