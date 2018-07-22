package main

import (
	"errors"
	"os"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
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
// for handing off to the front-end to show to the user.
type formattedPost struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Date          string `json:"date"`
	Tag           string `json:"tag"`
	Content       string `json:"content"`
	Rolls         []Roll `json:"rolls"`
	EditingWindow bool   `json:"canEdit"`
}

// newPostData is data required for creating a new post.
type newPostData struct {
	ID      int    `json:"-"`
	Tag     string `json:"-"`
	Content string `json:"content"`
}

// newPasswordData is data required for changing a user's password.
type newPasswordData struct {
	ID          int    `json:"-"`
	OldPassword string `json:"old"`
	NewPassword string `json:"new"`
}

// updateUserData is data required to update a user's database model.
type updateUserData struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PostsPerPage string `json:"postsPerPage"`
	NewestAtTop  bool   `json:"newestAtTop"`
	Tag          string `json:"tag"`
}

// addRollData is data required for storing a new roll.
type addRollData struct {
	ID     int    `json:"-"`
	String string `json:"roll"`
}

// editPostData is data required for editing a post.
type editPostData struct {
	ID      int    `json:"-"`
	Content string `json:"content"`
}

// invalidLoginsData is data required for invalidating logins.
type invalidLoginsData struct {
	ID   int
	UUID string
}

// registerNewAccount takes a struct of data and creates a new user
// account so long as one with the same name does not already exist.
func registerNewAccount(data *registerData) (string, User, error) {
	joinCode := os.Getenv("JOIN_CODE")
	if joinCode == "" {
		joinCode = "join"
	}
	if data.Code != joinCode {
		return "", User{}, errors.New("Join code mismatch")
	}
	db := database()
	defer db.Close()

	existingUsers := []User{}
	if err := db.Select(&existingUsers, querySelectUserByName, data.Name); err != nil {
		return "", User{}, err
	}
	if len(existingUsers) != 0 {
		return "", User{}, errors.New("Username not unique")
	}

	hashedPassword, err := createPasswordHash(data.Password)
	if err != nil {
		return "", User{}, err
	}
	_, err = db.Exec(queryCreateUser, data.Name, hashedPassword, data.Email)
	if err != nil {
		return "", User{}, err
	}
	u, err := getUserByName(data.Name)
	if err != nil {
		return "", User{}, err
	}
	uuid, err := createSession(db, u)
	if err != nil {
		return "", User{}, err
	}
	return uuid, u, nil
}

// login takes a struct of data, checks the user's password against the
// hashed password from the database, and if it matches, creates a new
// session for the user, the uuid of which is returned for storage in
// the localSession on the front-end along with the user.
func login(data *loginData) (string, User, error) {
	db := database()
	defer db.Close()
	u, err := getUserByName(data.Name)
	if err != nil {
		return "", User{}, err
	}
	passwordMatch := checkHashAgainstPassword(u.Password, data.Password)
	if !passwordMatch {
		return "", User{}, errors.New("Password mismatch")
	}
	uuid, err := createSession(db, u)
	if err != nil {
		return "", User{}, err
	}
	return uuid, u, nil
}

// getAllPostIDs returns a slice of all the database post IDs.
func getAllPostIDs() ([]int, error) {
	db := database()
	defer db.Close()
	ids := []int{}
	if err := db.Select(&ids, querySelectAllPostIDs); err != nil {
		return nil, err
	}
	return ids, nil
}

// createNewPost takes a struct of data and creates a new post by that
// user with that content.
func createNewPost(data *newPostData) error {
	db := database()
	defer db.Close()
	impacted, err := db.Exec(queryCreatePost, data.ID, timestamp(), data.Tag, data.Content)
	newPostID, err := impacted.LastInsertId()
	if err != nil {
		return err
	}
	_, err = db.Exec(querySavePendingRoll, newPostID, data.ID)
	return err
}

