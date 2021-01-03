Note: Outdated public clone of CodeKN/CodeKN. The current version of CodeKN is not publicly available yet.

# CodeKN

KN is a public and open-source project used to present the development process of the real-world project.

The purpose of the project is to build a platform for developers. As a first step, it gathers articles and engineering post URLs from the authority websites and provides a useful user interface for developers.

At the same time, I also want to show the process of creating software projects for beginners. It's planned to include project examples for Go, Python, and front-end technologies like JavaScript and React for web, and Flutter for mobile.

Tech stack: Go-lang, Docker, Kubernetes, MySQL.

In [Wiki](https://github.com/CoderVlogger/codekn/wiki) you can find link to YouTube videos with coding sessions and tutorials.

## Includes

### ProfX

#### Purpose

Back-end service to gather information from different sources and store in easy to use structured way.

#### Current state

Deployed to Kubernetes cluster and uses MySQL database cluster to store the information.

### Flash

#### Purpose:

API for the back-end data.

#### Current state:

Dummy API endpoint which exposes env of the pod.
