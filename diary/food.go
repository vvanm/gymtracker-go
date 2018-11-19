package diary

import (
	"github.com/gin-gonic/gin"
	"github.com/ravendb/ravendb-go-client"
	"github.com/vvanm/gymtracker-go/food"
	"github.com/vvanm/gymtracker-go/raven"
	"github.com/vvanm/gymtracker-go/util/jwt"
)

func AddFood(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)
	date := c.Param("date")
	diaryKey := claims.ID + "/diaries/" + date

	var f food.Food
	c.ShouldBind(&f)

	//create patch
	patch := ravendb.PatchRequest_forScript(`this.foods.push(args.payload)`)

	//add values
	patch.SetValues(map[string]interface{}{
		"payload": f,
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

func DeleteFood(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)
	date := c.Param("date")
	diaryKey := claims.ID + "/diaries/" + date

	//create patch
	patch := ravendb.PatchRequest_forScript(` this.foods.splice(args.i,1)`)

	//add values
	patch.SetValues(map[string]interface{}{
		"i": c.Param("i"),
	})

	patchOp := ravendb.NewPatchOperation(diaryKey, nil, patch, nil, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, gin.H{})

}
