package models

import (
	"time"

	db "github.com/geo-afk/Online-Doctor-Appointment/app/postgres"
)

func (c ContactDetail) ToBdType() (contactParam db.CreateContactParams) {

	contactParam = db.CreateContactParams{
		PrimaryNumber:   ToPgText(c.PrimaryNumber),
		SecondaryNumber: ToPgText(c.SecondaryNumber),
		Email:           ToPgText(c.Email),
	}
	return
}

func (u User) ToBdType(primaryContactId, emergencyContactId int32) db.RegisterUserParams {

	var userParam = db.RegisterUserParams{
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Location:         ToPgText(u.Location),
		UserRole:         db.URole(u.UserRole),
		ContactID:        ToPgInt(primaryContactId),
		EmergencyContact: ToPgInt(emergencyContactId),
		CreatedAt:        ToPgTime(u.CreatedAt),
	}

	return userParam

}

func (a Auth) ToBdType(userId int32) (authParam db.CreateUserAuthParams) {

	authParam = db.CreateUserAuthParams{
		UserID:   userId,
		UserName: a.UserName,
		Password: a.Password,
	}

	return
}

func GetIsUserLoggedInParams(userId int32, userRole string) db.IsUserLoggedInParams {

	return db.IsUserLoggedInParams{
		UserID:   userId,
		UserRole: db.URole(userRole),
	}
}

func (a Appointment) GetBookAppointmentParams() db.BookAppointmentParams {

	status := db.NullApStatus{
		ApStatus: db.ApStatus(a.AppointmentStatus),
	}
	return db.BookAppointmentParams{
		UserID:            a.UserID,
		Reason:            a.Reason,
		BookedAt:          ToPgTime(a.BookedAt),
		AppointmentStatus: status,
		CreatedAt:         ToPgTime(a.CreatedAt),
	}
}

func CreateSession(s Session) db.InsertSessionParams {

	return db.InsertSessionParams{
		ID:           s.Id,
		UserID:       s.UserId,
		UserRole:     db.URole(s.UserRole),
		RefreshToken: s.RefreshToken,
		IsRevoked:    s.IsRevoked,
		CreatedAt:    ToPgTime(s.CreatedAt),
		ExpiresAt:    ToPgTime(s.ExpiresAt),
	}
}

func CreateCreateRequestParams(requestType, token, email string, expiresAt time.Time) db.CreateRequestParams {

	return db.CreateRequestParams{
		RequestType: requestType,
		Token:       token,
		ExpiresAt:   ToPgTime(expiresAt),
		UserEmail:   email,
	}
}

func CreateForgetPasswordParams(password, token string) db.ForgetPasswordParams {

	return db.ForgetPasswordParams{
		Password: password,
		Token:    token,
	}
}

func CreateChangePasswordParams(password string, user_id int32) db.ChangePasswordParams {

	return db.ChangePasswordParams{
		Password: password,
		ID:       user_id,
	}
}
