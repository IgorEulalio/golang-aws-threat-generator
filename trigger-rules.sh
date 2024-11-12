#!/bin/bash

# Check if the correct number of arguments are provided
if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <USER_NAME> <POLICY_ARN>"
  exit 1
fi

# Assign arguments to variables
USER_NAME="$1"
POLICY_ARN="$2"

# Loop through the API calls indefinitely
while true; do
  echo "Listing groups for user $USER_NAME..."
  aws iam list-groups-for-user --user-name "$USER_NAME" --no-cli-pager

  echo "Listing attached user policies for $USER_NAME..."
  aws iam list-attached-user-policies --user-name "$USER_NAME" --no-cli-pager

  echo "Listing inline user policies for $USER_NAME..."
  aws iam list-user-policies --user-name "$USER_NAME" --no-cli-pager

  echo "Getting policy details for $POLICY_ARN..."
  aws iam get-policy --policy-arn "$POLICY_ARN" --no-cli-pager

  echo "Getting policy version for $POLICY_ARN..."
  VERSION_ID=$(aws iam get-policy --policy-arn "$POLICY_ARN" --query 'Policy.DefaultVersionId' --output text)
  aws iam get-policy-version --policy-arn "$POLICY_ARN" --version-id "$VERSION_ID" --no-cli-pager

  echo "Waiting for 1 seconds before the next loop..."
  sleep 1
done

