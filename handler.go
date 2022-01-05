package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// type Period struct {
// 	Start string `json:"start"`
// 	End   string `json:"end"`
// }

/* 아래 항목이 swagger에 의해 문서화 된다. */
// HelloHandler godoc
// @Summary namespace count
// @Description 특정기간동안 검색된 namespace
// @name get-string-by-int
// @Accept json
// @Produce json
// @Param namespace path string true "namespace name"
// @Param start query string true "start"
// @Param end query string true "end"
// @Router /{namespace}/count [get]
// @Success 200
// @Failure 400
func GetCountNamespace(c *gin.Context) {
	metricdb := GetCollection()

	//var period Period

	namespace := c.Param("namespace")
	start := c.Query("start")
	end := c.Query("end")
	//c.ShouldBindJSON(&period)

	// filter := bson.M{"$and": []bson.M{
	// 	{"value": 99},
	// 	//{"@timestamp": bson.M{"$gt": period.Start, "$lt": period.End}},
	// 	{"@timestamp": bson.M{"$gt": start, "$lt": end}},
	// }}

	filter := bson.M{"$and": []bson.M{
		{"kubernetes.container_name": namespace},
		{"log": bson.M{"$gt": start, "$lt": end}},
	}}

	//cursor, err := metricdb.Find(context.TODO(), filter)
	count, err := metricdb.CountDocuments(context.TODO(), filter)
	// cursor, _ := uesrsCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println("에러!")
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace_name": namespace,
		"count":          count,
		"start":          start,
		"end":            end,
	})
}
