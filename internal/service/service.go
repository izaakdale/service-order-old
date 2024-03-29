package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/izaakdale/lib/publisher"
	"github.com/izaakdale/lib/server"
	"github.com/izaakdale/service-order/internal/datastore"
	"github.com/izaakdale/service-order/schema/order"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

var (
	name = "service-order"
	spec specification
)

type (
	specification struct {
		Host        string `envconfig:"HOST"`
		Port        string `envconfig:"PORT"`
		GRPCPort    string `envconfig:"GRPC_PORT"`
		AwsRegion   string `envconfig:"AWS_REGION" default:"eu-west-2"`
		TableName   string `envconfig:"TABLE_NAME" required:"true"`
		AWSEndpoint string `envconfig:"AWS_ENDPOINT"`
		TopicArn    string `envconfig:"TOPIC_ARN" required:"true"`
		QueueURL    string `envconfig:"QUEUE_URL" required:"true"`
	}

	Service struct {
		Name       string
		HttpServer *http.Server
		GrpcServer *grpc.Server
	}

	GServer struct {
		order.OrderServiceServer
	}
)

func New() *Service {
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

	datastore.Init(getAwsDynamoClient(cfg, spec.AWSEndpoint), spec.TableName)

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

	gsrv := grpc.NewServer()
	ls := &GServer{}
	order.RegisterOrderServiceServer(gsrv, ls)

	return &Service{name, srv, gsrv}
}

func (s *Service) Run() {
	log.Printf("service %s starting up", s.Name)
	go s.RunGRPC()
	log.Fatal(s.HttpServer.ListenAndServe())
}

func (s *Service) RunGRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", spec.Host, spec.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen on %v\n", err)
	}
	log.Fatal(s.GrpcServer.Serve(lis))
}

func (g *GServer) GetOrder(ctx context.Context, o *order.OrderRequest) (*order.Order, error) {
	log.Printf("grpc order request: %s", o.Id)
	return datastore.Fetch(o.Id)
}

// allows use of localstack
func getAwsDynamoClient(cfg aws.Config, endpoint string) *dynamodb.Client {
	if endpoint != "" {
		return dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolver(dynamodb.EndpointResolverFromURL(endpoint)))
	}

	return dynamodb.NewFromConfig(cfg)
}
