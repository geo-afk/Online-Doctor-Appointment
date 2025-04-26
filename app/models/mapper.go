package models

import (
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

func GetIsUserLoggedInParams(userName, userRole string) db.IsUserLoggedInParams {

	return db.IsUserLoggedInParams{
		UserName: userName,
		UserRole: db.URole(userRole),
	}
}
