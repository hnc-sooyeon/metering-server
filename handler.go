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
// @Summary container count
// @Description 특정기간동안 검색된 container
// @name get-string-by-int
// @Accept json
// @Produce json
// @Param container path string true "container name"
// @Param start query string true "start"
// @Param end query string true "end"
// @Router /count/{container} [get]
// @Success 200
// @Failure 400
func GetCountNamespace(c *gin.Context) {
	metricdb := GetCollection()

	//var period Period

	container := c.Param("container")
	start := c.Query("start")
	end := c.Query("end")
	//c.ShouldBindJSON(&period)

	// filter := bson.M{"$and": []bson.M{
	// 	{"value": 99},
	// 	//{"@timestamp": bson.M{"$gt": period.Start, "$lt": period.End}},
	// 	{"@timestamp": bson.M{"$gt": start, "$lt": end}},
	// }}

	filter := bson.M{"$and": []bson.M{
		{"kubernetes.container_name": container},
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
		"container_name": container,
		"count":          count,
		"start":          start,
		"end":            end,
	})
}

/* 아래 항목이 swagger에 의해 문서화 된다. */
// HelloHandler godoc
// @Summary cpu usage
// @Description 현재 cpu 사용량
// @name get-string-by-int
// @Accept json
// @Produce json
// @Param container path string true "container name"
// @Router /cpu/{container} [get]
// @Success 200
// @Failure 400
func GetCpu(c *gin.Context) {
	client := InfluxConn()

	containerName := c.Param("container")
	// start := c.Query("start")
	// end := c.Query("end")

	// Get query client
	queryAPI := client.QueryAPI("primary")

	// query := `from(bucket: "primary")
	// |> range(start: 2022-01-10T07:00:01.396Z, stop: 2022-01-10T07:08:01.396Z)
	// |> filter(fn: (r) => r["container_name"] == "wordpress")
	// |> filter(fn: (r) => r["_field"] == "cpu_usage_nanocores")
	// |> yield(name: "mean")
	// |> mean(column: "_value")`

	query := fmt.Sprintf(`from(bucket: "primary")
	|> range(start: -10m)
	|> filter(fn: (r) => r["_measurement"] == "kubernetes_pod_container")
	|> filter(fn: (r) => r["_field"] == "cpu_usage_nanocores")
	|> filter(fn: (r) => r["container_name"] == "%s")
	|> last()`, containerName)

	//fmt.Print(query)
	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)

	if err != nil {
		panic(err)
	}

	// Iterate over query response
	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		// Access data
		fmt.Printf("value: %v\n", result.Record().Value())
		//fmt.Printf("%s\n", result.Record().Value())
	}

	// check for an error
	if result.Err() != nil {
		fmt.Printf("query parsing error: %\n", result.Err().Error())
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"value": result.Record().Value(),
	// })

	//always close client at the end
	defer client.Close()
}

/* 아래 항목이 swagger에 의해 문서화 된다. */
// HelloHandler godoc
// @Summary memory usage
// @Description 현재 memory 사용량
// @name get-string-by-int
// @Accept json
// @Produce json
// @Param container path string true "container name"
// @Router /memory/{container} [get]
// @Success 200
// @Failure 400
func GetMemory(c *gin.Context) {
	client := InfluxConn()

	container := c.Param("container")
	// start := c.Query("start")
	// end := c.Query("end")

	// Get query client
	queryAPI := client.QueryAPI("primary")

	// query := `from(bucket: "primary")
	// |> range(start: 2022-01-10T07:00:01.396Z, stop: 2022-01-10T07:08:01.396Z)
	// |> filter(fn: (r) => r["container_name"] == "wordpress")
	// |> filter(fn: (r) => r["_field"] == "cpu_usage_nanocores")
	// |> yield(name: "mean")
	// |> mean(column: "_value")`

	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), fmt.Sprintf(`from(bucket: "primary")
	|> range(start: -1m, stop: now())
	|> filter(fn: (r) => r["container_name"] == "%s")
	|> filter(fn: (r) => r["_field"] == "cpu_usage_nanocores")
	|> yield(name: "mean")
	|> last()`, container))

	if err != nil {
		panic(err)
	}

	// Iterate over query response
	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		// Access data
		fmt.Printf("value: %v\n", result.Record().Result())

		// c.JSON(http.StatusOK, gin.H{
		// 	"value": result.Record().Value(),
		// })
	}
	// check for an error
	if result.Err() != nil {
		fmt.Printf("query parsing error: %\n", result.Err().Error())
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"value": result.Record().Value(),
	// })

	//always close client at the end
	defer client.Close()
}
