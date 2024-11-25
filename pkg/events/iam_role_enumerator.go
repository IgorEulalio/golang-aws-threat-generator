package events

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/IgorEulalio/golang-threat-generator/pkg/client"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type IAMRoleEnumerator struct {
	AWSClient *client.AWSClient
}

func (i IAMRoleEnumerator) EnumerateRolesThatCanBeAssumed() ([]types.Role, error) {
	ctx := context.Background()

	listRolesOutput, err := i.AWSClient.IamClient.ListRoles(ctx, &iam.ListRolesInput{})
	if err != nil {
		return nil, err
	}

	principalArn, err := i.AWSClient.GetPrincipalArn()
	if err != nil {
		return nil, err
	}

	arn, err := transformAssumedRoleToRoleArn(*principalArn)
	if err != nil {
		return nil, err
	}

	roles := listRolesOutput.Roles

	var rolesThatCanBeAssumed []types.Role
	for _, role := range roles {

		canBe := canBeAssumedByPrincipal(role, &arn)
		if !canBe {
			continue
		}
		rolesThatCanBeAssumed = append(rolesThatCanBeAssumed, role)
	}
	return rolesThatCanBeAssumed, nil
}

func canBeAssumedByPrincipal(role types.Role, principalArn *string) bool {
	assumeRolePolicyDocumentObject, err := ParsePolicyFromEncodedString(*role.AssumeRolePolicyDocument)
	if err != nil {
		return false
	}

	if err != nil {
		return false
	}

	for _, statement := range assumeRolePolicyDocumentObject.Statement {
		if statement.Principal.AWS == *principalArn {
			return true
		}
		fmt.Printf("Principal: %v\n Statement: %s\n ###### \n", *principalArn, statement.Principal.AWS)
	}

	return false
}

func transformAssumedRoleToRoleArn(principalArn string) (string, error) {
	if !strings.Contains(principalArn, "assumed-role") {
		return principalArn, nil
	}

	parts := strings.Split(principalArn, ":")
	if len(parts) < 6 {
		return "", errors.New("invalid ARN format")
	}

	accountID := parts[4]
	roleDetails := parts[5]

	roleParts := strings.Split(roleDetails, "/")
	if len(roleParts) < 2 {
		return "", errors.New("invalid assumed-role format")
	}

	roleName := roleParts[1]
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)

	return roleArn, nil
}
