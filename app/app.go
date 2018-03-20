package app

import (
  "strconv"
  "github.com/gin-gonic/gin"
  _ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
 Id int64 `db:"id" json:"id"`
 Username string `db:"username" json:"username"`
 Firstname string `db:"firstname" json:"firstname"`
 Lastname string `db:"lastname" json:"lastname"`
 Password string `db:"password" json:"password"`
 Joindate string `db:"joindate" json:"joindate"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var dbmap = initDb()

func AuthUsers(userId string, password string, c *gin.Context) (string, bool) {
	log.Printf("fetching user by username " + userId + " pass: " + password)
	var user User
	derr := dbmap.SelectOne(&user, "SELECT * FROM user WHERE username=?  LIMIT 1", userId)
	if derr != nil{
		checkErr(derr, "username not found")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password) )

	if err == nil {
		return user.Username, true

	} else {
		checkErr(err, "user not found")
		return "", false
	}
}

func GetUsers(c *gin.Context) {
  var users []User
  _, err := dbmap.Select(&users, "SELECT * FROM user")
  if err == nil {
    c.JSON(200, users)
  } else {
    c.JSON(404, gin.H{"error": "no user(s) into the table"})
  }
  // curl -i http://localhost:8080/api/v1/users
}


func GetUser(c *gin.Context) {
  id := c.Params.ByName("id")
  var user User
  
  err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)
  if err == nil {
    user_id, _ := strconv.ParseInt(id, 0, 64)
  
    content := &User{
      Id: user_id,
      Firstname: user.Firstname,
      Lastname: user.Lastname,
      Username: user.Username,
      Joindate: user.Joindate,
    }
 
    c.JSON(200, content)
  } else {
    c.JSON(404, gin.H{"error": "user not found"})
  }
  // curl -i http://localhost:8080/api/v1/Users/1
}

func PostUser(c *gin.Context) {
  var user User
  c.BindJSON(&user)
  hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
  //gin.Logger()
	if err != nil {
		checkErr(err, "Hashing failed")
	}
	user.Password = string(hashedPass)

  if user.Username != "" && user.Password != "" {
    if insert, _ := dbmap.Exec(
    	`INSERT INTO user (username, firstname, lastname, joindate, password) VALUES (?, ?, ?, ?, ?)`,
    	user.Username, user.Firstname, user.Lastname, user.Joindate, user.Password); insert != nil {
      user_id, err := insert.LastInsertId()
      if err == nil {
        content := &User{
          Id: user_id,
          Username: user.Username,
          Firstname: user.Firstname,
          Lastname: user.Lastname,
          Joindate: user.Joindate,
        }
        c.JSON(201, content)
      } else {
        checkErr(err, "Insert failed")
      }
    }
  } else {
    c.JSON(422, gin.H{"error": "fields are empty"})
  }
  // curl -i -X POST -H "Content-Type: application/json" -d "{ \"event_status\": \"83\", \"event_name\": \"100\" }" http://localhost:8080/api/v1/users
}

func UpdateUser(c *gin.Context) {
  id := c.Params.ByName("id")
  var user User
  err := dbmap.SelectOne(&user, "SELECT * FROM user WHERE id=?", id)
  
  if err == nil {
    var json User
    c.BindJSON(&json)
    updateUser := User{
      Id: user.Id,
      Firstname: json.Firstname,
      Lastname: json.Lastname,
      Username: user.Username,
      Password: user.Password,
      Joindate: user.Joindate,
    }

    if updateUser.Firstname != "" && updateUser.Lastname != ""{
    _, err = dbmap.Update(&updateUser)

      if err == nil {
        c.JSON(200, updateUser)
      } else {
        checkErr(err, "Updated failed")
      }
    } else {
      c.JSON(422, gin.H{"error": "fields are empty"})
    }
  } else {
    c.JSON(404, gin.H{"error": "user not found"})
  }
  // curl -i -X PUT -H "Content-Type: application/json" -d "{ \"event_status\": \"83\", \"event_name\": \"100\" }" http://localhost:8080/api/v1/users/1
}

func DeleteUser(c *gin.Context) {
  id := c.Params.ByName("id")
  var user User
  err := dbmap.SelectOne(&user, "SELECT id FROM User WHERE id=?", id)
  
  if err == nil {
    _, err = dbmap.Delete(&user)
    
    if err == nil {
      c.JSON(200, gin.H{"id #" + id: " deleted"})
    } else {
      checkErr(err, "Delete failed")
    }
  } else {
    c.JSON(404, gin.H{"error": "user not found"})
  }
  // curl -i -X DELETE http://localhost:8080/api/v1/users/1
}
