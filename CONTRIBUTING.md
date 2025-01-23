<!-- omit in toc -->
# Contributing to RADAR-Kubernetes

First off, thanks for taking the time to contribute! ❤️

All types of contributions are encouraged and valued. See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles them. Please make sure to read the relevant section before making your contribution. It will make it a lot easier for us maintainers and smooth out the experience for all involved. The community looks forward to your contributions. 🎉

> And if you like the project, but just don't have time to contribute, that's fine. There are other easy ways to support the project and show your appreciation, which we would also be very happy about:
> - Star the project
> - Tweet about it
> - Refer this project in your project's readme
> - Mention the project at local meetups and tell your friends/colleagues

<!-- omit in toc -->
## Table of Contents

- [I Have a Question](#i-have-a-question)
- [I Want To Contribute](#i-want-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)
  - [Improving The Documentation](#improving-the-documentation)
- [Styleguides](#styleguides)
  - [Commit Messages](#commit-messages)
- [Join The Project Team](#join-the-project-team)



## I Have a Question

> If you want to ask a question, we assume that you have read the available [Documentation](https://github.com/RADAR-base/RADAR-Kubernetes/blob/main/README.md).

Before you ask a question, it is best to search for existing [Issues](https://github.com/RADAR-base/RADAR-Kubernetes/issues) and [Discussions](https://github.com/RADAR-base/RADAR-Kubernetes/discussions) that might help you. In case you have found a suitable issue and still need clarification, you can write your question in this issue.

If you then still feel the need to ask a question and need clarification, we recommend the following:

- Open a [Discussion](https://github.com/RADAR-base/RADAR-Kubernetes/discussions/new/choose).
- Provide as much context as you can about what you're running into.
- Provide project and platform versions (nodejs, npm, etc), depending on what seems relevant.

We will then take care of the issue as soon as possible.

## I Want To Contribute

> ### Legal Notice <!-- omit in toc -->
> When contributing to this project, you must agree that you have authored 100% of the content, that you have the necessary rights to the content and that the content you contribute may be provided under the project license.

### Reporting Bugs

<!-- omit in toc -->
#### Before Submitting a Bug Report

A good bug report shouldn't leave others needing to chase you up for more information. Therefore, we ask you to investigate carefully, collect information and describe the issue in detail in your report. Please complete the following steps in advance to help us fix any potential bug as fast as possible.

- Make sure that you are using the latest version.
- Determine if your bug is really a bug and not an error on your side e.g. using incompatible environment components/versions (Make sure that you have read the [documentation](https://github.com/RADAR-base/RADAR-Kubernetes/blob/main/README.md). If you are looking for support, you might want to check [this section](#i-have-a-question)).
- To see if other users have experienced (and potentially already solved) the same issue you are having, check if there is not already a bug report existing for your bug or error in the [bug tracker](https://github.com/RADAR-base/RADAR-Kubernetes/issues?q=label%3Abug).

<!-- omit in toc -->
#### How Do I Submit a Good Bug Report?

> You must never report security related issues, vulnerabilities or bugs including sensitive information to the issue tracker, or elsewhere in public. Instead sensitive bugs must be sent by email to radar-base@kcl.ac.uk and radar-base@thehyve.nl.

We use GitHub issues to track bugs and errors. If you run into an issue with the project:

- Open an [Issue](https://github.com/RADAR-base/RADAR-Kubernetes/issues/new).
- Explain the behavior you would expect and the actual behavior.
- Please provide as much context as possible and describe the *reproduction steps* that someone else can follow to recreate the issue on their own. This usually includes your code. For good bug reports you should isolate the problem and create a reduced test case.
- Provide the information you collected in the previous section.

Once it's filed:

- The project team will label the issue accordingly.
- A team member will try to reproduce the issue with your provided steps. If there are no reproduction steps or no obvious way to reproduce the issue, the team will ask you for those steps and mark the issue as `needs-repro`. Bugs with the `needs-repro` tag will not be addressed until they are reproduced.


### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for RADAR-Kubernetes, **including completely new features and minor improvements to existing functionality**. Following these guidelines will help maintainers and the community to understand your suggestion and find related suggestions.

<!-- omit in toc -->
#### Before Submitting an Enhancement

- Make sure that you are using the latest version.
- Read the [documentation](https://github.com/RADAR-base/RADAR-Kubernetes/blob/main/README.md) carefully and find out if the functionality is already covered, maybe by an individual configuration.
- Perform a [search](https://github.com/RADAR-base/RADAR-Kubernetes/issues) to see if the enhancement has already been suggested. If it has, add a comment to the existing issue instead of opening a new one.
- Find out whether your idea fits with the scope and aims of the project. Keep in mind that we want features that will be useful to the majority of our users and not just a small subset.

<!-- omit in toc -->
#### How Do I Submit a Good Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://github.com/RADAR-base/RADAR-Kubernetes/issues).

- Use a **clear and descriptive title** for the issue to identify the suggestion.
- Provide a **step-by-step description of the suggested enhancement** in as many details as possible.
- **Describe the current behavior** and **explain which behavior you expected to see instead** and why. At this point you can also tell which alternatives do not work for you.
- **Explain why this enhancement would be useful** to most RADAR-Kubernetes users. You may also want to point out the other projects that solved it better and which could serve as inspiration.

### Your First Code Contribution
<!-- TODO
include Setup of env, IDE and typical getting started instructions?

-->
In order to contribute to the repository, first, you need to make sure that you have:
- Read the README.md
- Have a good knowledge of Kubernetes, Helm and Helmfile as listed in [How to get started with tools around RADAR-Kubernetes](https://radar-base.atlassian.net/wiki/spaces/RAD/pages/2731638785/How+to+get+started+with+tools+around+RADAR-Kubernetes)
- Have a working installation of RADAR-Base for testing your changes

Then you can make a new fork or branch and make your changes there and after you have tested it and added neccessarry documentation create a pull request to `dev` branch. We test the installation automatically on supported Kubernetes versions so make sure that the Github actions run successfully. The changes will be reviewed and once merged, they'll be released in the next cycle.
If you're changing an existing code, make sure that it is either backwards compatible or the documentation shows a clear path of applying the changes without breaking the existing installations.


#### Development automation

This repository can be used for development automation for instance on a k3s or k3d (dockerized k3s) cluster. The example below shows how to deploy on a k3d cluster.

1. Install k3d (see [here](https://github.com/k3d-io/k3d#get))
2. Create a k3d cluster that is configured to run RADAR-base

```shell
k3d cluster create my-test-cluster --port '80:80@loadbalancer' --config=.github/ci_config/k3d-config.yaml
```

This example creates a cluster named `my-test-cluster` with a load balancer that forwards local port 80 to the cluster. The
configuration file `.github/ci_config/k3d-config.yaml` is used to configure the cluster. This cluster will be accessible
in _kubectl_ with context name _k3d-my-test-cluster_.

3. Initialize the RADAR-Kubernetes deployment. Run:

```shell
./bin/init
```

4. In file _etc/production.yaml_:

- set _kubeContext_ to _k3d-my-test-cluster_
- set _dev_deployment_ to _true_
- (optional) enable/disable components as needed with the __install_ fields

5. Install RADAR-Kubernetes on the k3d cluster:

```shell
helmfile sync
```

When installation is complete, you can access the applications at `http://localhost`.

#### Adding a new component to RADAR-Kuberentes
In order to add a new component you first need to add its helm chart to [radar-helm-charts)](https://github.com/RADAR-base/radar-helm-charts) repository. Refer to contributing guidelines of that repository for more information. Once the chart has been added you need to:
- Add a helmfile for it in `helmfile.d` directory. The helmfiles are seperated in a modular way in order to avoid having a huge file and also installing certain components in order. Have a look at the current helmfiles and if your component is related to one of them add your component in that file other file create a new file. If your component is a dependency to other components, like Kafka or PostgreSQL prefix the file name with a smaller number so it will be installed first, but if it's a standalone component, the prefix number can be higher.
- Add release to helmfile. Depending on the helm chart this can mostly be copy pasted from other releases and change names to your component. If you've added custom values files in `etc` directory make sure to reference them in the release.
- Add a basic configuration of it to `etc/base.yaml` which should include at least `_install`, `_chart_version` and `_extra_timeout` values. In order to keep the `base.yaml` short, only add configurations that the user will most likely to change during installation.
- If your component is dealing with credentials, the values in the helm charts that refer to that has to be added to `etc/base-secrets.yaml` file.
 - If the credentials isn't something external and can be auto-generated be sure to add it to `bin/generate-secrets`, following examples of the current credentials
- If the user has to input a file to the helm chart, add the relavant key to the `base.yaml.gotmpl` file.
- If the component that you're adding is an external component and you want it to have some default configuration, create a folder with its name in `etc` directory and add the default configuration there in a YAML file and refer to that configuration in the helmfile of the component.

#### Testing the changes
In order to test the changes locally you can use helmfile command to install the component in your cluster. You can make installation faster if you only select your component to install:
```
helmfile apply --file helmfile.d/name-of-the-helmfile.yaml --selector name=name-of-the-component
```
You can also use other the helmfile commands like `helmfile template` and `helmfile diff` to see what is being applied to the cluster.


### Improving The Documentation
<!-- TODO
Updating, improving and correcting the documentation

-->
Feel free to make a PR for any part of the documentation that is missing or isn't clear. If the documentation in question is in the wiki send an email to radar-base@kcl.ac.uk and radar-base@thehyve.nl so we can create an account for you to edit the documentation.


## Join The Project Team
<!-- TODO -->
It's highly recommended to join the [RADAR-Base slack community](https://docs.google.com/forms/d/e/1FAIpQLScKNZ-QonmxNkekDMLLbP-b_IrNHyDRuQValBy1BAsLOjEFpg/viewform) in order to be involved with the community. You can joing #radar-kubernetes to discuss with other developers and attend weekly development meetings.

<!-- omit in toc -->
## Attribution
This guide is based on the **contributing-gen**. [Make your own](https://github.com/bttger/contributing-gen)!
