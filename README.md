# Alertmanager-Webhook-Wechat
Promethues alertmanager webhook wechat 实现

## Usage
编译得到可执行二进制文件，即daemon-systemd下的二进制文件，创建服务器的systemd进程服务。

测试运行可使用：```./alertmanager-webhook-wechat-linux --RobotKey="xxxxxx-xxxxx-xxxxx-xxxxxx-xxxxxxx```

RobotKey即为企业微信生成机器人的key值，也是项目初始化设置的default key。

编辑alertmanager config或者prometheus config

## alertmanager config
```shell
sudo vim /etc/alertmanager/alertmanager.yml
```
```yaml 
global:
  resolve_timeout: 5m #处理超时时间，默认为5min
route:
  group_by: ['alertname']
  group_wait: 3s
  group_interval: 5s
  repeat_interval: 5m
  receiver: 'web.hook'
receivers:
  - name: 'web.hook'
    webhook_configs:
    - url: 'http://localhost:6666/webhook?key=xxxxxx-xxxxx-xxxxx-xxxxxx-xxxxxxx'
```

## prometheus config
```shell
sudo vim /etc/prometheus/alert.rules.yml
```
```yaml
groups:
- name: alert.rules
  rules:
    - alert: HostOutOfDiskSpace
      expr: (node_filesystem_avail{mountpoint="/"}  * 100) / node_filesystem_size{mountpoint="/"} < 50
      for: 1s
      labels:
        severity: warning
     annotations:
       summary: "Host out of disk space (instance {{ $labels.instance }})"
       description: "Disk is almost full (< 50% left)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}"
       wechatRobot: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=9a81393c-a141-4920-9e19-169a445db908"
```

## test
```curl 'http://127.0.0.1:6666/webhook'  -H 'Content-Type: application/json'    -d '
  {
    "receiver": "web.hook",
    "status": "firing",
    "alerts": [
      {
        "status": "firing",
        "labels": {
          "alertname": "HostOutOfDiskSpace",
          "instance": "localhost:9200",
          "severity": "warning"
        },
        "annotations": {
          "info": "Test message,ignore",
          "description": "Disk is almost full (< 50% left)
VALUE = 18.87424590841387
LABELS: map[mountpoint:/ device:/dev/mapper/ubuntu--vg-ubuntu--lv fstype:ext4 instance:localhost:9200 job:node_exporter_metrics]",
          "summary": "Host out of disk space (instance localhost:9200)"
        },
        "startsAt": "2021-04-14T11:33:33.6639785927Z",
        "endsAt": "0001-01-01T00:00:00Z",
        "generatorURL": "warning"
      }
    ],
    "groupLabels": {
      "alertname": "HostOutOfDiskSpace"
    },
    "commonLabels": {
      "alertname": "HostOutOfDiskSpace",
      "instance": "localhost:9200"
    },
    "commonAnnotations": {
      "info": "Test message,ignore",
      "summary": "Host out of disk space (instance localhost:9200)"
    },
    "externalURL": "http://localhost:9093",
    "version": "4",
    "groupKey": "{}/{alertname=~\"^(?:HostOutOfDiskSpace.*)$\"}:{alertname=\"HostOutOfDiskSpace\"}"
  }'
```
