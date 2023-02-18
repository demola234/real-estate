package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/demola234/real-estate/database"
	"github.com/demola234/real-estate/helpers"
	"github.com/demola234/real-estate/interfaces"
	"github.com/demola234/real-estate/models"
	"github.com/demola234/real-estate/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var users []models.User
		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var user models.User
			cursor.Decode(&user)
			users = append(users, user)
		}
		defer cancel()
		c.JSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("user_id")

		var foundUser models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundUser)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
		}
		if *foundUser.Otp_Verified {
			c.JSON(http.StatusOK, gin.H{
				"_id":           foundUser.ID,
				"token":         foundUser.Token,
				"refresh_token": foundUser.Refresh_Token,
				"first_name":    foundUser.First_Name,
				"last_name":     foundUser.Last_Name,
				"email":         foundUser.Email,
				"otp":           foundUser.Otp,
				"avatar":        foundUser.Avatar,
				"phone":         foundUser.Phone,
				"created_at":    foundUser.Created_at,
				"updated_at":    foundUser.Updated_at,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Please check your email to get your verify otp and verify your account"})

		}
	}

}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 200*time.Second)
		defer cancel()
		var users models.User

		phone := c.PostForm("phone")
		email := c.PostForm("email")
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")

		users.First_Name = &first_name
		users.Last_Name = &last_name
		// check if user already exists
		countEmail, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if countEmail > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user with this email already exists"})
			return
		}
		users.Email = &email

		countPhone, err := userCollection.CountDocuments(ctx, bson.M{"phone": phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if countPhone > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user with this phone number already exists"})
			return
		}
		users.Phone = &phone

		avatar, _, err := c.Request.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error uploading profile picture"})
			return
		}

		uploadImage, uploadErr := utils.UploadAvatar(ctx, avatar)
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading profile picture"})
			return
		}

		users.Avatar = &uploadImage

		otp := utils.GenerateOTP()
		otpExpiration := time.Now().Add(15 * time.Minute) // set expiration time to 15 minutes from now
		otpValid := false
		users.Otp = &otp
		users.Otp_Verified = &otpValid
		users.Otp_Expiration = otpExpiration

		sendErr := utils.SendOTP(email, otp, first_name)

		if sendErr != nil {
			c.JSON(http.StatusInternalServerError, sendErr.Error())
			return
		}

		users.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		users.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		users.ID = primitive.NewObjectID()
		users.User_id = users.ID.Hex()

		token, refreshToken, _ := helpers.GenerateAllTokens(*users.Email, *users.First_Name, *users.Last_Name, users.User_id)
		users.Token = &token
		users.Refresh_Token = &refreshToken

		insertedID, insertErr := userCollection.InsertOne(ctx, users)
		if insertErr != nil {
			msg := fmt.Sprintf("error inserting user: %v", insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": "please check your email to verify your account", "id": insertedID})
	}
}

func CreatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 200*time.Second)
		defer cancel()
		var password interfaces.Password
		var foundUser models.User

		if err := c.BindJSON(&password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		error := userCollection.FindOne(ctx, bson.M{"email": password.Email}).Decode(&foundUser)
		if error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		newPass := password.Password
		if foundUser.Password != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password already exists"})

		}

		passwordUser, err := utils.HashPassword(newPass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
			return
		}

		userCollection.UpdateOne(ctx, bson.M{"email": password.Email}, bson.M{"$set": bson.M{"password": passwordUser}})
		c.JSON(http.StatusOK, gin.H{"data": "password created successfully"})

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		//convert the login data from postman which is in JSON to golang readable format

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}

		//find a user with that email and see if that user even exists
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found, login seems to be incorrect"})
			return
		}
		//check if the user has verified their account
		if !*foundUser.Otp_Verified {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Please check your email to get your verify otp and verify your account"})
			return
		}

		//then you will verify the password
		passwordIsValid, msg := utils.VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		//if all goes well, then you'll generate tokens

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_id)

		// save push token
		if user.Push_Notification_Token != nil {
			_, err := userCollection.UpdateOne(ctx, bson.M{"user_id": foundUser.User_id}, bson.M{"$set": bson.M{"push_notification_token": user.Push_Notification_Token}})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		//update tokens - token and refresh token
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		//return statusOK
		c.JSON(http.StatusOK, gin.H{
			"_id":           foundUser.ID,
			"token":         foundUser.Token,
			"refresh_token": foundUser.Refresh_Token,
			"first_name":    foundUser.First_Name,
			"last_name":     foundUser.Last_Name,
			"email":         foundUser.Email,
			"avatar":        foundUser.Avatar,
			"phone":         foundUser.Phone,
			"created_at":    foundUser.Created_at,
			"updated_at":    foundUser.Updated_at,
		})
	}
}

