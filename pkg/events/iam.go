package events

import (
	"context"

	"github.com/IgorEulalio/golang-threat-generator/pkg/client"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type IAMEnumerator struct {
	AWSClient *client.AWSClient
}

func (i IAMEnumerator) EnumerateRolesThatCanBeAssumed() ([]types.Role, error) {
	ctx := context.Background()

	listRolesOutput, err := i.AWSClient.IamClient.ListRoles(ctx, &iam.ListRolesInput{})
	if err != nil {
		return nil, err
	}

	principalArn, err := i.AWSClient.GetPrincipalArn()
	if err != nil {
		return nil, err
	}

	roles := listRolesOutput.Roles

	var rolesThatCanBeAssumed []types.Role
	for _, role := range roles {

		canBe := canBeAssumedByPrincipal(role, principalArn)
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
	}

	return false
}
