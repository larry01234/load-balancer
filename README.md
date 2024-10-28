# Go Load Balancer Proof of Concept


This Proof of Concept (PoC) implements a simple load balancer in Go that distributes HTTP requests across multiple backend nodes. Each node has its own rate limits measured in Bytes Per Minute (BPM) and Requests Per Minute (RPM). The load balancer ensures that requests are forwarded to nodes without exceeding their respective rate limits.