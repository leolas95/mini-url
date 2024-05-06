import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda'
import * as apigw from 'aws-cdk-lib/aws-apigateway'

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
      proxy: false,
    });

    api.root.addResource('/urls', {
      defaultIntegration: new apigw.LambdaIntegration(miniUrlLambda),
    }).addMethod('POST');
  }
}
