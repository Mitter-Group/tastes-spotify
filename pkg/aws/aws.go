package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSCredentials struct {
	Region    string
	AccessKey string
	SecretKey string
}

func GetConfig(r *AWSCredentials, local bool) (aws.Config, error) {
	roleArn := os.Getenv("AWS_ROLE_ARN")
	tokenFilePath := os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE")
	if local {
		// localstackPort := "4566"
		// customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		// 	return aws.Endpoint{
		// 		PartitionID:   "aws",
		// 		URL:           fmt.Sprintf("http://localhost:%s", localstackPort),
		// 		SigningRegion: r.Region,
		// 	}, nil
		// })
		// return config.LoadDefaultConfig(context.TODO(),
		// 	config.WithEndpointResolverWithOptions(customResolver))

		return config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(r.Region),
			config.WithSharedConfigProfile("tastes-app"),
		)

	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(r.Region))
	if err != nil {
		panic("failed to load config, " + err.Error())
	}
	client := sts.NewFromConfig(cfg)
	credsCache := aws.NewCredentialsCache(stscreds.NewWebIdentityRoleProvider(
		client,
		roleArn,
		stscreds.IdentityTokenFile(tokenFilePath),
		func(o *stscreds.WebIdentityRoleOptions) {
			o.RoleSessionName = "aws"
		}))
	return config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(credsCache), config.WithRegion(r.Region))

}
