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
	uuid, err := registerNewAccount(&data)
	require.Nil(t, err)
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
	uuid, err := registerNewAccount(&data)
	require.Nil(t, err)
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
	uuid, err := registerNewAccount(&data)
	require.Nil(t, err)
	require.NotEmpty(t, uuid, "Returned uuid is empty")
	uuid, err = registerNewAccount(&data)
	require.NotNil(t, err, "No error thrown")
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
	uuid, err := registerNewAccount(&data)
	require.NotNil(t, err, "No error thrown")
	require.Contains(t, err.Error(), "Join code mismatch")
	require.Empty(t, uuid)
}

func TestLoginNoUser(t *testing.T) {
	newDB()
	data := loginData{
		Name:     "aaa",
		Password: "bbb",
	}
	uuid, err := login(&data)
	require.NotNil(t, err, "No error thrown")
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
	uuid, err := login(&data)
	require.Nil(t, err)
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
		Content: "Content",
		AuthID:  1,
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
		OldPassword: "",
		NewPassword: "",
		AuthID:      1,
	}
	err := changePassword(&data)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "Old password does not match")
}

func TestChangePassword(t *testing.T) {
	newDB()
	const newPass = "new-password"
	db := database()
	defer db.Close()
	addUser(db)
	data := newPasswordData{
		OldPassword: "password",
		NewPassword: newPass,
		AuthID:      1,
	}
	err := changePassword(&data)
	require.Nil(t, err)
	u, err := getUserByID(1)
	require.Nil(t, err)
	require.True(t, checkHashAgainstPassword(u.Password, newPass))
}
