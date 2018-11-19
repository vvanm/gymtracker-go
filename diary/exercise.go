package diary

import (
	"github.com/gin-gonic/gin"
	"github.com/ravendb/ravendb-go-client"
	"github.com/vvanm/gymtracker-go/raven"

	"github.com/vvanm/gymtracker-go/util/jwt"
)

func AddSetToExercise(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)
	date := c.Param("date")
	diaryKey := claims.ID + "/diaries/" + date

	var s Diary_Exercise_Set
	c.ShouldBind(&s)

	//create patch
	patch := ravendb.PatchRequest_forScript(`
		this.exercises = this.exercises.map(ex => {
			if(ex.exerciseID == args.exerciseID){
				ex.sets.push(args.set)
			}		
			return ex
		})
	`)

	//add values
	patch.SetValues(map[string]interface{}{
		"exerciseID": c.Param("exerciseID"),
		"set":        s,
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

func DeleteSetFromExercise(c *gin.Context) {

	claims := jwt.ClaimsFromCookie(c)
	date := c.Param("date")
	diaryKey := claims.ID + "/diaries/" + date

	//create patch

	patch := ravendb.PatchRequest_forScript(`
		this.exercises.map(ex => {
			if(ex.exerciseID == args.exerciseID){
				 ex.sets.splice(args.i,1)
			}
			return ex
		})`,
	)

	//add values
	patch.SetValues(map[string]interface{}{
		"i":          c.Param("i"),
		"exerciseID": c.Param("exerciseID"),
	})

	patchOp := ravendb.NewPatchOperation(diaryKey, nil, patch, nil, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, gin.H{})

}
