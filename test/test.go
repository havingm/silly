package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"silly/logger"
	"time"
)

func testYaml() {
	type _tHobby struct {
		Hb1     string
		Hb2     string
		HbOther []string
	}
	type _tPerson struct {
		Name  string `yaml:"name"`
		Age   int    `yaml:"age"`
		Email string `yaml:"email"`
		Hobby _tHobby
	}

	p := _tPerson{
		Name:  "Robin",
		Age:   8,
		Email: "my-email",
		Hobby: _tHobby{
			Hb1:     "badminton",
			Hb2:     "tabletennis",
			HbOther: []string{"hb1", "hb2", "hb3"},
		},
	}
	ps := make([]*_tPerson, 0)
	for i := 0; i < 10; i++ {
		ps = append(ps, &p)
	}
	data, err := yaml.Marshal(ps)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("marshaled data: \n", string(data))

	m := make(map[string]string)
	m["st1"] = "student1"
	m["st2"] = "student2"
	m["st3"] = "student3"
	data, err = yaml.Marshal(m)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("marshaled data: \n", string(data))

	type students struct {
		St1 string
		St2 string
		St3 string
	}
	std1 := &students{
		St1: "student1",
		St2: "student2",
		St3: "student3",
	}
	data, err = yaml.Marshal(std1)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("marshaled data: \n", string(data))
}

func testMongo() {
	// 创建一个上下文对象
	ctx := context.Background()
	clientOpts := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetServerSelectionTimeout(3 * time.Second)
	// 创建一个 MongoDB 客户端连接
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		logger.Error(err)
		return
	}

	type Person struct {
		Name string
		Age  int
	}
	p1 := &Person{
		"MyName1",
		29,
	}
	// 获取要操作的集合
	db := client.Database("Person")
	collection := db.Collection("PersonList")
	_, err = collection.InsertOne(ctx, p1)
	if err != nil {
		logger.Error(err)
		return
	}
	var p2 Person
	result1 := collection.FindOne(ctx, bson.M{"name": "MyName"})
	result1.Decode(&p2)
	logger.Info("p2,name: ", p2.Name, " age: ", p2.Age)
}

func main() {
	//testMongo()
	testYaml()
}
