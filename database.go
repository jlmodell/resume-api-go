package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB   string = "personal"
	COLL string = "resume"
)

var (
	client *mongo.Client
	db     *mongo.Database
	coll   *mongo.Collection
	uri    string
)

func init() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri = os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}
}

type Resume struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Email          string             `bson:"email" json:"email"`
	Pronouns       string             `bson:"pronouns" json:"pronouns"`
	Certifications []string           `bson:"certifications" json:"certifications"`
	Skills         []string           `bson:"skills" json:"skills"`
	Projects       []string           `bson:"projects" json:"projects"`
	Links          []string           `bson:"links" json:"links"`
}

type FieldList struct {
	Skills         []string `json:"skills,omitempty" bson:"skills,omitempty"`
	Certifications []string `json:"certifications,omitempty" bson:"certifications,omitempty"`
	Projects       []string `json:"projects,omitempty" bson:"projects,omitempty"`
	Links          []string `json:"links,omitempty" bson:"links,omitempty"`
}

func getResumeById(id string) (Resume, error) {
	var err error

	ctx := context.Background()

	db = client.Database(DB, nil)
	coll = db.Collection(COLL, nil)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	var resume Resume
	err = coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&resume)
	if err == mongo.ErrNoDocuments {
		log.Printf("No document was found with the id %s\n", id)
		return resume, err
	}
	if err != nil {
		return resume, err
	}

	return resume, nil
}

func putAdditionalItemInFieldSlice(id string, add string, field string) ([]string, error) {
	var err error
	var fieldList FieldList
	var fieldValue []string

	ctx := context.Background()

	db = client.Database(DB, nil)
	coll = db.Collection(COLL, nil)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fieldValue, err
	}

	add = strings.ToLower(add)

	filter := bson.M{
		"_id": objId,
		field: bson.M{
			"$not": bson.M{
				"$regex": add, "$options": "i",
			},
		},
	}

	update := bson.M{"$addToSet": bson.M{field: add}}

	after := options.After
	projection := bson.M{"_id": 0, field: 1}
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Projection:     &projection,
	}

	err = coll.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&fieldList)

	if err == mongo.ErrNoDocuments {
		return fieldValue, fmt.Errorf("no document was found with the id (%s) without '%s' in '%s'", id, add, field)
	}
	if err != nil {
		return fieldValue, err
	}

	if len(fieldList.Certifications) > 0 {
		fieldValue = fieldList.Certifications
	} else if len(fieldList.Skills) > 0 {
		fieldValue = fieldList.Skills
	} else if len(fieldList.Projects) > 0 {
		fieldValue = fieldList.Projects
	} else if len(fieldList.Links) > 0 {
		fieldValue = fieldList.Links
	}

	return fieldValue, nil
}

func delItemInFieldSlice(id string, item string, field string) ([]string, error) {
	var err error
	var fieldList FieldList
	var fieldValue []string

	ctx := context.Background()

	db = client.Database(DB, nil)
	coll = db.Collection(COLL, nil)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fieldValue, err
	}

	item = strings.ToLower(item)

	filter := bson.M{
		"_id": objId,
		field: bson.M{
			"$in": []string{item},
		},
	}

	update := bson.M{"$pull": bson.M{field: item}}

	after := options.After
	projection := bson.M{"_id": 0, field: 1}
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Projection:     &projection,
	}

	err = coll.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&fieldList)

	if err == mongo.ErrNoDocuments {
		return fieldValue, fmt.Errorf("no document was found with the id (%s) with '%s' in %s", id, item, field)
	}
	if err != nil {
		return fieldValue, err
	}

	if len(fieldList.Certifications) > 0 {
		fieldValue = fieldList.Certifications
	} else if len(fieldList.Skills) > 0 {
		fieldValue = fieldList.Skills
	} else if len(fieldList.Projects) > 0 {
		fieldValue = fieldList.Projects
	} else if len(fieldList.Links) > 0 {
		fieldValue = fieldList.Links
	}

	return fieldValue, nil
}
