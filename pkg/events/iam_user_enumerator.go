package events

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/IgorEulalio/golang-threat-generator/pkg/client"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type IAMUserEnumerator struct {
	client *client.AWSClient
}

func (i IAMUserEnumerator) EnumerateUserAndPolicy() error {

	ctx := context.Background()

	iamClient := client.GetAWSClient().IamClient
	users, err := iamClient.ListUsers(ctx, &iam.ListUsersInput{})
	if err != nil {
		fmt.Printf("Error listing users: %v", err)
		return err
	}

	policyList := []string{
		"arn:aws:iam::aws:policy/AdministratorAccess",
		"arn:aws:iam::aws:policy/AmazonS3FullAccess",
		"arn:aws:iam::aws:policy/AmazonEC2FullAccess",
	}

	for _, u := range users.Users {
		for _, p := range policyList {
			fmt.Printf("Executing enumeration for user %v and policy %v", *u.UserName, p)
			cmd := exec.Command("/trigger-rules.sh", *u.UserName, p)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error executing command: %v", err)
				return err
			}
			fmt.Printf("Output: %v", string(output))
		}
	}
	return nil
}
