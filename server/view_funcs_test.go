package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterNewAccountDefaultJoinCode(t *testing.T) {
	newDB()
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
	newDB()
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
	newDB()
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
	newDB()
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
	newDB()
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
	newDB()
	db := database()
	defer db.Close()
	addUser(db)
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
	newDB()
	posts, err := getAllPosts()
	require.Nil(t, err)
	require.Empty(t, posts, "Magically found posts")
}

func TestGetAllPosts(t *testing.T) {
	newDB()
	db := database()
	defer db.Close()
	addUser(db)
	addPost(db)
	posts, err := getAllPosts()
	require.Nil(t, err)
	require.Equal(t, len(posts), 1, "Incorrect number of posts returned")
	require.Equal(t, posts[0].Name, "username")
}

func TestCreatePost(t *testing.T) {
	newDB()
	db := database()
	defer db.Close()
	addUser(db)
	data := newPostData{
		ID:      1,
		Content: "Content",
	}
	err := createNewPost(&data)
	require.Nil(t, err)
}

func TestChangePasswordMismatch(t *testing.T) {
	newDB()
	db := database()
	defer db.Close()
	addUser(db)
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
	newDB()
	db := database()
	defer db.Close()
	addUser(db)
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
	newDB()
	const newPass = "new-password"
	db := database()
	defer db.Close()
	addUser(db)
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
