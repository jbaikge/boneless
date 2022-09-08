import * as cdk from 'aws-cdk-lib';
import * as constructs from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as cloudfrontOrigins from 'aws-cdk-lib/aws-cloudfront-origins';


export class StaticStack extends cdk.Stack {
  public readonly bucket: s3.Bucket;
  public readonly distribution: cloudfront.Distribution;

  constructor(scope: constructs.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    this.bucket = new s3.Bucket(this, 'StaticBucket', {
      bucketName: cdk.PhysicalName.GENERATE_IF_NEEDED,
      encryption: s3.BucketEncryption.S3_MANAGED,
      accessControl: s3.BucketAccessControl.PUBLIC_READ,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
      cors: [
        {
          allowedMethods: [
            s3.HttpMethods.GET,
            s3.HttpMethods.POST,
            s3.HttpMethods.PUT,
          ],
          allowedOrigins: [
            '*',
          ],
          allowedHeaders: [
            '*',
          ],
        },
      ],
    });

    const originAccessIdentity = new cloudfront.OriginAccessIdentity(this, 'StaticOAI');
    this.bucket.grantRead(originAccessIdentity);

    this.distribution = new cloudfront.Distribution(this, 'StaticDistribution', {
      defaultBehavior: {
        origin: new cloudfrontOrigins.S3Origin(this.bucket, {
          originAccessIdentity: originAccessIdentity,
        }),
      },
    });

    new cdk.CfnOutput(this, 'StaticUrl', {
      value: this.distribution.distributionDomainName,
      description: 'CloudFront distribution domain',
    });
  }
}
