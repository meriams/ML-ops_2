&nbsp;

Once you are done intrumenting your Java application, you can run it using the below commands
&nbsp;

### Step 1: Run OTel Collector
 Run this command inside the `otelcol-contrib` directory that you created in the install Otel Collector step

```bash
./otelcol-contrib --config ./config.yaml &> otelcol-output.log & echo "$!" > otel-pid
```
&nbsp;

#### (Optional Step): View last 50 lines of `otelcol` logs
```bash
tail -f -n 50 otelcol-output.log
```

#### (Optional Step): Stop `otelcol`
```bash
kill "$(< otel-pid)"
```
&nbsp;

### Step 2: Run your application

```bash
java -javaagent:<path>/opentelemetry-javaagent.jar -jar {{MYAPP}}.jar
```
&nbsp;

<path> - update it to the path where you downloaded the Java JAR agent in previous step

&nbsp;

If you encounter any difficulties, please consult the [troubleshooting section](https://signoz.io/docs/instrumentation/springboot/#troubleshooting-your-installation) for assistance.