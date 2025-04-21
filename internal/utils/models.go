package utils

import "time"

type ApStatus string
type URole string

// type ContactDetail struct {
// 	ID              int32
// 	PrimaryNumber   pgtype.Text
// 	SecondaryNumber pgtype.Text
// 	Email           pgtype.Text
// }
//
// type Doctor struct {
// 	ID                int32
// 	Specialty         string
// 	Experience        int3`2
// 	Note              pgtype.Text
// 	AppointmentStatus NullApStatus
// 	IsAvailable       pgtype.Bool
// 	AppointmentCount  pgtype.Int4
// 	CreatedAt         pgtype.Timestamp
// }

// type MedicalHistory struct {
// 	ID            int32
// 	PatientID     pgtype.Int4
// 	Description   pgtype.Text
// 	DiagnosedDate pgtype.Timestamp
// 	Status        NullMedicalStatus
// 	Medication    pgtype.Int4
// 	Allergies     pgtype.Int4
// 	Surgeries     pgtype.Int4
// }

type Auth struct {
	UserID   int32  `json:"user_id,omitempty"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type Doctor struct {
	ID                int32
	Specialty         string
	Experience        int32
	Note              string
	AppointmentStatus ApStatus
	IsAvailable       bool
	AppointmentCount  int
	CreatedAt         time.Time
}

// type MedicalHistory struct {
// 	ID            int32
// 	PatientID     pgtype.Int4
// 	Description   pgtype.Text
// 	DiagnosedDate pgtype.Timestamp
// 	Status        NullMedicalStatus
// 	Medication    pgtype.Int4
// 	Allergies     pgtype.Int4
// 	Surgeries     pgtype.Int4
// }
// type MedicationHistory struct {
// 	ID          int32
// 	Description pgtype.Text
// 	Dosage      pgtype.Text
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
// 	Description pgtype.Text
// 	DateDone    pgtype.Timestamp
// 	Hospital    pgtype.Text
// }

type User struct {
	ID               int32
	FirstName        string
	LastName         string
	Location         string
	UserRole         URole
	ContactID        int
	EmergencyContact int
	CreatedAt        time.Time
}
