package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

// main is the entry point for the app
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
	r.GET("/post", viewPosts)
	r.GET("/post/:id", middlewareLoggedIn(), viewGetSinglePost)
	r.PUT("/post/:id", middlewareLoggedIn(), viewEditPost)
	r.POST("/post", middlewareLoggedIn(), viewCreatePost)
	r.GET("/roll", middlewareLoggedIn(), viewGetPendingDice)
	r.POST("/roll", middlewareLoggedIn(), viewRollDice)
	r.GET("/profile", middlewareLoggedIn(), viewGetProfile)
	r.PUT("/profile", middlewareLoggedIn(), viewUpdateUser)
	r.PUT("/profile/password", middlewareLoggedIn(), viewChangePassword)

	r.Run(":5000")
}

// middlewareLoggedIn returns a Gin handler function that asserts
// the the incoming HTTP requests to the restricted endpoints are
// from an authenticated user. If they aren't, the the request is aborted
// with a HTTP status code that explains why the user was not allowed
// access to that endpoint.
func middlewareLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		s := Session{}
		db := database()
		defer db.Close()
		if err := db.Get(&s, querySelectSessionByUUID, header); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		u, err := getUserByID(s.UserID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("authName", u.Name)
		c.Set("authID", u.ID)
		c.Next()
	}
}

// endpoint: GET /
func viewIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Index page",
	})
}

// endpoint: GET /register
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

// endpoint: GET /login
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
	})
}

// endpoint: GET /post
func viewPosts(c *gin.Context) {
	posts, err := getAllPosts()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get all posts from the database",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// endpoint: POST /post
func viewCreatePost(c *gin.Context) {
	data := newPostData{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	data.ID = c.GetInt("authID")
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

// endpoint: GET /roll
func viewGetPendingDice(c *gin.Context) {
	authID := c.GetInt("authID")
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

// endpoint: POST /roll
func viewRollDice(c *gin.Context) {
	data := addRollData{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	data.ID = c.GetInt("authID")
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

// endpoint: GET /profile
func viewGetProfile(c *gin.Context) {
	u, err := getUserByID(c.GetInt("authID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get profile for user",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, u)
}

// endpoint: PUT /profile
func viewUpdateUser(c *gin.Context) {
	data := updateUserData{}
	if err := c.ShouldBindWith(&data, binding.JSON); err != nil {
		abortError(c, err)
		return
	}
	data.ID = c.GetInt("authID")
	u, err := updateUserInformation(&data)
	if err != nil {
		abortError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

// endpoint: PUT /profile/password
func viewChangePassword(c *gin.Context) {
	data := newPasswordData{}
	if err := c.BindJSON(&data); err != nil {
		abortError(c, err)
		return
	}
	data.ID = c.GetInt("authID")
	if err := changePassword(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Updating password failed",
			"error":   err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// endpoint: GET /post/:id
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
	if !isPostWithinEditWindow(&post) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.JSON(http.StatusOK, post)
}

// endpoint: PUT /post/:id
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
	p, err := getPostByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if p.UserID != c.GetInt("authID") {
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
