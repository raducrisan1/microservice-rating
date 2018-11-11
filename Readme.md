# The Rating Microservice

## Overview

The Rating microservice fetches its data from StockInfo. It uses gRPC for this purpose and the message syntax is defined in the file `stockinfo.proto`. To generate the proxy and stub code around this, use the protoc tool like below:

```bash
protoc -I=stockinfo --go_out=plugins=grpc:stockinfo stockinfo/stockinfo.proto
```

Note: you need to have protoc-gen-go plugin installed together with protoc tool. For further information, please consult: [GRPC Go Quick Start](https://grpc.io/docs/quickstart/go.html)