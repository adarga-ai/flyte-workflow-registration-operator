# Flyte workflow registration operator

Kubernetes operator that registers Flyte workflows using a CRD.

Currently, the only way to register Flyte workflows is by using the Flyte CLI. This is a manual process that requires
access to the Flyte Admin API and does not play nicely with how Kubernetes works. This operator aims to automate the
process of registering workflows in Flyte by using a Kubernetes Custom Resource Definition (CRD) 

This operator introduces a `FlyteRegistration` object type to Kubernetes. When a `FlyteRegistration` object is created, 
the operator will register the specified workflow in Flyte. Similarly, when a `FlyteRegistration` object is updated, the
operator will update the specified workflow in Flyte.

The workflows should be packaged using the [Flyte CLI](https://docs.flyte.org/en/latest/api/flytekit/pyflyte.html#pyflyte-package)
and uploaded to one of the supported repositories.

The operator supports two repositories for the workflow package: JFrog Artifactory and OCI repositories. The source is s
specified in the `workflowPackageUri` configuration field.

## Getting started

The operator can be installed in a Kubernetes cluster by using the provided Helm chart (example using an OCI registry):

```sh
helm install flyte-workflow-registration-operator oci://registry-1.docker.io/adarga/flyte-workflow-registration-operator-chart \
  --version 1.0.0 \
  --set controllerManager.manager.env.flyteAdminEndpoint=dns://flyte.mycompany.com \
  --set controllerManager.manager.env.ociRegistry=1234.dkr.ecr.eu-west-1.amazonaws.com \
  --set controllerManager.manager.env.ociAuthStrategy=ecr
```

## Download strategy

To be able to retrieve the workflow package from the specified URI, the operator needs to know how to authenticate with
the source. The `downloadStrategy` field can be set to `jfrog` or `oci`. If set to `jfrog`, the operator will
use the `artifactoryUsername` and `artifactoryPassword` fields to authenticate with the JFrog Artifactory instance. If set to
`oci`, the operator will use the `ociUsername` and `ociPassword` fields to authenticate with the OCI registry.

### OCI authentication strategies

The `ociAuthStrategy` field can be set to `ecr` or `static`.
With the `ecr` strategy the operator will try to use the AWS SDK to authenticate with the registry pulling credentials
from the environment. With the `static` strategy the operator will use the `ociUsername` and `ociPassword` fields to
authenticate with the registry (it does not need to be an ECR registry in that case it can be any OCI compatible registry).

## Flyte credentials

You will also need to provision a `Secret` in your environment named `flyte-credentials` that contains a `clientId` and
a `clientSecret` for the Flyte Admin API. These will be used by the operator to communicate with the Flyte Admin API.

## CRD

This is the definition of the `FlyteRegistration` CRD. You will need to provide one instance of this CRD for each
workflow you want to register.

```yaml
apiVersion: flyte.backend/v1
kind: FlyteRegistration
metadata:
  labels:
    app.kubernetes.io/name: flyteregistration
    app.kubernetes.io/instance: flyteregistration
    app.kubernetes.io/part-of: flyte
  name: my-workflow-registration
spec:
  # Name of the workflow to register
  workflowName: my-workflow
  # Version of the workflow to register
  workflowVersion: 1.2.3
  # Flyte project to register the workflow in
  workflowProject: data-warehouse
  # Flyte domain to register the workflow in
  workflowDomain: development
  # URI of the workflow package
  workflowPackageUri: adarga/data-warehouse-workflows-flyte

```

## Notes

This is not an officially supported Adarga product.