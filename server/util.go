package main

import (
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

const timeFormat string = "Jan _2, 2006 @ 15:04:05"
const editWindow time.Duration = time.Duration(30) * time.Minute

var regexDice = regexp.MustCompile(`(\d+)d(\d+)`)
var regexMod = regexp.MustCompile(`([+-])(\d+)`)
var regexD20 = regexp.MustCompile(`^.*: 1d20[^d]*$`)

// abortError sets the Gin header and body with a generic
// error message that includes the passed error's message.
func abortError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": "Internal error occurred",
		"error":   err.Error(),
	})
}

// timestamp returns the current time as a formatted string.
func timestamp() string {
	return time.Now().UTC().Format(timeFormat)
}

// readTimestamp converts a string timestamp to a time.Time struct.
func readTimestamp(str string) (time.Time, error) {
	return time.Parse(timeFormat, str)
}

// isPostWithinEditWindow returns true if the post is within the
// window of "recently posted" so that the user can edit it.
func isPostWithinEditWindow(p *Post) bool {
	t, err := readTimestamp(p.Date)
	if err != nil {
		return false
	}
	dateIsPast := time.Now().UTC().After(t)
	return dateIsPast && time.Since(t) < editWindow
}

// createPasswordHash takes a user's raw password string and
// returns the hashed version of it by running it through bcrypt.
func createPasswordHash(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), 0)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// checkHashAgainstPassword takes the hashed and raw passwords and
// calls bcrypt to see if they match.
func checkHashAgainstPassword(hashed, raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)) == nil
}

// createUUID creates and returns a new UUID.
func createUUID() (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u4.String(), err
}

// getUserByID searches for the database user model by the passed id.
func getUserByID(id int) (User, error) {
	db := database()
	defer db.Close()
	u := User{}
	err := db.Get(&u, querySelectUserByID, id)
	return u, err
}

// getUserByName searches for the database user model by the passed name.
func getUserByName(name string) (User, error) {
	db := database()
	defer db.Close()
	u := User{}
	err := db.Get(&u, querySelectUserByName, name)
	return u, err
}

// createSession creates a new database model for a new session for a user
// and returns the new session UUID.
func createSession(db *sqlx.DB, u User) (string, error) {
	uuid, err := createUUID()
	if err != nil {
		return "", err
	}
	if _, err := db.Exec(queryCreateSession, u.ID, uuid); err != nil {
		return "", err
	}
	return uuid, nil
}

// rollDice takes in a string from the frontend and returns an int value
// of what the rolls come out to. No database interaction.
func rollDice(str string) (int, error) {
	// declaration of vars and regex setup
	finalValue := 0
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// remove spaces from the string
	str = strings.Replace(str, " ", "", -1)
	// split into the separate groups of dice
	dice := strings.Split(str, ",")
	// iterate through all dice groups
	for _, die := range dice {
		dieResult := 0
		// split out the number before and after the 'd'
		groups := regexDice.FindStringSubmatch(die)
		count, err := strconv.Atoi(groups[1])
		if err != nil {
			return 0, err
		}
		sides, err := strconv.Atoi(groups[2])
		if err != nil {
			return 0, err
		}
		// "roll" the dice
		for i := 0; i < count; i++ {
			dieResult += rng.Intn(sides) + 1
		}
		// find the mod
		groups = regexMod.FindStringSubmatch(die)
		if len(groups) > 0 {
			// get the mod value and direction
			delta, err := strconv.Atoi(groups[2])
			if err != nil {
				return 0, err
			}
			// apply mod
			if groups[1] == "+" {
				dieResult += delta
			} else {
				dieResult -= delta
			}
		}
		// add the result of rolling these die with the rolling total
		finalValue += dieResult
	}
	return finalValue, nil
}

// injectD20Crits iterates through a list of roll structs and sets
// the `IsD20Crit` boolean to true when the dice is a single d20 roll
// with optional modifiers that yields a roll of 20.
func injectD20Crits(rolls []Roll) {
	for i := 0; i < len(rolls); i++ {
		if regexD20.Match([]byte(rolls[i].String)) {
			value := rolls[i].Value
			diceString := strings.Replace(strings.Split(rolls[i].String, ":")[1], " ", "", -1)
			groups := regexMod.FindStringSubmatch(diceString)
			if len(groups) > 0 {
				delta, err := strconv.Atoi(groups[2])
				if err != nil {
					continue
				}
				if groups[1] == "+" {
					value -= delta
				} else {
					value += delta
				}
			}
			if value == 20 {
				rolls[i].IsD20Crit = true
			}
		}
	}
}
