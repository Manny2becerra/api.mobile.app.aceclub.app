package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbClient struct {
	Client             *mongo.Client
	Establishments     *mongo.Collection
	EstablishmentItems *mongo.Collection
	IngredientsLookup  *mongo.Collection
	S3BucketName       string
	S3Client           *s3.Client
}

func (db *DbClient) Close() {
	db.Client.Disconnect(context.TODO())
}

func Connect(connectionString string, databaseName string, s3BucketName string, s3Region string, s3AccessKey string, s3AccessSecret string) (c *DbClient, e error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, errors.New("failed to connect to database " + err.Error())
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, errors.New("failed to ping database " + err.Error())
	}

	fmt.Println("Databases connected successfully")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(s3Region), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3AccessKey, s3AccessSecret, "")))

	if err != nil {
		log.Fatalf("unable to load AWS SDK config")
	}

	s3Client := s3.NewFromConfig(cfg)

	return &DbClient{
		Client:             client,
		Establishments:     client.Database(databaseName).Collection("establishments"),
		EstablishmentItems: client.Database(databaseName).Collection("establishmentItems"),
		IngredientsLookup:  client.Database(databaseName).Collection("ingredients-lookup-v3"),
		S3BucketName:       s3BucketName,
		S3Client:           s3Client,
	}, nil

}
