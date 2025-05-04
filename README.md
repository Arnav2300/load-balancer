# load-balancer

Known issues:
The interval for healthcheck is 5s. Requests will fail if sent before health check if service is down.
