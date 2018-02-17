package main

import (
	"errors"
	"os"
)

// registerData is data required for creating a new user.
type registerData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}

// loginData is data required for logging a user into the app.
type loginData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// formattedPost is a processed Post struct, suitable
// for handing off to the front-end to show to the uesr.
type formattedPost struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

// newPostData is data required for creating a new post.
type newPostData struct {
	Content string `json:"content"`
	AuthID  int    `json:"-"`
}

// newPasswordData is data required for changing a user's password.
type newPasswordData struct {
	OldPassword string `json:"old"`
	NewPassword string `json:"new"`
	AuthID      int    `json:"-"`
}

// registerNewAccount takes a struct of data and creates a new user
// account so long as one with the same name does not already exist.
func registerNewAccount(data *registerData) (string, error) {
	joinCode := os.Getenv("JOIN_CODE")
	if joinCode == "" {
		joinCode = "join"
	}
	if data.Code != joinCode {
		return "", errors.New("Join code mismatch")
	}
	db := database()
	defer db.Close()

	existingUsers := []User{}
	if err := db.Select(&existingUsers, querySelectUserByName, data.Name); err != nil {
		return "", err
	}
	if len(existingUsers) != 0 {
		return "", errors.New("Username not unique")
	}

	hashedPassword, err := createPasswordHash(data.Password)
	if err != nil {
		return "", err
	}
	_, err = db.Exec(queryCreateUser, data.Name, hashedPassword, data.Email)
	if err != nil {
		return "", err
	}
	u, err := getUserByName(data.Name)
	if err != nil {
		return "", err
	}
	uuid, err := createSession(u)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

// login takes a struct of data, checks the user's password against the
// hashed password from the database, and if it matches, creates a new
// session for the user, the uuid of which is returned for storage in
// the localSession on the front-end.
func login(data *loginData) (string, error) {
	db := database()
	defer db.Close()
	u, err := getUserByName(data.Name)
	if err != nil {
		return "", err
	}
	passwordMatch := checkHashAgainstPassword(u.Password, data.Password)
	if !passwordMatch {
		return "", errors.New("Password mismatch")
	}
	uuid, err := createSession(u)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

// getAllPosts pulls all the posts out of the database, copies them
// each into a new struct, injecting the poster's name instead of
// their id, and returns the slice.
func getAllPosts() ([]formattedPost, error) {
	posts := []Post{}
	users := []User{}
	userMap := make(map[int]User)
	retVal := []formattedPost{}
	db := database()
	defer db.Close()
	if err := db.Select(&posts, querySelectPosts); err != nil {
		return nil, err
	}
	if err := db.Select(&users, queryselectUsers); err != nil {
		return nil, err
	}
	for _, u := range users {
		userMap[u.ID] = u
	}
	for _, p := range posts {
		err := insertRolls(&p)
		if err != nil {
			return nil, err
		}
		retVal = append(retVal, formattedPost{
			ID:      p.ID,
			Name:    userMap[p.UserID].Name,
			Date:    p.Date,
			Content: p.Content,
		})
	}
	return retVal, nil
}

// createNewPost takes a struct of data and creates a new post by that
// user with that content.
func createNewPost(data *newPostData) error {
	db := database()
	defer db.Close()
	_, err := db.Exec(queryCreatePost, data.AuthID, timestamp(), data.Content)
	return err
}

// changePassword changes a user's password if the supplied old password
// matches what's already in the database.
func changePassword(data *newPasswordData) error {
	db := database()
	defer db.Close()
	u, err := getUserByID(data.AuthID)
	if err != nil {
		return err
	}
	oldPasswordMatches := checkHashAgainstPassword(u.Password, data.OldPassword)
	if !oldPasswordMatches {
		return errors.New("Old password does not match")
	}
	newHash, err := createPasswordHash(data.NewPassword)
	if err != nil {
		return err
	}
	_, err = db.Exec(queryUpdatePassword, newHash, u.ID)
	return err
}
