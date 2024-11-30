# Multipathd Exporter for Prometheus

The **Multipathd Exporter** is a Prometheus exporter designed to expose metrics related to multipath devices from `multipathd`, which is the daemon responsible for managing Multipath I/O on Linux systems. This exporter collects data from the multipathd daemon and exposes it in a format that Prometheus can scrape.

## Features

- Exposes metrics related to multipath devices, such as path status, device status, and load balancing information.
- Supports multiple multipath device statistics.
- Can be integrated with Prometheus to monitor the health and performance of multipath devices.
- Allows monitoring of multipath failures, active paths, and path weights.

## Requirements

- **Linux system** with `multipathd` service running.
- **Prometheus** for scraping the exporter metrics.
- **Go** for building the exporter from source (optional).
- **root** or sufficient privileges to access multipathd statistics.

## Installation

You can install the Multipathd Exporter either by downloading pre-built binaries or by building it from the source.

### Build from Source

1. Clone the repository

```bash
git clone https://github.com/wr8fdy/multipath-exporter.git
cd multipathd_exporter
```

2. Build the exporter using Go.

```bash
go build -o multipathd_exporter
```

3. Move the binary to a directory in your system's `PATH` (e.g., `/usr/local/bin`).

```bash
sudo mv multipathd_exporter /usr/local/bin/
```

## Usage

To start the Multipathd Exporter, run the following command:

```bash
multipathd_exporter
```

By default, the exporter will listen on port 9101:

- HTTP server will be available at http://localhost:9101/metrics.

You can change the listening port by using the --web.listen-address flag:

```bash
multipathd_exporter --web.listen-address=":9200"
```

Systemd Service (Optional)

To run the exporter as a systemd service, you can create a service file. Below is an example systemd unit file for the exporter.

Create a file /etc/systemd/system/multipathd_exporter.service:

```
[Unit]
Description=Multipathd Exporter for Prometheus
After=network.target

[Service]
ExecStart=/usr/local/bin/multipathd_exporter
Restart=always
User=nobody
Group=nogroup

[Install]
WantedBy=multi-user.target
```

Reload systemd and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable multipathd_exporter
sudo systemctl start multipathd_exporter
```

## Metrics Exposed

```
multipath_path_info{dev="sdx",group_id="1",host_adapter="x.x.x.x",map_uuid="uuid",target_wwnn="wwnn"} 1
multipath_path_device_state{dev="sdb",group_id="1",map_uuid="uuid",state="running"} 1
multipath_path_check_state{dev="sdb",group_id="1",map_uuid="uuid",state="ready"} 1
multipath_path_state{dev="sdb",group_id="1",map_uuid="uuid",state="active"} 1
multipath_map_info{name="name",paths="2",sysfs="dm-x",uuid="uuid",vendor="Vendor"} 1
multipath_map_faults{uuid="uuid"} 11
multipath_map_state{state="active",uuid="uuid"} 1
multipath_group_state{group_id="1",map_uuid="uuid",state="active"} 1
```
