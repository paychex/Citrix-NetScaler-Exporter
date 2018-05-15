# Prometheus exporter for Citrix NetScaler
This exporter collects statistics from Citrix NetScaler and makes them available for Prometheus to pull.  As the NetScaler is an appliance it's not recommended to run the exporter directly on it, so it will need to run elsewhere.

## NetScaler configuration
The exporter works via a local NetScaler user account.  It would be preferable to configure a specific user for this which only has permissions to retrieve stats and specific configuration details.

If you lean towards the NetScaler CLI, you want to do something like the following (obviously changing the username as you see fit).

````
# Create a new Command Policy which is only allowed to run the stat command
add system cmdPolicy stat ALLOW (^stat.*|show ns license|show serviceGroup)

# Create a new user.  Disabling externalAuth is important as if it is enabled a user created in AD (or other external source) with the same name could login
add system user stats "password" -externalAuth DISABLED # Change the password to reflect whatever complex password you want

# Bind the local user account to the new Command Policy
bind system user stats stat 100
````

## Usage
You can monitor multiple NetScaler instances by passing in the URL, username, and password as command line flags to the exporter.  If you're running multiple exporters on the same server, you'll also need to change the port that the exporter binds to.

| Flag      | Description                                                                                               | Default Value | Env Name            |
| --------- | --------------------------------------------------------------------------------------------------------- | ------------- | ------------------- |
| url       | Base URL of the NetScaler management interface.  Normally something like https://mynetscaler.internal.com | none          | NETSCALER_URL       |
| username  | Username with which to connect to the NetScaler API                                                       | none          | NETSCALER_USERNAME  |
| password  | Password with which to connect to the NetScaler API                                                       | none          | NETSCALER_PASSWORD  |
| bind_port | Port to bind the exporter endpoint to                                                                     | 9280          | NETSCALER_BIND_PORT |
| multi     | Enable multi query endpoint                                                                               | false         | NETSCALER_MULTI     |


Run the exporter manually using the following command:

````
Citrix-NetScaler-Exporter.exe -url https://mynetscaler.internal.com -username stats -password "my really strong password"
````

This will run the exporter using the default bind port.  If you need to change the port, append the ``-bind_port`` flag to the command.

You may also define environment variables for each of the flags and skip the command line switches

### Running as a service
Ideally you'll run the exporter as a service.  There are many ways to do that, so it's really up to you.  If you're running it on Windows I would recommend [NSSM](https://nssm.cc/).

