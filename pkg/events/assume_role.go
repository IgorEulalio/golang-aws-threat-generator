package events

import (
	"context"
	"encoding/json"
	"io"

	"github.com/IgorEulalio/golang-threat-generator/pkg/client"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AssumeRole struct {
	AWSClient *client.AWSClient
}

func (a AssumeRole) AssumeByArn(roleArn string) (*sts.AssumeRoleOutput, error) {
	ctx := context.Background()

	input := sts.AssumeRoleInput{
		RoleArn:         &roleArn,
		RoleSessionName: aws.String("sysdig-demo-session"),
	}

	role, err := a.AWSClient.StsClient.AssumeRole(ctx, &input)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func DecodeBodyIntoAssumeRole(payload io.Reader) (*string, error) {
	var assumeRole AssumeRoleObject

	d := json.NewDecoder(payload)

	err := d.Decode(&assumeRole)
	if err != nil {
		return nil, err
	}

	return &assumeRole.RoleArn, nil
}

type AssumeRoleObject struct {
	RoleArn string `json:"role_arn"`
}
