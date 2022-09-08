# bandwagonhost cloud exporter
A Prometheus exporter for [bandwagonhost cloud](https://bandwagonhost.com/).

Mertrics api exposing bandwagonhost server information. [mertrics example](./doc/mertrics_example)

## usage
Example
```
wget https://github.com/weiqiang333/bandwagonhost_cloud_exporter/releases/download/v0.1/bandwagonhost_cloud_exporter-v0.1-linux-amd64.tar.gz
mkdir /usr/local/bandwagonhost_cloud_exporter
tar -zxf bandwagonhost_cloud_exporter-linux-amd64.tar.gz -C /usr/local/bandwagonhost_cloud_exporter
chmod +x /usr/local/bandwagonhost_cloud_exporter/bandwagonhost_cloud_exporter
/usr/local/bandwagonhost_cloud_exporter/bandwagonhost_cloud_exporter --config.file /usr/local/bandwagonhost_cloud_exporter/config/bandwagonhost_cloud_exporter.yaml
    # Don't forget to modify your config file /usr/local/bandwagonhost_cloud_exporter/config/bandwagonhost_cloud_exporter.yaml
```
Flags
```
      --config.file string        exporter config file (default "./config/bandwagonhost_cloud_exporter.yaml")
      --exporter.address string   The address on which to expose the web interface and generated Prometheus metrics. (default ":9103")
```

---
## prometheus


## grafana