package events

import (
	"context"
	"fmt"

	"github.com/IgorEulalio/golang-threat-generator/pkg/client"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type IAMEnumerator struct {
	AWSClient *client.AWSClient
}

func (i IAMEnumerator) EnumerateRolesThatCanBeAssumed() error {
	ctx := context.Background()

	listRolesOutput, err := i.AWSClient.IamClient.ListRoles(ctx, &iam.ListRolesInput{})
	if err != nil {
		return err
	}

	roles := listRolesOutput.Roles
	for _, role := range roles {
		// Do something with the role
		fmt.Println(role.RoleName)
	}

	return nil
}
