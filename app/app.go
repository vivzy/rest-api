package app

import (
  "strconv"
  "github.com/gin-gonic/gin"
_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"rest-api/app/models"
)



func AuthUsers(userId string, password string, c *gin.Context) (string, bool) {
	log.Printf("fetching user by username " + userId + " pass: " + password)
	user, ret := models.FetchUserFromUsername(userId)
	if ret == nil && user.Id > 0 {

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if err == nil {
			return user.Username, true

		} else {
			models.CheckErr(err, "user not found")
			return "", false
		}
	}else{
		models.CheckErr(ret, "user not found")
		return "", false

	}
}

func GetUsers(c *gin.Context) {
  users, err := models.FetchAllUsers()
  if err == nil {
    c.JSON(200, users)
  } else {
    c.JSON(404, gin.H{"error": "no user(s) into the table"})
  }
}


func GetUser(c *gin.Context) {
  id := c.Params.ByName("id")
  user, err := models.FetchUserFromId(id)
  if err == nil {
    user_id, _ := strconv.ParseInt(id, 0, 64)
  
    content := &models.User{
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
}

func PostUser(c *gin.Context) {
  var user models.User
  c.BindJSON(&user)
  u, ret := models.FetchUserFromUsername(user.Username)
  if ret == nil {
  	models.LogErr(ret, "username already exists : " + u.Username)
  	c.JSON(500, gin.H{"error": "username already exists"})
  }else {

	  hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	  //gin.Logger()
	  if err != nil {
		  models.CheckErr(err, "Hashing failed")
	  }
	  user.Password = string(hashedPass)

	  if user.Username != "" && user.Password != "" {
		  insertedUser, err := models.InsertUser(user)
			  if err == nil {
				  content := insertedUser
				  c.JSON(201, content)
			  } else {
				  models.CheckErr(err, "Insert failed")
			  }
	  }else {
		  c.JSON(422, gin.H{"error": "fields are empty"})
	  }
  }
}

func UpdateUser(c *gin.Context) {
  id := c.Params.ByName("id")
  user, err := models.FetchUserFromId(id)
  
  if err == nil {
    var json models.User
    c.BindJSON(&json)
    updateUser := models.User{
      Id: user.Id,
      Firstname: json.Firstname,
      Lastname: json.Lastname,
      Username: user.Username,
      Password: user.Password,
      Joindate: user.Joindate,
    }

    if updateUser.Firstname != "" && updateUser.Lastname != ""{
    _, err = models.UpdateUser(updateUser)

      if err == nil {
        c.JSON(200, updateUser)
      } else {
        models.CheckErr(err, "Updated failed")
      }
    } else {
      c.JSON(422, gin.H{"error": "fields are empty"})
    }
  } else {
    c.JSON(404, gin.H{"error": "user not found"})
  }
}

func DeleteUser(c *gin.Context) {
  id := c.Params.ByName("id")
  user, err := models.FetchUserFromId(id)
  
  if err == nil {
    _, err = models.DeleteUser(user)
    
    if err == nil {
      c.JSON(200, gin.H{"id #" + id: " deleted"})
    } else {
      models.CheckErr(err, "Delete failed")
    }
  } else {
    c.JSON(404, gin.H{"error": "user not found"})
  }
}
