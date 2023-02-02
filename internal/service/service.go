package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/izaakdale/lib/publisher"
	"github.com/izaakdale/lib/server"
	"github.com/izaakdale/service-order/internal/datatore"
	"github.com/kelseyhightower/envconfig"
)

type (
	specification struct {
		Host        string `envconfig:"HOST"`
		Port        string `envconfig:"PORT"`
		AwsRegion   string `envconfig:"AWS_REGION" default:"eu-west-2"`
		TableName   string `envconfig:"TABLE_NAME" required:"true"`
		AWSEndpoint string `envconfig:"AWS_ENDPOINT"`
		TopicArn    string `envconfig:"TOPIC_ARN" required:"true"`
		QueueURL    string `envconfig:"QUEUE_URL" required:"true"`
	}

	Service struct {
		Name   string
		Server *http.Server
	}
)

func New(name string) *Service {
	var spec specification
	err := envconfig.Process("", &spec)
	if err != nil {
		panic(err)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), func(o *config.LoadOptions) error {
		o.Region = spec.AwsRegion
		return nil
	})
	if err != nil {
		panic(err)
	}

	datatore.Init(getAwsDynamoClient(cfg, spec.AWSEndpoint), spec.TableName)

	err = publisher.Initialise(cfg, spec.TopicArn, publisher.WithEndpoint(spec.AWSEndpoint))
	if err != nil {
		panic(err)
	}

	srv, err := server.New(
		Router(),
		server.WithHost(spec.Host),
		server.WithPort(spec.Port),
		server.WithTimeouts(time.Second, time.Second))
	if err != nil {
		panic(err)
	}

	return &Service{name, srv}
}

func (s *Service) Run() {
	log.Printf("service %s starting up", s.Name)
	log.Fatal(s.Server.ListenAndServe())
}

// allows use of localstack
func getAwsDynamoClient(cfg aws.Config, endpoint string) *dynamodb.Client {
	if endpoint != "" {
		return dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL(endpoint)))
	}

	return dynamodb.NewFromConfig(cfg)
}
