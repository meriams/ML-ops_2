### Step 1: Add fluentforward receiver

Add the fluentforward receiver in the receivers section of `config.yaml` file of the **`otecol-contrib`** directory that you created in the previous step

```bash
receivers:
  fluentforward:
    endpoint: 0.0.0.0:24224
```
You can read more about fluentforward receiver [here](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/fluentforwardreceiver).

We have used port `24224` for listening in fluentforward protocol, you can change it to the port you want.

Modify your config.yaml and add the above receiver
```bash
service:
    ....
    logs:
        receivers: [otlp, fluentforward]
        processors: [batch]
        exporters: [otlp]
```

### Step 2: Update the fluentD config file
Add the following to your **fluentD config** to forward the logs to otel collector.
```bash
<match <directive>>
  @type forward
  send_timeout 60s
  recover_wait 10s
  hard_timeout 60s

  <server>
    name myserver1
    host <host>
    port 24224
  </server>
</match>
```
 Replace <directive> with your directive name. 
 Also we are assuming that you are running the fluentD binary on the host. If not, the value of host might change depending on your environment.

- For MacOS - host is `host.docker.internal`
- For other systems - host is IP address of your system