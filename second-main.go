package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Distance struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type Duration struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type Elements struct {
	Distance Distance `json:"distance"`
	Duration Duration `json:"duration"`
}

type Rows struct {
	Elements []Elements `json:"elements"`
}

type DirectionMatrix struct {
	Origin []string `json:"origin_addresses"`
	Desc   []string `json:"destination_addresses"`
	Rows   []Rows   `json:"rows"`
}

const mongoUrl, dbName, collectionName string = "mongodb://localhost:27017", "distanceSet", "distanceMatrix"

func initMongoConnection() *mongo.Client {
	var client *mongo.Client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(mongoUrl)
	client, _ = mongo.Connect(ctx, clientOptions)

	return client
}

func queryDbAndPrint() {
	client := initMongoConnection()
	collection := client.Database(dbName).Collection(collectionName)
	// D := bson.D{{"_id", "ObjectId('620b87038353634a90babf8f')"}}
	//objectId, err := primitive.ObjectIDFromHex("620c734ce6c3964ae58bb2b3")
	cursor, err := collection.Find(context.TODO(), bson.D{{Key: "lat", Value: 1}, {Key: "long", Value: 2}})
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Result=", results)
	for i := 0; i < len(results); i++ {
		fmt.Println(results[i])
	}
}

func addOnetoMongoDB(directionMatrix DirectionMatrix) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := initMongoConnection()
	collection := client.Database(dbName).Collection(collectionName)
	result, _ := collection.InsertOne(ctx, directionMatrix)
	fmt.Println(result)
}

func praseDirectionMatrixJSONFile(filename string) (DirectionMatrix, error) {
	var directionMatrix DirectionMatrix
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return directionMatrix, err
	}
	defer jsonFile.Close()

	//var directionMatrix DirectionMatrix

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &directionMatrix)
	fmt.Println(directionMatrix)

	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)
	fmt.Println(result["rows"].([]interface{})[0].(map[string]interface{})["elements"].([]interface{})[0].(map[string]interface{})["distance"].(map[string]interface{})["text"])

	return directionMatrix, nil
	// ([0](map[string][]interface{}))["elements"])
}

func getDirectionTest() {
	//url := "https://maps.googleapis.com/maps/api/directions/json?origin=Bangkok&destination=Rayong&key=AIzaSyAWzh2dPlTfViun5CCHavrE3gYv3rBkBGs"
	url := "https://maps.googleapis.com/maps/api/distancematrix/json?origins='Belle%20Condominium%20Bangkok'&destinations='Bangkok%20Cristian%20Collage'&key=AIzaSyAWzh2dPlTfViun5CCHavrE3gYv3rBkBGs"
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(req)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Return Body = ", string(body))
}

func main() {
	// fmt.Println(returnStr())
	// fmt.Println(a111.ReturnStr())
	//getDirectionTest()
	_, error := praseDirectionMatrixJSONFile("mock_direction_matrix.json")
	if error != nil {
		fmt.Println(error)
		return
	}

	//addOnetoMongoDB(directionMatrix)
	queryDbAndPrint()

	elements := Elements{
		Distance: Distance{Text: "t1", Value: 1},
		Duration: Duration{Text: "t2", Value: 2},
	}
	fmt.Println(elements)
}
