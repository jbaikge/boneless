package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GocmsStackProps struct {
	awscdk.StackProps
}

func NewGocmsStack(scope constructs.Construct, id string, props *GocmsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// table := awsdynamodb.NewTable(stack, jsii.String("GoCMS"), &awsdynamodb.TableProps{
	// 	PartitionKey: &awsdynamodb.Attribute{
	// 		Name: jsii.String("Id"),
	// 		Type: awsdynamodb.AttributeType_STRING,
	// 	},
	// 	BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	// 	TableClass:  awsdynamodb.TableClass_STANDARD,
	// })

	// table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
	// 	IndexName: jsii.String("FrontendIndex"),
	// 	PartitionKey: &awsdynamodb.Attribute{
	// 		Name: jsii.String("Class"),
	// 		Type: awsdynamodb.AttributeType_STRING,
	// 	},
	// 	SortKey: &awsdynamodb.Attribute{
	// 		Name: jsii.String("Slug"),
	// 		Type: awsdynamodb.AttributeType_STRING,
	// 	},
	// })

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get current working directory: %v", err)
	}
	assetDir := filepath.Join(filepath.Dir(currentDir), "assets")

	pingPongFunction := awslambda.NewFunction(stack, jsii.String("PingPong"), &awslambda.FunctionProps{
		Code:    awslambda.NewAssetCode(jsii.String(filepath.Join(assetDir, "ping-pong")), nil),
		Handler: jsii.String("ping-pong"),
		Runtime: awslambda.Runtime_GO_1_X(),
		Timeout: awscdk.Duration_Seconds(jsii.Number(300)),
	})

	// AWS API Gateway V1 Implementation
	api := awsapigateway.NewLambdaRestApi(stack, jsii.String("GoCMSAPI"), &awsapigateway.LambdaRestApiProps{
		Handler: pingPongFunction,
	})

	pingPongResource := api.Root().AddResource(jsii.String("ping-pong"), nil)
	pingPongGetIntegration := awsapigateway.NewLambdaIntegration(
		pingPongFunction,
		&awsapigateway.LambdaIntegrationOptions{},
	)
	pingPongResource.AddMethod(jsii.String("GET"), pingPongGetIntegration, nil)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewGocmsStack(app, "GocmsStack", &GocmsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
