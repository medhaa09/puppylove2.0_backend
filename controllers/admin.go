package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pclubiitk/puppylove2.0_backend/db"
	"github.com/pclubiitk/puppylove2.0_backend/models"
	"github.com/pclubiitk/puppylove2.0_backend/utils"
)

var Db db.PuppyDb
var permit bool = true

func AdminLogin(c *gin.Context) {
	info := new(models.AdminLogin)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	if info.Id != os.Getenv("ADMIN_ID") {
		c.JSON(http.StatusForbidden, gin.H{"error": "This action will be reported."})
		return
	}

	if info.Pass != os.Getenv("ADMIN_PASS") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Password."})
		return
	}

	token, err := generateJWTToken(info.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}
	expirationTime := time.Now().Add(time.Hour * 24)
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  expirationTime,
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		HttpOnly: true,
		Secure:   false, // Set this to true if you're using HTTPS, false for HTTP
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"message": "Admin logged in successfully !!"})
}

func AddNewUser(c *gin.Context) {
	// TODO: Modify this function to handle multiple concatenated json inputs

	// TODO: Implement admin authentication logic
	// Authenticate the admin here

	// Validate the input format
	info := new(models.AddNewUser)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	// Create user
	for _, user := range info.TypeUserNew {

		newUser := models.User{
			Id:            user.Id,
			Name:          user.Name,
			Email:         user.Email,
			Gender:        user.Gender,
			Pass:          "",
			PubK:          "",
			PrivK:         "",
			AuthC:         utils.RandStringRunes(15),
			Data:          "",
			Submit:        false,
			Matches:       "",
			ReceivedSongs: "",
			Dirty:         false,
			Publish:       false,
			Code:          "",
			About:         "",
			Intrests:      "{}",
		}

		// Insert the user into the database

		if err := Db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
}

func DeleteUser(c *gin.Context) {
	// TODO: Implement admin authentication logic
	// Authenticate the admin here

	// Validate the input format
	info := new(models.TypeUserNew)
	if err := c.BindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input data format."})
		return
	}

	newUser := models.User{
		Id:     info.Id,
		Name:   info.Name,
		Email:  info.Email,
		Gender: info.Gender,
	}

	if err := Db.Unscoped().Delete(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User Deleted successfully."})
}

func DeleteAllUsers(c *gin.Context) {
	// TODO: Implement admin authentication logic
	// Authenticate the admin here

	newUser := models.User{}
	if err := Db.Unscoped().Where("1 = 1").Delete(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "All Users Deleted successfully."})
}

// func PublishResults(c *gin.Context) {
// 	if !models.PublishMatches {
// 		var matchdb models.MatchTable
// 		var matches []models.MatchTable
// 		records := Db.Model(&matchdb).Where("").Find(&matches)
// 		if records.Error != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Some error occured while calculating matches"})
// 			return
// 		}
// 		matchesMap := make(map[string][]string)
// 		for _, match := range matches {
// 			roll1 := match.Roll1
// 			roll2 := match.Roll2
// 			//song12 := match.SONG12
// 			song21 := match.SONG21
// 			var userdb models.User
// 			var userdb1 models.User
// 			record := Db.Model(&userdb).Where("id = ?", roll1).First(&userdb)
// 			if record.Error != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while connecting to db"})
// 				return
// 			}
// 			record = Db.Model(&userdb1).Where("id = ?", roll2).First(&userdb1)
// 			if record.Error != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while connecting to db"})
// 				return
// 			}
// 			// if userdb.Publish && userdb1.Publish {
// 			// 	matchesMap[roll1] = append(matchesMap[roll1], fmt.Sprintf("%s (Song: %s)", roll2, song12))
// 			// 	matchesMap[roll2] = append(matchesMap[roll2], fmt.Sprintf("%s (Song: %s)", roll1, song21))
// 			// }
// 			if userdb.Publish && userdb1.Publish {
// 				matchesMap[roll1] = append(matchesMap[roll1], roll2)
// 				matchesMap[roll2] = append(matchesMap[roll2], roll1)
// 			}
// 		}
// 		for key := range matchesMap {
// 			var userdb models.User
// 			record := Db.Model(&userdb).Where("id = ?", key).First(&userdb)
// 			if record.Error != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating matches of " + key})
// 				return
// 			}
// 			tempMap := make(map[string]bool)
// 			for _, match := range matchesMap[key] {
// 				tempMap[match] = true
// 			}
// 			results := []string{}
// 			for key := range tempMap {
// 				results = append(results, key)
// 			}
// 			userdb.Matches = strings.Join(results, ",")
// 			record = Db.Save(&userdb)
// 			if record.Error != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating matches of " + key})
// 				return
// 			}

//			}
//			models.PublishMatches = true
//			c.JSON(http.StatusOK, gin.H{"msg": "Published Matches"})
//			return
//		}
//		c.JSON(http.StatusOK, gin.H{"msg": "Matches already published"})
//	}
func PublishResults(c *gin.Context) {
	if !models.PublishMatches {
		var matchdb models.MatchTable
		var matches []models.MatchTable
		records := Db.Model(&matchdb).Find(&matches)
		if records.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Some error occurred while calculating matches"})
			return
		}

		matchesMap := make(map[string][]string)
		songsMap := make(map[string][]string) // Map to store received songs with senders

		for _, match := range matches {
			roll1 := match.Roll1
			roll2 := match.Roll2
			song21 := match.SONG21 // Song sent by roll2 to roll1

			var userdb models.User
			var userdb1 models.User

			record := Db.Model(&userdb).Where("id = ?", roll1).First(&userdb)
			if record.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while connecting to DB"})
				return
			}
			record = Db.Model(&userdb1).Where("id = ?", roll2).First(&userdb1)
			if record.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while connecting to DB"})
				return
			}

			if userdb.Publish && userdb1.Publish {
				matchesMap[roll1] = append(matchesMap[roll1], roll2)
				matchesMap[roll2] = append(matchesMap[roll2], roll1)

				// Store received song with sender info only if song21 is not empty
				if song21 != "" {
					songsMap[roll1] = append(songsMap[roll1], fmt.Sprintf("%s:%s", roll2, song21))
				}
			}
		}

		for key := range matchesMap {
			var userdb models.User
			record := Db.Model(&userdb).Where("id = ?", key).First(&userdb)
			if record.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating matches of " + key})
				return
			}

			// Remove duplicate matches
			tempMatchMap := make(map[string]bool)
			for _, match := range matchesMap[key] {
				tempMatchMap[match] = true
			}
			results := []string{}
			for match := range tempMatchMap {
				results = append(results, match)
			}
			userdb.Matches = strings.Join(results, ",")

			// Ensure that if no song is received, store an empty string
			if len(songsMap[key]) == 0 {
				userdb.ReceivedSongs = ""
			} else {
				userdb.ReceivedSongs = strings.Join(songsMap[key], ",")
			}

			// Save updated user data
			record = Db.Save(&userdb)
			if record.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating matches of " + key})
				return
			}
		}

		models.PublishMatches = true
		c.JSON(http.StatusOK, gin.H{"msg": "Published Matches"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Matches already published"})
}

func TogglePermit(c *gin.Context) {
	permit = !permit
	c.JSON(http.StatusOK, gin.H{"permitStatus": permit})
}
