package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

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
	return time.Now().UTC().Format("Jan _2, 2006 @ 15:04:05")
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
func createSession(u User) (string, error) {
	db := database()
	defer db.Close()
	uuid, err := createUUID()
	if err != nil {
		return "", err
	}
	if _, err := db.Exec(queryDeleteSessionsForUser, u.ID); err != nil {
		return "", err
	}
	if _, err := db.Exec(queryCreateSession, u.ID, uuid); err != nil {
		return "", err
	}
	return uuid, nil
}

// textFormatWithDiceRolls takes a post's raw, just-submitted content and uses
// a RNG to replace the user's "dice rolls" with actual values. This modified
// post content is then returned.
func textFormatWithDiceRolls(text string) (string, error) {
	regexBBCode, err := regexp.Compile(`(?i)\[dice=([\w ]+)\]([\dd\+\- ]+)\[/dice\]`)
	if err != nil {
		return "", err
	}
	regexDice, err := regexp.Compile(`(\d+)d(\d+)`)
	if err != nil {
		return "", err
	}
	regexMod, err := regexp.Compile(`([+-])(\d)+`)
	if err != nil {
		return "", err
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, rolls := range regexBBCode.FindAllStringSubmatch(text, -1) {
		rollValue := 0
		valueOfDice := 0
		originalRolltext := fmt.Sprintf("[dice=%s]%s[/dice]", rolls[1], rolls[2])
		rollText := strings.Replace(rolls[2], " ", "", -1)

		for _, dice := range regexDice.FindAllStringSubmatch(rollText, -1) {
			rollCount, err := strconv.Atoi(dice[1])
			if err != nil {
				return "", err
			}
			diceSides, err := strconv.Atoi(dice[2])
			if err != nil {
				return "", err
			}
			for i := 0; i < rollCount; i++ {
				rVal := rng.Intn(diceSides) + 1
				rollValue += rVal
				valueOfDice += rVal
			}
		}
		for _, mod := range regexMod.FindAllStringSubmatch(rollText, -1) {
			val, err := strconv.Atoi(mod[2])
			if err != nil {
				return "", err
			}
			if mod[1] == "+" {
				rollValue += val
			} else {
				rollValue -= val
			}
		}
		diceRollResults := regexDice.ReplaceAllString(rolls[2], "")
		text = strings.Replace(text, originalRolltext, fmt.Sprintf("<br><span style=\"color: green;\">%s: %s ⇒ (%d)%s -> %d</span><br>", rolls[1], rolls[2], valueOfDice, diceRollResults, rollValue), 1)
	}
	return text, nil
}

// insertRolls takes ata from a user and replaces the "dice rolls"
// with RNG values and modifies the passed struct with that new content.
func insertRolls(p *newPostData) error {
	content, err := textFormatWithDiceRolls(p.Content)
	if err != nil {
		return err
	}
	p.Content = content
	return nil
}
