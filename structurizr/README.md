
# Structurizr

This folder contains the source code for the Structurizr project, which is a web application that allows users to create and share software architecture diagrams.

The source code has been generated with the help of the IA and I have amended it to fit more closely to the project I have built. The diagrams are generated using the Structurizr DSL, which is a simple text-based language for describing software architecture.

You can have a look at the different diagrams on the [Structurizr Playground](https://playground.structurizr.com/) if you'd like: just download the Structurizr DSL file from this repository and import it into the website.

You will then be able to view the diagrams and explore the architecture from the perspective of the system, the containers, the Kubernetes components and the backend execution flow.

## Screenshots

Click to enlarge

#### System Context Diagram

<img width="900" height="270" alt="SystemContext-dark" src="https://github.com/user-attachments/assets/e7dd98fc-5941-466f-8c6c-f679e4580440" />

#### Container Diagram

<img width="900" height="500" alt="Containers-dark" src="https://github.com/user-attachments/assets/f987cbae-76c2-4a45-b0df-1edc973fe0f1" />

#### Backend Execution Flow Diagram

<img width="900" height="270" alt="BackendComponents-dark" src="https://github.com/user-attachments/assets/a862becf-3cb5-480f-9ae8-860a2ccc85bf" />

#### Kubernetes Diagram

<img width="900" height="480" alt="KubernetesDeployment-dark" src="https://github.com/user-attachments/assets/91d0c1b7-7d2d-40af-8f13-aed749d6d25e" />

The accuracy of the diagrams is may change over the evolution of the project, in particular when it comes to the Kubernetes cluster.

It is currently depicting the v0.2.0 release.

## The Necessity of Diagrams

As a distributed system, ft_lgtm is a project that assembles multiple services to provide a complete solution.

Starting from a simple web application with a client and a server, I have added an IPFS node to store and publish the code snippets, an OpenTelemetry Collector to collect data from the backend, and a monitoring stack (Prometheus, Tempo, Loki and Grafana) to format and visualise the data.

Each service has its own set of rules, ports, protocols and dependencies, and must be able to fit seamlessly into the system to collaborate with the other components.

The diagrams are a way to visualise the architecture of the system and understand how the different components interact with each other, using which protocols, to exchange what data.

