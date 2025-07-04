import * as cdk from "aws-cdk-lib";
import * as ecr from "aws-cdk-lib/aws-ecr";
import { Construct } from "constructs";

export interface EcrStackProps extends cdk.StackProps {
  repositoryName: string;
}

export class EcrStack extends cdk.Stack {
  readonly repository: ecr.Repository;

  constructor(scope: Construct, id: string, props: EcrStackProps) {
    super(scope, id, props);

    this.repository = new ecr.Repository(this, "ParticlesRepository", {
      repositoryName: props.repositoryName,
      imageTagMutability: ecr.TagMutability.MUTABLE,
      imageScanOnPush: true,
      lifecycleRules: [
        {
          description: "Keep only latest 5 images",
          maxImageCount: 5,
        },
      ],
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    // Output the repository URI
    new cdk.CfnOutput(this, "RepositoryURI", {
      value: this.repository.repositoryUri,
      description: "ECR Repository URI",
    });

    new cdk.CfnOutput(this, "RepositoryName", {
      value: this.repository.repositoryName,
      description: "ECR Repository Name",
    });
  }
}
