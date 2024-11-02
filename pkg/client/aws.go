package client

import (
	"context"
	"sync"

	configuration "github.com/IgorEulalio/golang-threat-generator/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSClient struct {
	cfg       aws.Config
	IamClient *iam.Client
}

var (
	awsClient *AWSClient
	once      sync.Once
)

func Init() error {
	if awsClient != nil {
		return nil
	}

	var initErr error

	once.Do(func() {
		awsConfig := configuration.LoadConfig().AwsConfig
		var cfg aws.Config
		var err error

		if awsConfig.AccessKey != "" && awsConfig.SecretKey != "" && awsConfig.Region != "" {
			cfg, err = config.LoadDefaultConfig(
				context.TODO(),
				config.WithRegion(awsConfig.Region),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsConfig.AccessKey, awsConfig.SecretKey, "")),
			)
		} else {
			cfg, err = config.LoadDefaultConfig(context.TODO())
		}
		if err != nil {
			initErr = err
			return
		}

		if awsConfig.RoleArn != "" {
			stsClient := sts.NewFromConfig(cfg)
			assumeRoleOptions := func(o *stscreds.AssumeRoleOptions) {
				if awsConfig.ExternalID != "" {
					o.ExternalID = aws.String(awsConfig.ExternalID)
				}
			}
			provider := stscreds.NewAssumeRoleProvider(stsClient, awsConfig.RoleArn, assumeRoleOptions)
			cfg.Credentials = aws.NewCredentialsCache(provider)
		}

		// Perform a dry run to validate credentials
		stsClient := sts.NewFromConfig(cfg)
		_, err = stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
		if err != nil {
			initErr = err
			return
		}

		iamClient := iam.NewFromConfig(cfg)
		awsClient = &AWSClient{
			cfg:       cfg,
			IamClient: iamClient,
		}
	})

	return initErr
}

func GetAWSClient() *AWSClient {
	return awsClient
}

func (client AWSClient) GetRegion() string {
	return client.cfg.Region
}
