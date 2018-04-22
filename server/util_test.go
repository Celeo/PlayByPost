package main

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimestamp(t *testing.T) {
	require.NotEmpty(t, timestamp())
}

func TestReadTimeStamp(t *testing.T) {
	val, err := readTimestamp("Feb 25, 2018 @ 23:33:04")
	require.Nil(t, err)
	require.Equal(t, val, time.Date(2018, 2, 25, 23, 33, 4, 0, time.UTC))
}

func TestIsPostWithinEditWindow(t *testing.T) {
	tests := []struct {
		String   string
		Expected bool
	}{
		{"Jan 01, 2018 @ 00:00:00", false},
		{time.Now().UTC().Add(time.Duration(-1) * time.Second).Format(timeFormat), true},
		{time.Now().UTC().Add(time.Duration(-1) * time.Second).Add(time.Duration(-30) * time.Minute).Format(timeFormat), false},
		{time.Now().UTC().Add(editWindow).Format(timeFormat), false},
	}
	for _, test := range tests {
		p := Post{Date: test.String}
		actual := isPostWithinEditWindow(&p)
		require.Equal(t, actual, test.Expected, "Should be '%t' for string '%s'", test.Expected, test.String)
	}
}

func TestCreatePasswordHash(t *testing.T) {
	h, err := createPasswordHash("test")
	require.Nil(t, err)
	require.NotEmpty(t, h)
}

func TestCreateUUID(t *testing.T) {
	h, err := createUUID()
	require.Nil(t, err)
	require.NotEmpty(t, h)
}

func TestCheckHashAgainstPassword(t *testing.T) {
	raw := "test"
	hashed, err := createPasswordHash(raw)
	require.Nil(t, err)
	m := checkHashAgainstPassword(hashed, raw)
	require.True(t, m)
}

func TestGetUserByID(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	u, err := getUserByID(1)
	require.Nil(t, err)
	require.NotNil(t, u)
	require.Equal(t, u.ID, 1)
}

func TestGetUserByIDNoUser(t *testing.T) {
	setDBToTest()
	u, err := getUserByID(1)
	require.NotNil(t, err)
	require.NotNil(t, u)
	require.Equal(t, u.ID, 0)
}

func TestGetUserByName(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	u, err := getUserByName("username")
	require.Nil(t, err)
	require.NotNil(t, u)
	require.Equal(t, u.ID, 1)
}

func TestGetUserByNameNoUser(t *testing.T) {
	setDBToTest()
	u, err := getUserByName("username")
	require.NotNil(t, err)
	require.NotNil(t, u)
	require.Equal(t, u.ID, 0)
}

func TestRollDice(t *testing.T) {
	rolls := []string{
		"abc: 1d20",
		"abc: 2d20",
		"abc: 1d6+3",
		"abc: 1d4 + 1, 2d6 -3, 1d2",
		"abc: 1d2- 100",
	}
	for _, roll := range rolls {
		r, err := rollDice(roll)
		require.Nil(t, err)
		require.NotNil(t, r)
	}
}

func TestCreateSession(t *testing.T) {
	setDBToTest()
	db := database()
	defer db.Close()
	addTestUser(db)
	u, err := getUserByID(1)
	require.Nil(t, err)
	uuid, err := createSession(db, u)
	require.Nil(t, err)
	require.NotNil(t, uuid)
}

func TestInjectD20Crits(t *testing.T) {
	yes := []Roll{
		{
			String: "something: 1d20",
			Value:  20,
		},
		{
			String: "something: 1d20 + 100",
			Value:  120,
		},
		{
			String: "something: 1d20 - 4",
			Value:  16,
		},
	}
	injectD20Crits(yes)
	for index, roll := range yes {
		require.True(t, roll.IsD20Crit, "index: "+strconv.Itoa(index))
	}
	no := []Roll{
		{
			String: "something: 1d6",
			Value:  4,
		},
		{
			String: "something: 1d20 + 1d4 + 1",
			Value:  5,
		},
		{
			String: "something: 1d20",
			Value:  5,
		},
	}
	injectD20Crits(no)
	for _, roll := range no {
		require.False(t, roll.IsD20Crit)
	}
}