func VerifyOtp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var verifyOtp interfaces.VerifyOtp
		var foundUser models.User
		if err := c.ShouldBindJSON(&verifyOtp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		email := verifyOtp.Email
		otp := verifyOtp.Otp

		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if *foundUser.Otp_Verified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you have already verified your account"})
			return
		}

		if *foundUser.Otp != otp {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wrong otp"})
			return
		}

		if time.Now().After(foundUser.Otp_Expiration) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "otp expired"})
			return
		}

		if *foundUser.Otp == otp {
			*foundUser.Otp_Verified = true
		}

		userCollection.UpdateByID(ctx, foundUser.ID, bson.M{"$set": bson.M{"otp_verified": true}})

		c.JSON(http.StatusOK, gin.H{"data": "otp verified successfully"})

	}
}

func ResendOtp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var resendOtp interfaces.ResendOtp
		var foundUser models.User
		if err := c.ShouldBindJSON(&resendOtp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to resend otp, please try again later"})
			return
		}

		email := resendOtp.Email

		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to resend otp, please try again later"})
			return
		}

		if *foundUser.Otp_Verified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you have already verified your account"})
			return
		}

		newOtp := utils.GenerateOTP()
		otpExpiration := time.Now().Add(15 * time.Minute) // set expiration time to 15 minutes from now

		sendErr := utils.SendOTP(email, newOtp, *foundUser.First_Name)

		if sendErr != nil {
			c.JSON(http.StatusInternalServerError, sendErr.Error())
			return
		}

		updated, userErr := userCollection.UpdateByID(ctx, foundUser.ID, bson.M{"$set": bson.M{"otp": newOtp, "otp_expiration": otpExpiration}})
		if userErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": userErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": "otp resent successfully", "updated": updated})
	}
}

func ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var forgetPassword interfaces.ForgetPassword
		var foundUser models.User
		if err := c.ShouldBindJSON(&forgetPassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		email := forgetPassword.Email

		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !*foundUser.Otp_Verified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you have not verified your account"})
			return
		}

		newOtp := utils.GenerateOTP()
		otpExpiration := time.Now().Add(15 * time.Minute) // set expiration time to 15 minutes from now

		sendErr := utils.SendOTP(email, newOtp, *foundUser.First_Name)

		if sendErr != nil {
			c.JSON(http.StatusInternalServerError, sendErr.Error())
			return
		}

		updated, userErr := userCollection.UpdateByID(ctx, foundUser.ID, bson.M{"$set": bson.M{"otp": newOtp, "otp_expiration": otpExpiration}})
		if userErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": userErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": "otp resent successfully", "updated": updated})

	}
}

func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var resetPassword interfaces.ResetPassword
		var foundUser models.User

		if err := c.ShouldBindJSON(&resetPassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		email := resetPassword.Email

		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !*foundUser.Otp_Verified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you have not verified your account"})
			return
		}

		if *foundUser.Otp != resetPassword.Otp {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid otp"})
			return
		}

		if time.Now().After(foundUser.Otp_Expiration) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "otp expired"})
			return
		}

		passwordIsValid, msg := utils.VerifyPassword(resetPassword.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		//return statusOK
		c.JSON(http.StatusOK, gin.H{"data": "password reset successfully"})

	}
}

func GoogleOauth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var googleAuth interfaces.GoogleAuth
		var user models.User
		var foundUser models.User
		var body interfaces.GoogleOAuthPayload

		if err := c.ShouldBindJSON(&googleAuth); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authorizationCode := googleAuth.Code
		pushToken := googleAuth.Token
		client := http.Client{}

		var url = "https://www.googleapis.com/oauth2/v3/userinfo"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authorizationCode)

		res, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		requestBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		json.Unmarshal(requestBody, &body)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		countEmail, err := userCollection.CountDocuments(ctx, bson.M{"email": body.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if countEmail <= 0 {
			user.Avatar = &body.Picture
			user.First_Name = &body.GivenName
			user.Last_Name = &body.FamilyName
			user.Email = &body.Email
			otpVerified := true
			user.Otp_Verified = &otpVerified
			user.Otp = &body.Sub

			user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.ID = primitive.NewObjectID()
			user.User_id = user.ID.Hex()

			token, refreshToken, _ := helpers.GenerateAllTokens(body.Email, *user.First_Name, *user.Last_Name, user.User_id)
			user.Token = &token
			user.Refresh_Token = &refreshToken

			_, err := userCollection.UpdateOne(ctx, bson.M{"user_id": foundUser.User_id}, bson.M{"$set": bson.M{"push_notification_token": pushToken}})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			_, insertErr := userCollection.InsertOne(ctx, user)
			if insertErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"_id":           user.ID,
				"token":         user.Token,
				"refresh_token": user.Refresh_Token,
				"first_name":    user.First_Name,
				"last_name":     user.Last_Name,
				"email":         user.Email,
				"avatar":        user.Avatar,
				"phone":         user.Phone,
				"created_at":    user.Created_at,
				"updated_at":    user.Updated_at,
			})
		} else {
			//find a user with that email and see if that user even exists
			err := userCollection.FindOne(ctx, bson.M{"email": body.Email}).Decode(&foundUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found, login seems to be incorrect"})
				return
			}

			//if all goes well, then you'll generate tokens

			token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_id)

			// save push token
			if user.Push_Notification_Token != nil {
				_, err := userCollection.UpdateOne(ctx, bson.M{"user_id": foundUser.User_id}, bson.M{"$set": bson.M{"push_notification_token": user.Push_Notification_Token}})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			//update tokens - token and refresh token
			helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

			//return statusOK
			c.JSON(http.StatusOK, gin.H{
				"_id":           foundUser.ID,
				"token":         foundUser.Token,
				"refresh_token": foundUser.Refresh_Token,
				"first_name":    foundUser.First_Name,
				"last_name":     foundUser.Last_Name,
				"email":         foundUser.Email,
				"avatar":        foundUser.Avatar,
				"phone":         foundUser.Phone,
				"created_at":    foundUser.Created_at,
				"updated_at":    foundUser.Updated_at,
			})

		}
	}
}