// updateUserInformation takes a struct of data and updates the database
// user model that matches the ID with the new information. This method
// does not handle password changes.
func updateUserInformation(data *updateUserData) (User, error) {
	db := database()
	defer db.Close()
	existing, err := getUserByID(data.ID)
	if err != nil {
		return User{}, err
	}
	pppStr := data.PostsPerPage
	if pppStr == "0" {
		pppStr = strconv.Itoa(existing.PostsPerPage)
	}
	ppp, err := strconv.Atoi(pppStr)
	if err != nil {
		return User{}, err
	}
	name := data.Name
	if len(name) == 0 {
		name = existing.Name
	}
	_, err = db.Exec(queryUpdateUser, name, data.Email, ppp, data.NewestAtTop, data.Tag, data.ID)
	if err != nil {
		return User{}, err
	}
	return getUserByID(data.ID)
}

// changePassword changes a user's password if the supplied old password
// matches what's already in the database.
func changePassword(data *newPasswordData) error {
	db := database()
	defer db.Close()
	u, err := getUserByID(data.ID)
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

// getPendingDice returns all pending rolls for the specified user
func getPendingDice(authID int) ([]Roll, error) {
	db := database()
	defer db.Close()
	rolls := []Roll{}
	err := db.Select(&rolls, queryGetPendingRollsForUser, authID)
	return rolls, err
}

// addPendingDie takes a user's roll, computes the value, and stores
// it into the database. All the user's pending rolls are returned.
func addPendingDie(data *addRollData) ([]Roll, error) {
	val, err := rollDice(data.String)
	if err != nil {
		return nil, err
	}
	db := database()
	defer db.Close()
	_, err = db.Exec(queryInsertPendingRoll, data.ID, data.String, val)
	if err != nil {
		return nil, err
	}
	return getPendingDice(data.ID)
}

// getRawPostByID returns the post struct matching the id.
func getRawPostByID(id int) (Post, error) {
	db := database()
	defer db.Close()
	p := Post{}
	err := db.Get(&p, querySelectSinglePost, id)
	return p, err
}

// getPostByID returns a single Post struct from the database by its id.
func getPostByID(id int) (formattedPost, error) {
	db := database()
	defer db.Close()
	post := Post{}
	rolls := []Roll{}
	user := User{}
	if err := db.Get(&post, querySelectSinglePost, id); err != nil {
		return formattedPost{}, err
	}
	if err := db.Select(&rolls, querySelectRollsForPostID, id); err != nil {
		return formattedPost{}, err
	}
	if err := db.Get(&user, querySelectUserByID, post.UserID); err != nil {
		return formattedPost{}, err
	}
	injectD20Crits(rolls)
	return formattedPost{
		ID:            post.ID,
		Name:          user.Name,
		Date:          post.Date,
		Content:       post.Content,
		Tag:           post.Tag,
		Rolls:         rolls,
		EditingWindow: isPostWithinEditWindow(&post),
	}, nil
}

// editPost takes data to modify the content of a post, and saves the
// changes to the database. No modification of rolls connected to the
// post are made.
func editPost(data *editPostData) error {
	db := database()
	defer db.Close()
	_, err := db.Exec(queryEditPost, data.Content, data.ID)
	return err
}

// clearLogins deletes all stored sessions from the data for the user
// other than the single uuid that's passed in as part of the data.
func clearLogins(data *invalidLoginsData) error {
	db := database()
	defer db.Close()
	_, err := db.Exec(queryInvalidLogins, data.ID, data.UUID)
	return err
}

// Returns an array of all post ids that match the fuzzy search content string.
func searchPosts(needle string) ([]int, error) {
	se := search.New(language.English, search.IgnoreCase)
	allPosts := []Post{}
	retPosts := make([]int, 0)
	if len(needle) == 0 {
		return retPosts, nil
	}
	db := database()
	defer db.Close()
	if err := db.Select(&allPosts, querySelectAllPosts); err != nil {
		return retPosts, err
	}
	for _, post := range allPosts {
		if iS, iE := se.IndexString(post.Content, needle); iS != -1 && iE != -1 {
			retPosts = append(retPosts, post.ID)
		}
	}
	return retPosts, nil
}
