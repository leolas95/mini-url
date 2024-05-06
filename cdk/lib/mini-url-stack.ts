import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda'
import * as apigw from 'aws-cdk-lib/aws-apigateway'
import {CfnOutput} from "aws-cdk-lib";

export class MiniUrlStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const miniUrlLambda = new lambda.Function(this, 'MiniUrlLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      code: lambda.Code.fromAsset('../deployment/deployment.zip'),
      handler: 'bootstrap',
      architecture: lambda.Architecture.ARM_64,
    });

    const api = new apigw.LambdaRestApi(this, 'LeoApi', {
      handler: miniUrlLambda,
    });

    new CfnOutput(this, 'ApiURL', {
      value: api.url,
      description: 'URL of the API gateway',
    });

    new CfnOutput(this, 'Lambda name', {
      value: miniUrlLambda.functionName
    });

    new CfnOutput(this, 'Lambda ARN', {
      value: miniUrlLambda.functionArn
    });
  }
}
