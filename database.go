package main

import (
	"context"
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

func putSkillOnResume(id string, skill string) error {
	var err error

	ctx := context.Background()

	db = client.Database(DB, nil)
	coll = db.Collection(COLL, nil)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	skill = strings.ToLower(skill)

	filter := bson.M{
		"_id": objId,
		"skills": bson.M{
			"$not": bson.M{
				"$regex": skill, "$options": "i",
			},
		},
	}

	_, err = coll.UpdateOne(ctx, filter, bson.M{"$addToSet": bson.M{"skills": skill}})
	if err == mongo.ErrNoDocuments {
		log.Printf("No document was found with the id %s\n", id)
		return err
	}
	if err != nil {
		return err
	}

	return nil
}
