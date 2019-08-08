package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (g game) getArea(c *gin.Context) (area, error) {
	aid, err := g.getAreaID(c)
	if err != nil {
		return area{}, errors.WithMessage(err, "unable to get area")
	}

	a, found := g.grid.area(aid.row, aid.column)
	if !found {
		return area{}, errors.WithMessage(errValidation, "area not found")
	}

	return a, nil
}

func (g game) getAreaID(c *gin.Context) (areaID, error) {
	if g.NumPlayers == 2 {
		obj := struct {
			Row    int `json:"row" binding:"min=1,max=6"`
			Column int `json:"column" binding:"min=1,max=8"`
		}{}

		err := c.ShouldBindJSON(&obj)
		if err != nil {
			return areaID{}, errors.WithMessage(err, "unable to get area")
		}
		return areaID{row: obj.Row, column: obj.Column}, nil
	}

	obj := struct {
		Row    int `json:"row" binding:"min=1,max=7"`
		Column int `json:"column" binding:"min=1,max=8"`
	}{}

	err := c.ShouldBindJSON(&obj)
	if err != nil {
		return areaID{}, errors.WithMessage(err, "unable to get area")
	}

	return areaID{row: obj.Row, column: obj.Column}, nil
}
