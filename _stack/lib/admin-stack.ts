import * as cdk from 'aws-cdk-lib';
import * as constructs from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as s3Deployment from 'aws-cdk-lib/aws-s3-deployment';
import * as path from 'path';
import * as fs from 'fs-extra';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as cloudfrontOrigins from 'aws-cdk-lib/aws-cloudfront-origins';


export class AdminStack extends cdk.Stack {
  constructor(scope: constructs.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Admin frontend
    // Ref: https://aws-cdk.com/deploying-a-static-website-using-s3-and-cloudfront/
    // Ref: https://github.com/aws-samples/cdk-build-bundle-deploy-example/blob/main/cdk-bundle-static-site-example/lib/static-site-stack.ts
    const adminBucket = new s3.Bucket(this, 'AdminBucket', {
      bucketName: cdk.PhysicalName.GENERATE_IF_NEEDED,
      encryption: s3.BucketEncryption.S3_MANAGED,
      accessControl: s3.BucketAccessControl.PRIVATE,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
    });

    const adminOriginAccessIdentity = new cloudfront.OriginAccessIdentity(this, 'AdminOAI');
    adminBucket.grantRead(adminOriginAccessIdentity);

    const distribution = new cloudfront.Distribution(this, 'AdminDistribution', {
      defaultRootObject: 'index.html',
      defaultBehavior: {
        origin: new cloudfrontOrigins.S3Origin(adminBucket, {originAccessIdentity: adminOriginAccessIdentity}),
      }
    });

    const buildDir = path.resolve(__dirname, '..', '..', '_frontend-admin', 'build');
    new s3Deployment.BucketDeployment(this, 'AdminDeployment', {
      destinationBucket: adminBucket,
      distribution: distribution,
      memoryLimit: 128,
      sources: [
        s3Deployment.Source.asset(buildDir),
      ],
    });

    new cdk.CfnOutput(this, 'CloudFrontUrl', {
      value: distribution.distributionDomainName,
      description: 'CloudFront distribution domain',
    });
  }
}
