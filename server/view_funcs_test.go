package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterNewAccountDefaultJoinCode(t *testing.T) {
	setDBToTest()
	data := registerData{
		Name:     "aaa",
		Password: "bbb",
		Email:    "ccc",
		Code:     "join",
	}
	uuid, u, err := registerNewAccount(&data)
	require.Nil(t, err)
	require.NotEqual(t, u.ID, 0)
	require.NotEmpty(t, uuid, "Returned uuid is empty")
}

func TestRegisterNewAccountDefaultNewCode(t *testing.T) {
	setDBToTest()
	const code = "foobar"
	os.Setenv("JOIN_CODE", code)
	data := registerData{
		Name:     "aaa",
		Password: "bbb",
		Email:    "ccc",
		Code:     code,
	}
	uuid, u, err := registerNewAccount(&data)
	require.Nil(t, err)
	require.NotEqual(t, u.ID, 0)
	require.NotEmpty(t, uuid, "Returned uuid is empty")
	os.Unsetenv("JOIN_CODE")
}

func TestRegisterNewAccountUserExists(t *testing.T) {
	setDBToTest()
	data := registerData{
		Name:     "aaa",
		Password: "bbb",
		Email:    "ccc",
		Code:     "join",
	}
	uuid, u, err := registerNewAccount(&data)
	require.Nil(t, err)
	require.NotEqual(t, u.ID, 0)
	require.NotEmpty(t, uuid, "Returned uuid is empty")
	uuid, u, err = registerNewAccount(&data)
	require.NotNil(t, err, "No error thrown")
	require.Equal(t, u.ID, 0)
	require.Contains(t, err.Error(), "Username not unique")
}

func TestRegisterNewAccountInvalidJoinCode(t *testing.T) {
	setDBToTest()
	data := registerData{
		Name:     "aaa",
		Password: "bbb",
		Email:    "ccc",
		Code:     "invalid-code",
	}
	uuid, u, err := registerNewAccount(&data)
	require.NotNil(t, err, "No error thrown")
	require.Contains(t, err.Error(), "Join code mismatch")
	require.Equal(t, u.ID, 0)
	require.Empty(t, uuid)
}

func TestLoginNoUser(t *testing.T) {
	setDBToTest()
	data := loginData{
		Name:     "aaa",
		Password: "bbb",
	}
	uuid, u, err := login(&data)
	require.NotNil(t, err, "No error thrown")
	require.Equal(t, u.ID, 0)
	require.Empty(t, uuid)
}

func TestLoinWithUser(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	data := loginData{
		Name:     "username",
		Password: "password",
	}
	uuid, u, err := login(&data)
	require.Nil(t, err)
	require.NotEqual(t, u.ID, 0)
	require.NotEmpty(t, uuid, "UUID is blank")
}

func TestGetAllPostsNoPosts(t *testing.T) {
	setDBToTest()
	posts, err := getAllPosts()
	require.Nil(t, err)
	require.Empty(t, posts, "Magically found posts")
}

func TestGetAllPosts(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	addTestPost(db)
	posts, err := getAllPosts()
	require.Nil(t, err)
	require.Equal(t, len(posts), 1, "Incorrect number of posts returned")
	require.Equal(t, posts[0].Name, "username")
}

func TestCreatePost(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	data := newPostData{
		ID:      1,
		Content: "Content",
	}
	err := createNewPost(&data)
	require.Nil(t, err)
}

func TestChangePasswordMismatch(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	data := newPasswordData{
		ID:          1,
		OldPassword: "",
		NewPassword: "",
	}
	err := changePassword(&data)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "Old password does not match")
}

func TestUpdateUser(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	data := updateUserData{
		ID:           1,
		Name:         "new-name",
		Email:        "new-email",
		PostsPerPage: "5",
		NewestAtTop:  true,
	}
	u, err := updateUserInformation(&data)
	require.Nil(t, err)
	require.Equal(t, u.ID, 1)
	require.Equal(t, u.Name, "new-name")
	require.Equal(t, u.Email, "new-email")
	require.Equal(t, u.PostsPerPage, 5)
	require.Equal(t, u.NewestAtTop, true)
}

func TestChangePassword(t *testing.T) {
	setDBToTest()
	const newPass = "new-password"
	db := database()
	defer db.Close()
	addTestUser(db)
	data := newPasswordData{
		ID:          1,
		OldPassword: "password",
		NewPassword: newPass,
	}
	err := changePassword(&data)
	require.Nil(t, err)
	u, err := getUserByID(1)
	require.Nil(t, err)
	require.True(t, checkHashAgainstPassword(u.Password, newPass))
}

func TestAddingGettingDice(t *testing.T) {
	setDBToTest()
	const roll = "abc: 1d20 + 3"
	db := database()
	defer db.Close()
	addTestUser(db)
	data := addRollData{
		ID:     1,
		String: roll,
	}
	rolls, err := addPendingDie(&data)
	require.Nil(t, err)
	require.Equal(t, len(rolls), 1)
	require.Equal(t, rolls[0].String, roll)
}

func TestSaveRollsOnPostCreate(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	pendingRollData := []addRollData{
		{
			ID:     1,
			String: "abc: 1d20",
		},
		{
			ID:     1,
			String: "abc: 1d20 + 3",
		},
		{
			ID:     1,
			String: "abc: 2d6 + 3, 1d4 - 1",
		},
	}
	for i := 0; i < len(pendingRollData); i++ {
		accum, err := addPendingDie(&pendingRollData[i])
		require.Nil(t, err)
		require.Equal(t, len(accum), i+1)
	}
	addTestPost(db)
	posts, err := getAllPosts()
	require.Nil(t, err)
	require.Equal(t, len(posts), 1)
	for _, roll := range posts[0].Rolls {
		require.Equal(t, roll.Pending, false)
	}
	require.True(t, posts[0].EditingWindow)
}

func TestGetPostByID(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	addTestPost(db)
	post, err := getPostByID(1)
	require.Nil(t, err)
	require.Equal(t, post.ID, 1)
}

func TestClearLogins(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	u, err := getUserByID(1)
	require.Nil(t, err)
	require.NotNil(t, u)
	for i := 0; i < 3; i++ {
		uuid, err := createSession(db, u)
		require.Nil(t, err)
		require.NotNil(t, uuid)
	}
	sessions := []Session{}
	err = db.Select(&sessions, queryGetAllSessions)
	require.Nil(t, err)
	savedUUID := sessions[0].UUID
	clearLogins(&invalidLoginsData{1, savedUUID})
	newSessions := []Session{}
	err = db.Select(&newSessions, queryGetAllSessions)
	require.Nil(t, err)
	require.Equal(t, len(newSessions), 1)
	require.Equal(t, newSessions[0].ID, 1)
	require.Equal(t, newSessions[0].UUID, savedUUID)
}
