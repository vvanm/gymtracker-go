package routine

import (
	"github.com/gin-gonic/gin"
	"github.com/ravendb/ravendb-go-client"
	"github.com/vvanm/gymtracker-go/raven"
	"github.com/vvanm/gymtracker-go/util/helpers"
	"github.com/vvanm/gymtracker-go/util/jwt"
)

type Routine struct {
	RoutineID string `json:"routineID,omitempty"`
	//
	Name      string   `json:"name,omitempty"`
	Exercises []string `json:"exercises"`
}

func AddNew(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)

	var r Routine
	c.ShouldBind(&r)

	r.RoutineID = helpers.UUID()
	r.Exercises = []string{}

	//create payload
	payload := map[string]interface{}{
		"payload": r,
	}

	//create patches
	patch := ravendb.PatchRequest_forScript(`this.entries.push(args.payload)`)
	patchIfMissing := ravendb.PatchRequest_forScript(`this.entries = [args.payload];this["@metadata"] = {"@collection" : "routines"};`)

	//add values
	patch.SetValues(payload)
	patchIfMissing.SetValues(payload)

	//execute
	patchOp := ravendb.NewPatchOperation(claims.ID+"/routines", nil, patch, patchIfMissing, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	c.JSON(200, gin.H{"routineID": r.RoutineID})
	return

}
