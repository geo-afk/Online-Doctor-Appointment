package utils

import db "github.com/geo-afk/Online-Doctor-Appointment/internal/postgres"

func AuthToUserLoginParams(auth Auth) (userLogin db.UserLoginParams) {

	userLogin = db.UserLoginParams{
		UserName: auth.UserName,
		Password: GeneratePassword(auth.Password),
	}

	return

}
