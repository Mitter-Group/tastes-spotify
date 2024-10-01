package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSCredentials struct {
	Region    string
	AccessKey string
	SecretKey string
}

func GetConfig(r *AWSCredentials, env string) (aws.Config, error) {
	roleArn := os.Getenv("AWS_ROLE_ARN")
	tokenFilePath := os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE")
	fmt.Println("aws r:", r)
	fmt.Println("env: ", env)
	switch env {
	case "local":
		return config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(r.Region),
			config.WithSharedConfigProfile("tastes-app"),
		)
	case "dev":
		// TODO: revisar porque no toma los valores desde: r *AWSCredentials
		return config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(r.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			)),
		)
	case "prod":
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
		return config.LoadDefaultConfig(context.TODO(),
			config.WithCredentialsProvider(credsCache),
			config.WithRegion(r.Region))
	default:
		return aws.Config{}, fmt.Errorf("unsupported environment: %s", env)
	}
}
