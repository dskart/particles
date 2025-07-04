import { CfnParameter } from "aws-cdk-lib";
import { Construct } from "constructs";

export class CfnParameters {
  readonly ecrRepositoryName: CfnParameter;
  readonly imageTag: CfnParameter;
  readonly instanceType: CfnParameter;
  readonly secretName: CfnParameter;

  constructor(scope: Construct) {
    this.ecrRepositoryName = new CfnParameter(scope, "ECRRepositoryName", {
      type: "String",
      default: "particles",
      description: "The ECR repository name for the particles docker image",
    });
    this.ecrRepositoryName.overrideLogicalId("EcrRepositoryName");

    this.imageTag = new CfnParameter(scope, "ImageTag", {
      type: "String",
      default: "latest",
      description: "Image tag for the particles docker image",
    });
    this.imageTag.overrideLogicalId("ImageTag");

    this.instanceType = new CfnParameter(scope, "InstanceType", {
      type: "String",
      default: "t3a.micro",
      description: "The instance type for the particles ECS tasks",
    });
    this.instanceType.overrideLogicalId("InstanceType");

    this.secretName = new CfnParameter(scope, "SecretName", {
      type: "String",
      default: "particles/ssh-host-key",
      description:
        "The name of the AWS Secrets Manager secret containing the SSH host key",
    });
    this.secretName.overrideLogicalId("SecretName");
  }
}
