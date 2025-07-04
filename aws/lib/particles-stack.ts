import * as cdk from "aws-cdk-lib";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as elbv2 from "aws-cdk-lib/aws-elasticloadbalancingv2";
import * as ecr from "aws-cdk-lib/aws-ecr";
import * as secretsmanager from "aws-cdk-lib/aws-secretsmanager";
import * as iam from "aws-cdk-lib/aws-iam";
import { Construct } from "constructs";
import { CfnParameters } from "./cfn-params";
import { Cluster } from "./cluster";

export interface ParticlesStackProps extends cdk.StackProps {
  stackName?: string;
}

export class ParticlesStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: ParticlesStackProps) {
    super(scope, id, props);

    // Configuration
    const sshPort = 2222;
    const nlbPort = 22;

    // Parameters
    const params = new CfnParameters(this);

    // Reference existing secret
    const sshHostKeySecret = secretsmanager.Secret.fromSecretNameV2(
      this,
      "SshHostKeySecret",
      params.secretName.valueAsString
    );

    // Use default VPC
    const vpc = ec2.Vpc.fromLookup(this, "DefaultVpc", {
      isDefault: true,
    });

    // Security Group for ECS instances
    const ecsSecurityGroup = new ec2.SecurityGroup(this, "EcsSecurityGroup", {
      vpc,
      description: "Security group for ECS instances",
      allowAllOutbound: true,
    });

    ecsSecurityGroup.addIngressRule(
      ec2.Peer.anyIpv4(),
      ec2.Port.tcp(sshPort),
      "Allow SSH particle simulation traffic"
    );

    // ECS Cluster
    const cluster = new Cluster(this, "ParticlesCluster", {
      vpc,
      clusterName: "particles-cluster",
      instanceType: params.instanceType,
      securityGroup: ecsSecurityGroup,
    });

    // Task Definition
    const taskDefinition = new ecs.Ec2TaskDefinition(
      this,
      "ParticlesTaskDefinition",
      {
        networkMode: ecs.NetworkMode.BRIDGE,
      }
    );

    // Add permissions to read secret
    if (taskDefinition.taskRole instanceof iam.Role) {
      taskDefinition.taskRole.addToPolicy(
        new iam.PolicyStatement({
          effect: iam.Effect.ALLOW,
          actions: ["secretsmanager:GetSecretValue"],
          resources: [sshHostKeySecret.secretArn],
        })
      );
    }

    // ECR Repository and Image
    const image = ecs.ContainerImage.fromEcrRepository(
      ecr.Repository.fromRepositoryName(
        this,
        "EcrRepository",
        params.ecrRepositoryName.valueAsString
      ),
      params.imageTag.valueAsString
    );

    // Container Definition
    const container = taskDefinition.addContainer("particles-container", {
      image,
      memoryLimitMiB: 948,
      cpu: 2048,
      command: ["serve", "--port", sshPort.toString()],
      logging: ecs.LogDrivers.awsLogs({
        streamPrefix: "particles",
      }),
      environment: {
        AWS_DEFAULT_REGION: cdk.Aws.REGION,
        PARTICLES__API__SSH_HOST_KEY: "sm:" + params.secretName.valueAsString,
      },
    });

    container.addPortMappings({
      containerPort: sshPort,
      hostPort: sshPort,
      protocol: ecs.Protocol.TCP,
    });

    // ECS Service
    new ecs.Ec2Service(this, "ParticlesService", {
      cluster: cluster.ecsCluster,
      taskDefinition,
      desiredCount: 1,
      serviceName: "particles-service",
      capacityProviderStrategies: [
        {
          capacityProvider: cluster.capacityProvider.capacityProviderName,
          weight: 1,
        },
      ],
      minHealthyPercent: 0,
      maxHealthyPercent: 100,
    });

    // Network Load Balancer
    const nlb = new elbv2.NetworkLoadBalancer(this, "ParticlesNLB", {
      vpc,
      internetFacing: true,
      loadBalancerName: "particles-nlb",
    });

    // Target Group
    const targetGroup = new elbv2.NetworkTargetGroup(
      this,
      "ParticlesTargetGroup",
      {
        port: sshPort,
        protocol: elbv2.Protocol.TCP,
        vpc,
        targetType: elbv2.TargetType.INSTANCE,
        healthCheck: {
          enabled: true,
          protocol: elbv2.Protocol.TCP,
          port: sshPort.toString(),
          healthyThresholdCount: 2,
          unhealthyThresholdCount: 2,
          timeout: cdk.Duration.seconds(10),
          interval: cdk.Duration.seconds(30),
        },
      }
    );

    targetGroup.addTarget(cluster.autoScalingGroup);

    // Listener
    nlb.addListener("ParticlesListener", {
      port: nlbPort,
      protocol: elbv2.Protocol.TCP,
      defaultTargetGroups: [targetGroup],
    });

    // Outputs
    new cdk.CfnOutput(this, "LoadBalancerDNS", {
      value: nlb.loadBalancerDnsName,
      description: "DNS name of the Network Load Balancer",
    });

    new cdk.CfnOutput(this, "SSHCommand", {
      value: `ssh ${nlb.loadBalancerDnsName}`,
      description: "Command to connect to the particle simulation",
    });
  }
}
