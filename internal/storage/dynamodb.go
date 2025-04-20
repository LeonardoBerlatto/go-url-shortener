package storage

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/leonardoberlatto/go-url-shortener/internal/models"
)

const (
	tableName = "UrlMappings"
)

type DynamoDBStorage struct {
	client *dynamodb.Client
}

func NewDynamoDBStorage(endpoint, region, accessKeyID, secretAccessKey string) (*DynamoDBStorage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if ep := endpoint; ep != "" {
			o.BaseEndpoint = aws.String(ep)
		}
	})

	return &DynamoDBStorage{
		client: client,
	}, nil
}

func (d *DynamoDBStorage) Store(ctx context.Context, mapping models.URLMapping) error {
	item, err := attributevalue.MarshalMap(mapping)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(ShortID)"),
	})

	if err != nil {
		if _, ok := err.(*types.ConditionalCheckFailedException); ok {
			return ErrorConflict
		}
		return err
	}

	return nil
}

func (d *DynamoDBStorage) Get(ctx context.Context, shortID string) (models.URLMapping, error) {
	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ShortID": &types.AttributeValueMemberS{Value: shortID},
		},
	})

	if err != nil {
		return models.URLMapping{}, err
	}

	if result.Item == nil {
		return models.URLMapping{}, ErrorNotFound
	}

	var mapping models.URLMapping
	err = attributevalue.UnmarshalMap(result.Item, &mapping)
	if err != nil {
		return models.URLMapping{}, err
	}

	return mapping, nil
}

func (d *DynamoDBStorage) Delete(ctx context.Context, shortID string) error {
	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ShortID": &types.AttributeValueMemberS{Value: shortID},
		},
	})
	return err
}

func (d *DynamoDBStorage) CheckExists(ctx context.Context, shortID string) (bool, error) {
	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ShortID": &types.AttributeValueMemberS{Value: shortID},
		},
		ProjectionExpression: aws.String("ShortID"),
	})

	if err != nil {
		return false, err
	}

	return result.Item != nil, nil
}

func (d *DynamoDBStorage) IncrementHits(ctx context.Context, shortID string) error {
	_, err := d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ShortID": &types.AttributeValueMemberS{Value: shortID},
		},
		UpdateExpression: aws.String("ADD Hits :inc"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	return err
}
