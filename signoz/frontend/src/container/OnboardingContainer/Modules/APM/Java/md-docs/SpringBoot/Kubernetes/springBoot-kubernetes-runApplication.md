&nbsp;

Once you are done intrumenting your Java application, you can run it using the below commands


```bash
java -javaagent:<path>/opentelemetry-javaagent.jar -jar {{MYAPP}}.jar
```

<path> - update it to the path where you downloaded the Java JAR agent in previous step

&nbsp;

**Note:**
- In case you're dockerising your application, make sure to dockerise it along with OpenTelemetry instrumentation done in previous step.

&nbsp;

If you encounter any difficulties, please consult the [troubleshooting section](https://signoz.io/docs/instrumentation/springboot/#troubleshooting-your-installation) for assistance.