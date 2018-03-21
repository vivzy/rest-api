package models

type User struct {
	Id int64 `db:"id" json:"id"`
	Username string `db:"username" json:"username" binding:"required"`
	Firstname string `db:"firstname" json:"firstname" binding:"required"`
	Lastname string `db:"lastname" json:"lastname"`
	Password string `db:"password" json:"password" binding:"required"`
	Joindate string `db:"joindate" json:"joindate"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var dbmap = initDb()

func FetchUserFromUsername(username string) (u User, ret error) {
	var user User
	derr := dbmap.SelectOne(&user, "SELECT * FROM user WHERE username=?  LIMIT 1", username)
	if derr != nil{
		LogErr(derr, "username not found")
		return User{}, derr
	}else{
		return user, nil
	}
}

func FetchUserFromId(userId string) (u User, ret error){
	var user User
	derr := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id = ?  ", userId)
	if derr != nil{
		LogErr(derr, "user not found")
		return User{}, derr
	}else{
		return user, nil
	}
}

func FetchAllUsers() (u []User, ret error){
	var users []User
	_, err := dbmap.Select(&users, "SELECT username, firstname, lastname, joindate FROM user")
	return users, err
}

func InsertUser(user User) (finalUser User, ret error){
	if insert, err := dbmap.Exec(
		`INSERT INTO user (username, firstname, lastname, joindate, password) VALUES (?, ?, ?, ?, ?)`,
		user.Username, user.Firstname, user.Lastname, user.Joindate, user.Password); insert != nil {
		user_id, err := insert.LastInsertId()
		if err == nil {
			content := User{
				Id:        user_id,
				Username:  user.Username,
				Firstname: user.Firstname,
				Lastname:  user.Lastname,
				Joindate:  user.Joindate,
			}
			return content, nil
		}else{
			return User{}, err
		}
	}else{
		return User{}, err
	}
}

func UpdateUser(user User) (finalUser User, ret error){
	_,err := dbmap.Update(&user)
	return user, err
}

func DeleteUser(user User) (delUser User, ret error){
	_,err := dbmap.Delete(&user)
	return user, err
}