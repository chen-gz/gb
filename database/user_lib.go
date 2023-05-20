package database

// This  is first version of user database design.

// The Role determine user is admin or not.
// The Level determine user's permission.
// The group can be defined by admin. The post will apply filter based on group.

type UserData struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	Level int    `json:"level"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

func GetUser(email string) UserData {
	// get user data from database
	if email == "chen-gz@outlook.com" {
		return UserData{
			Email: "chen-gz@outlook.com",
			Role:  "admin",
			Level: 100,
			Name:  "Guangzong",
			Group: "admin",
		}
	}
	return UserData{}
}
