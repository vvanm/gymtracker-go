package user

import (
	"github.com/gin-gonic/gin"
	"github.com/vvanm/gymtracker-go/util/auth"
	"github.com/vvanm/gymtracker-go/util/jwt"

	"net/http"
)

var projectionSession = `
	declare function projection(user,ID){
		let exercises = load(ID+"/exercises")
		let foods = load(ID+"/foods");
		let routines = load(ID+"/routines")

		return {
				ID,
				password : user.password,
				foods : foods.entries !== null ? foods.entries : [],
				exercises : exercises.entries !== null ? exercises.entries : [],
				routines : routines.entries !== null ? routines.entries : []
		}
	}
`

func GetSession(c *gin.Context) {
	//Find claims
	claims := jwt.ClaimsFromCookie(c)

	if claims == nil {
		c.JSON(404, gin.H{"errorMsg": "noSession"})
		return
	}

	u, err := GetUser(
		projectionSession + `
		from index 'users/search' as entry
		where entry.ID == '` + claims.ID + `'
		select projection(entry,Id())
		`,
	)

	if err != nil {
		c.JSON(400, gin.H{"errorMsg": err.Error()})
		return
	}

	if u == nil {
		c.JSON(400, gin.H{"errorMsg": "no user found"})
		return
	}

	u.Password = ""

	c.JSON(200, u)

}

func SignUp(c *gin.Context) {
	//receive post data
	var u User
	c.BindJSON(&u)
	err := u.Create()
	if err != nil {
		c.JSON(400, gin.H{"errorMsg": err.Error()})
	}

	c.JSON(200, gin.H{})

}

func SignIn(c *gin.Context) {
	//Receive post data
	var postU User
	c.BindJSON(&postU)

	u, err := GetUser(
		projectionSession + `
			from index 'users/search' as entry
			where entry.email ==  '` + postU.Email + `'
			select projection(entry,Id())	
		`,
	)

	if err != nil {
		c.JSON(400, gin.H{"errorMsg": err.Error()})
	}

	//Check password match
	if !auth.PasswordMatch(postU.Password, u.Password) {
		c.JSON(400, gin.H{"errorMsg": "failAuth"})
		return
	}

	u.Password = ""
	//Create claims
	claims := jwt.NewClaims(u.ID)

	//Create token
	token, _ := jwt.NewToken(claims)

	//Create cookie
	cookie := jwt.NewCookie(token)

	//Set cookie
	http.SetCookie(c.Writer, cookie)

	c.JSON(200, u)

}
