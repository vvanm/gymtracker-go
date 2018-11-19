package main

import (
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/vvanm/gymtracker-go/diary"
	"github.com/vvanm/gymtracker-go/exercise"
	"github.com/vvanm/gymtracker-go/food"
	"github.com/vvanm/gymtracker-go/raven"
	"github.com/vvanm/gymtracker-go/routine"
	"github.com/vvanm/gymtracker-go/user"
	"github.com/vvanm/gymtracker-go/util/jwt"
	"os"
)

func main() {

	raven.SetupStore()

	port := "4000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	//Set up gin
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//Add cors
	r.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET,PUT,POST,DELETE,PATCH",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		Credentials:    true,
	}))

	//No route handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"errorMsg": "routeNotFound"})
	})

	//user endpoints
	r.GET("/session", user.GetSession)
	r.PUT("/sign-up", user.SignUp)
	r.POST("/sign-in", user.SignIn)

	//basic auth group
	rAuth := r.Group("/")
	rAuth.Use(jwt.BaseAuth())
	{
		//foods
		rAuth.POST("/foods", food.AddNew)
		rAuth.DELETE("/foods/:foodID", food.Delete)
		rAuth.PATCH("/foods/:foodID", food.Update)

		//exercises
		rAuth.POST("/exercises", exercise.AddNew)
		rAuth.PATCH("/exercises/:exerciseID", exercise.Update)
		rAuth.DELETE("/exercises/:exerciseID", exercise.Delete)

		//routines
		rAuth.POST("/routines", routine.AddNew)
		rAuth.POST("/routines/:routineID/exercises/:exerciseID", routine.AddExerciseToRoutine)
		rAuth.DELETE("/routines/:routineID/exercises/:exerciseID", routine.DelExerciseFromRoutine)
		rAuth.PATCH("/routines/:routineID/exercises/:exerciseID/:dir", routine.MoveExercise)

		//diary
		rAuth.GET("/diary/:date", diary.GetDiary)

		rAuth.POST("/diary/:date/food", diary.AddFood)
		rAuth.DELETE("/diary/:date/food/:i", diary.DeleteFood)

		rAuth.PUT("/diary/:date/routine", diary.AddRoutine)

		rAuth.POST("/diary/:date/exercises/:exerciseID/set", diary.AddSetToExercise)
		rAuth.DELETE("/diary/:date/exercises/:exerciseID/set/:i", diary.DeleteSetFromExercise)
	}

	//run the server
	r.Run(":" + port)

}
