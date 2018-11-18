/*
 * The primary-profile should be an IAM user that is attached to one or more cross-account roles.
 * This script will make the necessary API calls to obtain keys and a session token for the role name you pass in.
 * To make use of the keys and token you can use eval:
 * eval ./assume_role -role <yourRoleNameHere>
 */
package main

import (
	"flag"
	"fmt"
	"github.com/zieckey/goini"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

func main() {
	primaryProfile := flag.String("p", "ops", "The primary profile (from $HOME/.aws/credentials) used to assume roles")
	credsFile := flag.String("c", os.Getenv("HOME")+"/.aws/credentials", "The location of your AWS Credentials file")
	role := flag.String("r", "", "The role you wish to assume in the sub-account")
	terraform := flag.Bool("t", true, "Output terraform env vars") 
	flag.Parse()
	if *role == "" {
		fmt.Println("role must be defined")
		os.Exit(1)
	}


	// Load credentials and retrieve ARN.
	creds := credentials.NewSharedCredentials(*credsFile, *primaryProfile)
	
	ini := goini.New()
	err := ini.ParseFile(*credsFile)
	if err != nil {
		fmt.Println("Failed to parse ARN from creds file: ", err.Error())
		os.Exit(1)
	}
	roleARN, ok := ini.SectionGet(*role, "role_arn")
	if !ok {
		fmt.Println("Unable to read " + *role + " from the credentials file")
		os.Exit(1)
	}


	// Assume the role.	
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
	})
	if err != nil {
		fmt.Println("Failed to create session: ", err.Error())
		os.Exit(1)
	}

	svc := sts.New(sess)
	stsInput := &sts.AssumeRoleInput {
		DurationSeconds: aws.Int64(3600),
		ExternalId: aws.String(*role),
		RoleArn: aws.String(roleARN),
		RoleSessionName: aws.String(*role),
	}
	resp, err := svc.AssumeRole(stsInput)
	if err != nil {
		fmt.Println("Failed to assume role: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("export AWS_ACCESS_KEY_ID=" + *resp.Credentials.AccessKeyId)
	fmt.Println("export AWS_SECRET_ACCESS_KEY=" + *resp.Credentials.SecretAccessKey)
	fmt.Println("export AWS_SESSION_TOKEN=" + *resp.Credentials.SessionToken)
	if *terraform {
		fmt.Println("export TF_VAR_aws_access_key=" + *resp.Credentials.AccessKeyId)
		fmt.Println("export TF_VAR_aws_secret_key=" + *resp.Credentials.SecretAccessKey)
		fmt.Println("export TF_VAR_sessiontoken=" + *resp.Credentials.SessionToken)
	}
}
