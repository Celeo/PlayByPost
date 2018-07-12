package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

const contextAuthID string = "authID"
const contextAuthName string = "authName"
const contextAuthUUID string = "uuid"

// main is the entry point for the app.
func main() {
	createTables()
	r := gin.Default()
	r.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type",
		Credentials:    true,
	}))

	r.GET("/", viewIndex)
	r.POST("/login", viewLogin)
	r.POST("/register", viewRegister)
	r.GET("/posts", viewAllPostIDs)
	r.POST("/posts", middlewareLoggedIn(), viewCreatePost)
	r.POST("/posts/search/:needle", middlewareLoggedIn(), viewSearchPosts)
	r.GET("/post/:id", viewGetSinglePost)
	r.PUT("/post/:id", middlewareLoggedIn(), viewEditPost)
	r.GET("/roll", middlewareLoggedIn(), viewGetPendingDice)
	r.POST("/roll", middlewareLoggedIn(), viewRollDice)
	r.GET("/profile", middlewareLoggedIn(), viewGetProfile)
	r.PUT("/profile", middlewareLoggedIn(), viewUpdateUser)
	r.PUT("/profile/password", middlewareLoggedIn(), viewChangePassword)
	r.POST("/profile/invalidate", middlewareLoggedIn(), viewInvalidLogins)

	r.Run(":5000")
}

// middlewareLoggedIn returns a Gin handler function that asserts
// that the incoming HTTP requests to the restricted endpoints are
// from an authenticated user. If they aren't, the the request is aborted
// with a HTTP status code that explains why the user was not allowed
// access to that endpoint.
func middlewareLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.GetHeader("Authorization")
		if uuid == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		s := Session{}
		db := database()
		defer db.Close()
		if err := db.Get(&s, querySelectSessionByUUID, uuid); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		u, err := getUserByID(s.UserID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(contextAuthID, u.ID)
		c.Set(contextAuthName, u.Name)
		c.Set(contextAuthUUID, uuid)
		c.Next()
	}
}

func viewIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Index page",
	})
}

func viewRegister(c *gin.Context) {
	data := registerData{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	uuid, user, err := registerNewAccount(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not register account",
			"error":   err.Error,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "Registration successful",
		"uuid":         uuid,
		"name":         data.Name,
		"postsPerPage": user.PostsPerPage,
		"newestAtTop":  user.NewestAtTop,
	})
}

func viewLogin(c *gin.Context) {
	data := loginData{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	uuid, user, err := login(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not complete the login",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"uuid":         uuid,
		"name":         data.Name,
		"postsPerPage": user.PostsPerPage,
		"newestAtTop":  user.NewestAtTop,
		"tag":          user.Tag,
	})
}

func viewAllPostIDs(c *gin.Context) {
	ids, err := getAllPostIDs()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get all ids from the database",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ids)
}

func viewCreatePost(c *gin.Context) {
	data := newPostData{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	data.ID = c.GetInt(contextAuthID)
	u, err := getUserByID(data.ID)
	if err != nil {
		abortError(c, err)
		return
	}
	data.Tag = u.Tag
	if err := createNewPost(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "New post could not be created",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Message posted successfully",
	})
}

func viewSearchPosts(c *gin.Context) {
	needle := c.Param("needle")
	posts, err := searchPosts(needle)
	if err != nil {
		abortError(c, err)
		return
	}
	c.JSON(http.StatusOK, posts)
}

func viewGetPendingDice(c *gin.Context) {
	authID := c.GetInt(contextAuthID)
	rolls, err := getPendingDice(authID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get pending dice",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, rolls)
}

func viewRollDice(c *gin.Context) {
	data := addRollData{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	data.ID = c.GetInt(contextAuthID)
	rolls, err := addPendingDie(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not roll and store dice",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, rolls)
}

func viewGetProfile(c *gin.Context) {
	u, err := getUserByID(c.GetInt(contextAuthID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get profile for user",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, u)
}

func viewUpdateUser(c *gin.Context) {
	data := updateUserData{}
	if err := c.ShouldBindWith(&data, binding.JSON); err != nil {
		abortError(c, err)
		return
	}
	data.ID = c.GetInt(contextAuthID)
	u, err := updateUserInformation(&data)
	if err != nil {
		abortError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

func viewChangePassword(c *gin.Context) {
	data := newPasswordData{}
	if err := c.BindJSON(&data); err != nil {
		abortError(c, err)
		return
	}
	data.ID = c.GetInt(contextAuthID)
	if err := changePassword(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Updating password failed",
			"error":   err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

func viewInvalidLogins(c *gin.Context) {
	data := invalidLoginsData{
		ID:   c.GetInt(contextAuthID),
		UUID: c.GetString(contextAuthUUID),
	}
	err := clearLogins(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not invalidate other logins",
			"error":   err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

func viewGetSinglePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "ID is not a number",
			"error":   err.Error(),
		})
		return
	}
	post, err := getPostByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get post from database",
			"error":   err.Error(),
		})
		return
	}
	// If a post is no longer editable by users, it's highly unlikely that
	// it's going to change, so add a caching header to the response.
	if !post.EditingWindow {
		c.Header("Cache-Control", "max-age=604800") // 7 days
	}
	c.JSON(http.StatusOK, post)
}

func viewEditPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "ID is not a number",
			"error":   err.Error(),
		})
		return
	}
	data := editPostData{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	data.ID = id
	p, err := getRawPostByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if p.UserID != c.GetInt(contextAuthID) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if !isPostWithinEditWindow(&p) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	err = editPost(&data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not save the modified post",
			"error":   err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}
