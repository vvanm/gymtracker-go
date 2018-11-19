package diary

import (
	"github.com/gin-gonic/gin"
	"github.com/ravendb/ravendb-go-client"
	"github.com/vvanm/gymtracker-go/raven"
	"github.com/vvanm/gymtracker-go/routine"
	"github.com/vvanm/gymtracker-go/util/jwt"
)

func AddRoutine(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)
	date := c.Param("date")
	diaryKey := claims.ID + "/diaries/" + date

	var r routine.Routine
	c.ShouldBind(&r)

	exercises := make([]Diary_Exercise, len(r.Exercises))

	for i, v := range r.Exercises {
		exercises[i] = Diary_Exercise{
			ExerciseID: v,
			Sets:       []Diary_Exercise_Set{},
		}
	}

	//create patch
	patch := ravendb.PatchRequest_forScript(`this.routineID = args.routineID;this.exercises = args.exercises`)

	//add values
	patch.SetValues(map[string]interface{}{
		"routineID": r.RoutineID,
		"exercises": exercises,
	},
	)

	patchOp := ravendb.NewPatchOperation(diaryKey, nil, patch, nil, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, gin.H{})

}
