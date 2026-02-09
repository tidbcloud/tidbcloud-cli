# Command format

This document describes the command format users need to provide. The following uses the example resource:


```
ticloud serverless example create -c --displayname
ticloud serverless example get -c --example-id
ticloud serverless example list -c --output
ticloud serverless example delete -c --example-id --force
```

Flags can be omitted and inferred from the TiDB Cloud SDK. The following are also accepted:

```
ticloud serverless example create
ticloud serverless example get
ticloud serverless example list
ticloud serverless example delete
```