# Simple app using Prometheus + Grafana monitoring
## How to use? (*Requires Docker*)
> Stupidiest method, maybe later ill do better way

1) `Clone` this repo and `cd` into it
2) Build and start app using docker
```
docker compose up -d --build
```
> Stop app with `docker compose down`
3) Open Grafana `http://localhost:3000`
> admin/admin for auth by default

4) Go to the connections page and add prometheus
> Use `http://prometheus:9090` by default

5) Now create new dashboard and paste this query for example
```
rate(http_requests_total{app="go-web"}[5m])
```
<img width="600" height="350" alt="image" src="https://github.com/user-attachments/assets/1cb898ff-47da-476e-9ab1-1db04957d03f" />

6) Also you can generate simple traffic for endpoints
```bash
for i in {1..50}; do curl http://localhost:8080/; done
for i in {1..50}; do curl http://localhost:8080/healthz; done
```