func FacebookOauth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var facebookAuth interfaces.FacebookAuth
		var user models.User
		var foundUser models.User
		var body interfaces.FaceResponse

		if err := c.ShouldBindJSON(&facebookAuth); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authorizationCode := facebookAuth.Code
		pushToken := facebookAuth.Token
		userId := facebookAuth.ID
		client := http.Client{}

		var url = fmt.Sprintf("https://graph.facebook.com/%v?fields=id,name,email,picture&access_token=%v", userId, authorizationCode)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		requestBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		json.Unmarshal(requestBody, &body)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		countEmail, err := userCollection.CountDocuments(ctx, bson.M{"email": body.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if countEmail <= 0 {
			user.Avatar = &body.Picture.Data.URL
			user.First_Name = &strings.Split(body.Name, " ")[0]
			user.Last_Name = &strings.Split(body.Name, " ")[1]
			user.Email = &body.Email

			otpVerified := true
			user.Otp_Verified = &otpVerified
			user.Push_Notification_Token = &pushToken
			user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.ID = primitive.NewObjectID()
			user.User_id = user.ID.Hex()

			token, refreshToken, _ := helpers.GenerateAllTokens(body.Email, *user.First_Name, *user.Last_Name, user.User_id)

			user.Token = &token
			user.Refresh_Token = &refreshToken

			_, err := userCollection.UpdateOne(ctx, bson.M{"user_id": foundUser.User_id}, bson.M{"$set": bson.M{"push_notification_token": pushToken}})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			_, insertErr := userCollection.InsertOne(ctx, user)
			if insertErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"_id":           user.ID,
				"token":         user.Token,
				"refresh_token": user.Refresh_Token,
				"first_name":    user.First_Name,
				"last_name":     user.Last_Name,
				"email":         user.Email,
				"avatar":        user.Avatar,
				"phone":         user.Phone,
				"created_at":    user.Created_at,
				"updated_at":    user.Updated_at,
			})
		} else {
			//find a user with that email and see if that user even exists
			err := userCollection.FindOne(ctx, bson.M{"email": body.Email}).Decode(&foundUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found, login seems to be incorrect"})
				return
			}

			//if all goes well, then you'll generate tokens

			token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_id)

			// save push token
			if user.Push_Notification_Token != nil {
				_, err := userCollection.UpdateOne(ctx, bson.M{"user_id": foundUser.User_id}, bson.M{"$set": bson.M{"push_notification_token": pushToken}})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}

			//update tokens - token and refresh token
			helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

			//return statusOK
			c.JSON(http.StatusOK, gin.H{
				"_id":           foundUser.ID,
				"token":         foundUser.Token,
				"refresh_token": foundUser.Refresh_Token,
				"first_name":    foundUser.First_Name,
				"last_name":     foundUser.Last_Name,
				"email":         foundUser.Email,
				"avatar":        foundUser.Avatar,
				"phone":         foundUser.Phone,
				"created_at":    foundUser.Created_at,
				"updated_at":    foundUser.Updated_at,
			})

		}

	}
}

func UpdateProfileImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId := c.Param("user_id")

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to resend otp, please try again later"})
			return
		}

		userFirstName := user.First_Name
		userLastName := user.Last_Name
		userPhone := user.Phone

		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		_, userErr := userCollection.UpdateByID(ctx, userId, bson.M{"$set": bson.M{"first_name": userFirstName, "last_name": userLastName, "phone": userPhone, "updated_at": user.Updated_at}})
		if userErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": userErr.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "User Updated Successfully!"})
	}
}

func LogOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId := c.Param("user_id")

		_, err := userCollection.UpdateByID(ctx, userId, bson.M{"$set": bson.M{"token": nil, "refresh_token": nil}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User Logged Out Successfully!"})
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId := c.Param("user_id")

		result, err := userCollection.DeleteOne(ctx, bson.M{"_id": userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": "user deleted successfully", "deleted": result})

	}
}