### Running in multi-query mode
While normally one runs one exporter per device, there are times where running one exporter for multiple Netscaler devices may make sense.  This setup works similar to the [SNMP exporter](https://github.com/prometheus/snmp_exporter).  Note that you will need to configure each Netscaler device to use the same username and password for stats.

When configuring Prometheus to scrape in this manner use the following Prometheus config snippet:

````YAML
scrape_configs:
  - job_name: 'netscaler'
    static_configs:
      - targets:
        - 192.168.1.2  # Netscaler device.
        - 192.168.1.3  # Netscaler device 2
        - 192.168.2.2  # Netscaler device 3, etc
    metrics_path: /query
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9280  # The netscaler exporter's real hostname:port running in "multi-query" mode
  - job_name: 'netscaler-exporter-stats' # gathers the exporter application process stats if you want this sort of information
    static_configs:
      - targets: 127.0.0.1:9280
````

## Exported metrics
### NetScaler

| Metric                                 | Metric Type | Unit    |
| -------------------------------------- | ----------- | ------- |
| CPU usage                              | Gauge       | Percent |
| Memory usage                           | Gauge       | Percent |
| Management CPU usage                   | Gauge       | Percent |
| Packet engine CPU usage                | Gauge       | Percent |
| /flash partition usage                 | Gauge       | Percent |
| /var partition usage                   | Gauge       | Percent |
| MB received per second                 | Gauge       | MB/s    |
| MB sent per second                     | Gauge       | MB/s    |
| HTTP requests per second               | Gauge       | None    |
| HTTP responses per second              | Gauge       | None    |
| Current client connections             | Gauge       | None    |
| Current established client connections | Gauge       | None    |
| Current server connections             | Gauge       | None    |
| Current established server connections | Gauge       | None    |

### Interfaces
For each interface, the following metrics are retrieved.

| Metric                                 | Metric Type | Unit    |
| -------------------------------------- | ----------- | ------- |
| Interface ID                           | N/A         | None    |
| Received bytes per second              | Gauge       | Bytes   |
| Transmitted bytes per second           | Gauge       | Bytes   |
| Received packets per second            | Gauge       | None    |
| Transmitted packets per second         | Gauge       | None    |
| Jumbo packets retrieved per second     | Gauge       | None    |
| Jumbo packets transmitted per second   | Gauge       | None    |
| Error packets received per second      | Gauge       | None    |
| Intrerface alias                       | N/A         | None    |

## Virtual Servers
For each virtual server, the following metrics are retrieved.

| Metric                     | Metric Type | Unit    |
| ---------------------------| ----------- | ------- |
| Name                       | N/A         | None    |
| Waiting requests           | Gauge       | None    |
| Health                     | Gauge       | Percent |
| Inactive services          | Gauge       | None    |
| Active services            | Gauge       | None    |
| Total hits                 | Counter     | None    |
| Hits rate                  | Gauge       | None    |
| Total requests             | Counter     | None    |
| Requests rate              | Gauge       | None    |
| Total responses            | Counter     | None    |
| Responses rate             | Gauge       | None    |
| Total request bytes        | Counter     | Bytes   |
| Request bytes rate         | Gauge       | Bytes/s |
| Total response bytes       | Counter     | Bytes   |
| Response bytes rate        | Gauge       | Bytes/s |
| Current client connections | Gauge       | None    |
| Current server connections | Gauge       | None    |

## Services
For each service, the following metrics are retrieved.

| Metric                         | Metric Type | Unit    |
| -------------------------------| ----------- | ------- |
| Name                           | N/A         | None    |
| Throughput                     | Counter     | MB      |
| Throughput rate                | Gauge       | MB/s    |
| Average time to first byte     | Gauge       | Seconds |
| State                          | Gauge       | None    |
| Total requests                 | Counter     | None    |
| Requests rate                  | Gauge       | None    |
| Total responses                | Counter     | None    |
| Responses rate                 | Gauge       | None    |
| Total request bytes            | Counter     | Bytes   |
| Request bytes rate             | Gauge       | Bytes/s |
| Total response bytes           | Counter     | Bytes   |
| Response bytes rate            | Gauge       | Bytes/s |
| Current client connections     | Gauge       | None    |
| Surge count                    | Gauge       | None    |
| Current server connections     | Gauge       | None    |
| Server established connections | Gauge       | None    |
| Current reuse pool             | Gauge       | None    |
| Max clients                    | Gauge       | None    |
| Current load                   | Gauge       | Percent |
| Service hits                   | Counter     | None    |
| Service hits rate              | Gauge       | None    |
| Active transactions            | Gauge       | None    |

## Service Groups
For each service group member, the following metrics are retrieved.

| Metric                         | Metric Type | Unit    |
| -------------------------------| ----------- | ------- |
| Average time to first byte     | Gauge       | Seconds |
| State                          | Gauge       | None    |
| Total requests                 | Counter     | None    |
| Requests rate                  | Gauge       | None    |
| Total responses                | Counter     | None    |
| Responses rate                 | Gauge       | None    |
| Total request bytes            | Counter     | Bytes   |
| Request bytes rate             | Gauge       | Bytes/s |
| Total response bytes           | Counter     | Bytes   |
| Response bytes rate            | Gauge       | Bytes/s |
| Current client connections     | Gauge       | None    |
| Surge count                    | Gauge       | None    |
| Current server connections     | Gauge       | None    |
| Server established connections | Gauge       | None    |
| Current reuse pool             | Gauge       | None    |
| Max clients                    | Gauge       | None    |

## Licensing

| Metric                         | Metric Type | Unit    |
| -------------------------------| ----------- | ------- |
| Model ID                       | Gauge       | None    |

## Downloading a release
<https://github.com/rokett/Citrix-NetScaler-Exporter/releases>

## Building the executable
All dependencies are version controlled, so building the project is really easy.

1. ``go get github.com/rokett/citrix-netscaler-exporter``.
2. From within the repository directory run ``go build``.
3. Hey presto, you have an executable.

## Dockerfile
A Dockerfile has been setup to create the exporter using golang:alpine3.6

This Dockerfile will create a container that will set the entrypoint as /Citrix-Netscaler-Exporter so you can just pass in the command line options
mentioned about to the container without needing to call the executable
