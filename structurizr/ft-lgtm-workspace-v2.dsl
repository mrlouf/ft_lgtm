workspace "ft_lgtm" "A k8s-hosted sandboxed code runner with an LGTM observability stack and IPFS publishing" {

    model {
        user = person "Developer" "Writes and runs code snippets in the browser editor"

        ft_lgtm = softwareSystem "ft_lgtm" "Sandbox IDE with observability and decentralized publishing" {

            frontend = container "Client" "Code editor, run button, output view" "React, TypeScript, CodeMirror, served via Nginx"

            backend = container "Backend" "Compiles and executes snippets in a WASM sandbox" "Go" {
                httpHandler = component "HTTP Handlers" "Routes /api/run, /api/publish, /api/health" "net/http"
                compilerIface = component "Compiler" "Compiles source into a WASM binary" "TinyGo (Go), Javy (JS)"
                executorIface = component "Executor" "Instantiates and runs the WASM module in a sandbox" "wazero"
                publisherIface = component "Publisher" "Publishes source + output as an IPFS directory" "Kubo RPC client"
                observability = component "Observability" "Tracer, Meter, Logger injected across the pipeline" "OpenTelemetry SDK"
            }

            otelCollector = container "OTel Collector" "Single ingestion point for traces, metrics and logs" "OpenTelemetry Collector"

            kubo = container "Kubo" "IPFS node: stores blocks, serves the public gateway, joins the DHT" "Kubo (IPFS)"

            prometheus = container "Prometheus" "Scrapes and stores metrics" "Prometheus"
            loki = container "Loki" "Stores logs, ingested natively via OTLP" "Loki"
            tempo = container "Tempo" "Stores traces, generates span metrics" "Tempo"
            grafana = container "Grafana" "Dashboards correlating metrics, logs and traces" "Grafana"
        }

        ipfsNetwork = softwareSystem "Public IPFS Network" "Public gateways and peers (ipfs.io, dweb.link, ...)" "External"

        # --- relationships ---
        user -> frontend "Writes and runs code, fetches snippets by CID" "HTTPS"
        frontend -> backend "Sends snippet / requests run" "JSON over HTTPS (/api)"

        httpHandler -> compilerIface "Compiles submitted source"
        compilerIface -> executorIface "Passes compiled WASM binary"
        executorIface -> publisherIface "Passes stdout/output on success"
        httpHandler -> observability "Starts request span, records metrics"
        compilerIface -> observability "Emits compile span"
        executorIface -> observability "Emits execute span"
        publisherIface -> observability "Emits publish span"

        publisherIface -> kubo "Adds directory (source.go + output.txt)" "Kubo RPC HTTP API"
        kubo -> ipfsNetwork "Announces CID via DHT, serves blocks" "libp2p / Bitswap"
        user -> ipfsNetwork "Fetches a published snippet from any public gateway" "HTTPS"

        observability -> otelCollector "Exports traces, metrics, logs" "OTLP/gRPC"
        otelCollector -> tempo "Exports traces" "OTLP/gRPC"
        otelCollector -> loki "Exports logs" "OTLP/HTTP"
        otelCollector -> prometheus "Exposes /metrics" "scraped by Prometheus (pull)"
        tempo -> prometheus "Pushes span-metrics (service graph, span-metrics)" "remote_write"

        grafana -> prometheus "Queries metrics" "PromQL"
        grafana -> loki "Queries logs" "LogQL"
        grafana -> tempo "Queries traces" "TraceQL"
        user -> grafana "Views dashboards" "HTTPS"

        deploymentEnvironment "Production" {
        developer = deploymentNode "Developer Machine" "" "Browser / curl" {
        }

        k3d = deploymentNode "k3d Cluster" "Runs on the pop-os host" "k3d" {
            lb = infrastructureNode "k3d LoadBalancer" "Host-mapped ports: 80 (public), 4001 tcp+udp (public, swarm), 127.0.0.1:5001 (loopback only, IPFS API)"

            ns = deploymentNode "Namespace: lgtm" {
                ingress = infrastructureNode "Ingress (Traefik)" "Host-based routing, port 80"

                frontendPod = deploymentNode "Pod: frontend" {
                    frontendInst = containerInstance frontend
                }
                svcFrontend = infrastructureNode "Service: frontend" "ClusterIP — :80"

                backendPod = deploymentNode "Pod: backend" {
                    backendInst = containerInstance backend
                }
                svcBackend = infrastructureNode "Service: backend" "ClusterIP — :4242"

                otelPod = deploymentNode "Pod: otel-collector" {
                    otelInst = containerInstance otelCollector
                }
                svcOtel = infrastructureNode "Service: otel-collector" "ClusterIP — :4317 (OTLP gRPC), :8889 (Prometheus scrape target), :55679 (zPages)"

                kuboPod = deploymentNode "Pod: kubo" {
                    kuboInst = containerInstance kubo
                }
                svcIpfs = infrastructureNode "Service: ipfs" "ClusterIP — :8080 (gateway, routed via Ingress)"
                svcIpfsExt = infrastructureNode "Service: ipfs-external" "LoadBalancer — :4001 tcp+udp (swarm, public), :5001 (API, loopback only)"

                promPod = deploymentNode "Pod: prometheus" {
                    promInst = containerInstance prometheus
                }
                svcProm = infrastructureNode "Service: prometheus" "ClusterIP — :9090"

                lokiPod = deploymentNode "Pod: loki" {
                    lokiInst = containerInstance loki
                }
                svcLoki = infrastructureNode "Service: loki" "ClusterIP — :3100 (HTTP/OTLP query+push), :9096 (gRPC, internal only)"

                tempoPod = deploymentNode "Pod: tempo" {
                    tempoInst = containerInstance tempo
                }
                svcTempo = infrastructureNode "Service: tempo" "ClusterIP — :4317 (OTLP gRPC receiver), :3200 (HTTP query)"

                grafanaPod = deploymentNode "Pod: grafana" {
                    grafanaInst = containerInstance grafana
                }
                svcGrafana = infrastructureNode "Service: grafana" "ClusterIP — :3000"
            }
        }

        developer -> lb "HTTP (app + gateway)" ":80"
        developer -> lb "IPFS API / WebUI (admin only)" "127.0.0.1:5001"
        lb -> ingress "" ":80"
        lb -> svcIpfsExt "swarm + API passthrough" ":4001, :5001"
        ingress -> svcFrontend "host: lgtm.local" ":80"
        ingress -> svcIpfs "host: ipfs.lgtm.local" ":80 -> :8080"
        svcFrontend -> frontendInst "" ":80"
        svcBackend -> backendInst "" ":4242"
        svcOtel -> otelInst "" ":4317 / :8889 / :55679"
        svcIpfs -> kuboInst "" ":8080"
        svcIpfsExt -> kuboInst "" ":4001 / :5001"
        svcProm -> promInst "" ":9090"
        svcLoki -> lokiInst "" ":3100 / :9096"
        svcTempo -> tempoInst "" ":4317 / :3200"
        svcGrafana -> grafanaInst "" ":3000"
        }
    }

    views {
        systemContext ft_lgtm "SystemContext" {
            include *
            autoLayout
        }

        container ft_lgtm "Containers" {
            include *
            autoLayout
        }

        component backend "BackendComponents" {
            include *
            autoLayout
        }

        deployment ft_lgtm "Production" "KubernetesDeployment" {
            include *
            autoLayout
        }

        styles {
            element "Person" {
                shape person
                background #08427b
                color #ffffff
            }
            element "Software System" {
                background #1168bd
                color #ffffff
            }
            element "Container" {
                background #438dd5
                color #ffffff
            }
            element "Component" {
                background #85bbf0
                color #000000
            }
            element "External" {
                background #999999
                color #ffffff
            }
        }
    }
}
