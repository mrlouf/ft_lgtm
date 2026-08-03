# Structurizr

This folder contains the source code for the Structurizr project, which is a web application that allows users to create and share software architecture diagrams. The source code has been generated with the help of the IA and I have amended it to fit more closely to the project I have built. The diagrams are generated using the Structurizr DSL, which is a simple text-based language for describing software architecture.

You can have a look at the different diagrams on the [Structurizr Playground](https://playground.structurizr.com/) if you'd like: just download the Structurizr DSL file from this repository and import it into the website. You will then be able to view the diagrams and explore the architecture from the perspective of the system, the containers, the Kubernetes components and the backend execution flow.

## The Necessity of Diagrams

As a distributed system, ft_lgtm is a project that assembles multiple services to provide a complete solution. Starting from a simple web application with a client and a server, I have added an IPFS node to store and publish the code snippets, an OpenTelemetry Collector to collect data from the backend, and a monitoring stack (Prometheus, Tempo, Loki and Grafana) to format and visualise the data.

Each service has its own set of rules, ports, protocols and dependencies, and must be able to fit seamlessly into the system to collaborate with the other components. The diagrams are a way to visualise the architecture of the system and understand how the different components interact with each other, using which protocols, to exchange what data.


