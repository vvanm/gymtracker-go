package routine

import (
	"github.com/gin-gonic/gin"
	"github.com/ravendb/ravendb-go-client"
	"github.com/vvanm/gymtracker-go/raven"
	"github.com/vvanm/gymtracker-go/util/jwt"
)

func AddExerciseToRoutine(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)

	//create patch
	patch := ravendb.PatchRequest_forScript(`
		this.entries.map(entry => {
			if(entry.routineID == args.routineID){
				entry.exercises.push(args.exerciseID)
			}
			return entry
		})`)

	//add values
	patch.SetValues(map[string]interface{}{
		"routineID":  c.Param("routineID"),
		"exerciseID": c.Param("exerciseID"),
	})

	//execute
	patchOp := ravendb.NewPatchOperation(claims.ID+"/routines", nil, patch, nil, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	c.JSON(200, gin.H{})
	return

}

func DelExerciseFromRoutine(c *gin.Context) {
	//create patch
	patch := ravendb.PatchRequest_forScript(`
		this.entries.map(entry => {
			if(entry.routineID == args.routineID){
				entry.exercises = entry.exercises.filter(ex => ex != args.exerciseID)
			}
			return entry
		})`,
	)

	//add values
	patch.SetValues(map[string]interface{}{
		"routineID":  c.Param("routineID"),
		"exerciseID": c.Param("exerciseID"),
	})

	claims := jwt.ClaimsFromCookie(c)

	//execute
	patchOp := ravendb.NewPatchOperation(claims.ID+"/routines", nil, patch, nil, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	c.JSON(200, gin.H{})
	return

}

func MoveExercise(c *gin.Context) {
	//create patch
	patch := ravendb.PatchRequest_forScript(`
		this.entries.map(entry => {
			if(entry.routineID == args.routineID){
				const i = entry.exercises.findIndex(ex => ex == args.exerciseID)
				const _i = args.dir == "up" ? i-1 : i+1
				const temp = entry.exercises[_i]
				
				entry.exercises[_i] = entry.exercises[i]
				entry.exercises[i] = temp 
			}
			return entry
		})
	`,
	)

	//add values
	patch.SetValues(map[string]interface{}{
		"routineID":  c.Param("routineID"),
		"exerciseID": c.Param("exerciseID"),
		"dir":        c.Param("dir"),
	})

	claims := jwt.ClaimsFromCookie(c)

	//execute
	patchOp := ravendb.NewPatchOperation(claims.ID+"/routines", nil, patch, nil, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	c.JSON(200, gin.H{})
	return

}
