package exercise

import (
	"github.com/gin-gonic/gin"
	"github.com/ravendb/ravendb-go-client"
	"github.com/vvanm/gymtracker-go/raven"
	"github.com/vvanm/gymtracker-go/util/helpers"
	"github.com/vvanm/gymtracker-go/util/jwt"
)

type Exercise struct {
	ExerciseID string `json:"exerciseID,omitempty"`
	//
	Name string `json:"name,omitempty"`
}

func Delete(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)

	patch := ravendb.PatchRequest_forScript(`this.entries = this.entries.filter(e => e.exerciseID != args.exerciseID)`)
	patch.SetValues(map[string]interface{}{
		"exerciseID": c.Param("exerciseID"),
	},
	)

	patchOp := ravendb.NewPatchOperation(claims.ID+"/exercises", nil, patch, nil, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	c.JSON(200, gin.H{})

}

func Update(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)

	var e Exercise
	c.ShouldBind(&e)

	e.ExerciseID = c.Param("exerciseID")

	//create patch
	patch := ravendb.PatchRequest_forScript(
		`this.entries = this.entries.map(entry => {
				if(entry.exerciseID == args.payload.exerciseID){
					entry = args.payload
				}
				return entry
			})`,
	)

	//add values
	patch.SetValues(map[string]interface{}{
		"payload": e,
	},
	)

	//execute
	patchOp := ravendb.NewPatchOperation(claims.ID+"/exercises", nil, patch, nil, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, gin.H{})

}

func AddNew(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)

	var e Exercise
	c.ShouldBind(&e)

	e.ExerciseID = helpers.UUID()

	//create payload
	payload := map[string]interface{}{
		"payload": e,
	}

	//create patches
	patch := ravendb.PatchRequest_forScript(`this.entries.push(args.payload)`)
	patchIfMissing := ravendb.PatchRequest_forScript(`this.entries = [args.payload];this["@metadata"] = {"@collection" : "exercises"};`)

	//add values
	patch.SetValues(payload)
	patchIfMissing.SetValues(payload)

	//execute
	patchOp := ravendb.NewPatchOperation(claims.ID+"/exercises", nil, patch, patchIfMissing, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	c.JSON(200, gin.H{"exerciseID": e.ExerciseID})
	return

}
