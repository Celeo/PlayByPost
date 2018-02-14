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

func abortError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": "Could not access database",
		"error":   err.Error(),
	})
}

func timestamp() string {
	return time.Now().UTC().Format("Jan _2, 2006 @ 15:04:05")
}

func createPasswordHash(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), 0)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func checkHashAgainstPassword(hashed, raw string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	return err == nil, err
}

func createUUID() (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u4.String(), err
}

func getUserByID(id int) (User, error) {
	db := database()
	defer db.Close()
	u := User{}
	err := db.Get(&u, querySelectUserByID, id)
	return u, err
}

func getUserByName(name string) (User, error) {
	db := database()
	defer db.Close()
	u := User{}
	err := db.Get(&u, querySelectUserByName, name)
	return u, err
}

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

// TODO need to hook up the database in here so rolls persist
func textFormatWithDiceRolls(text string) (string, []Roll, error) {
	regexBBCode, err := regexp.Compile(`\[dice=([\w ]+)\]([\dd\+ ]+)\[/dice\]`)
	if err != nil {
		return "", nil, err
	}
	regexDice, err := regexp.Compile(`(\d+)d(\d+)`)
	if err != nil {
		return "", nil, err
	}
	regexMod, err := regexp.Compile(`([+-])(\d)+`)
	if err != nil {
		return "", nil, err
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
				return "", nil, err
			}
			diceSides, err := strconv.Atoi(dice[2])
			if err != nil {
				return "", nil, err
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
				return "", nil, err
			}
			if mod[1] == "+" {
				rollValue += val
			} else {
				rollValue -= val
			}
		}
		diceRollResults := regexDice.ReplaceAllString(rolls[2], "")
		text = strings.Replace(text, originalRolltext, fmt.Sprintf("<br>%s: %s ⇒ (%d)%s -> %d<br>", rolls[1], rolls[2], valueOfDice, diceRollResults, rollValue), 1)
	}
	return text, []Roll{}, nil
}

type parsedRoll struct {
	Post Post
	Dice map[int]Roll
}

// TODO
func insertRolls(p Post) (parsedRoll, error) {
	diceMap := make(map[int]Roll)
	// content, dice, err := textFormatWithDiceRolls(p.Content)

	return parsedRoll{
		Post: p,
		Dice: diceMap,
	}, nil
}