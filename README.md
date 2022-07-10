# Fathom HTTP

HTTP version of [Fathom](https://github.com/jdart1/Fathom) — Syzygy tablebase probe tool. Server handles following requests:

1. POST /wdl
2. POST /root
3. POST /root-dtz
4. POST /root-wdl

## Build

You need to generate or download pregenerated Syzygy tablebase files.
If you're about to download them, then go to the _tablebases_ directory and follow [these instructions](./tablebases/README.md).

### Go

Just build it: `go mod download && go build`

### Docker

Build an image:

```sh
docker build -t fathom .
```

Run a container:

```sh
docker run \
    -p 8080:80 `# expose port` \
    -v "$(pwd)/tablebases:/app/tablebases" `# mount tablebase directory` \
    -it fathom
```

## Example Requests

### Probe WDL
```sh
curl -X POST 'http://localhost:8080/wdl' \
    --data-binary '{ "position": "4k3/7q/8/8/8/8/B61/4K3 b - - 0 1" }'
```

Response:

```json
{"res":"Win"}
```

### Probe root
```sh
curl -X POST 'http://localhost:8080/root' \
    --data-binary '{ "position": "4k3/7q/8/8/8/8/B61/4K3 b - - 0 1" }'
```

Response (pretty and shortened):

```json
{
    "res": {
        "move": "h7h1",
        "details": [
            {
                "move": "e8d7",
                "wdl": "Win",
                "dtz": 15
            },
            {
                "move": "e8e7",
                "wdl": "Win",
                "dtz": 15
            },
            ...
            {
                "move": "h7g8",
                "wdl": "Draw",
                "dtz": 0
            },
            {
                "move": "h7h8",
                "wdl": "Win",
                "dtz": 5
            }
        ]
    }
}
```

### Probe root DTZ

```sh
curl -X POST 'http://localhost:8080/root-dtz' \
    --data-binary '{ "position": "4k3/7q/8/8/8/8/B61/4K3 b - - 0 1", "useRule50": true }'
```

Response (pretty and shortened):

```json
{
    "res": [
        {
            "move": "e8d7",
            "pv": [],
            "score": 31744,
            "rank": 1000
        },
        {
            "move": "e8e7",
            "pv": [],
            "score": 31744,
            "rank": 1000
        },
        ...
        {
            "move": "h7h8",
            "pv": [],
            "score": 31744,
            "rank": 1000
        }
    ]
}
```

### Probe root WDL

_This is a fallback for the case that some or all DTZ tables are missing._

```sh
curl -X POST 'http://localhost:8080/root-wdl' \
    --data-binary '{ "position": "4k3/7q/8/8/8/8/B61/4K3 b - - 0 1", "useRule50": true }'
```

Response (pretty and shortened):

```json
{
    "res": [
        {
            "move": "e8d7",
            "pv": [],
            "score": 31744,
            "rank": 1000
        },
        {
            "move": "e8e7",
            "pv": [],
            "score": 31744,
            "rank": 1000
        },
        ...
        {
            "move": "h7h8",
            "pv": [],
            "score": 31744,
            "rank": 1000
        }
    ]
}
```

## CLI Arguments

| Parameter       | Type       | Default           | Description |
|----------       |-----       |--------           |------------
| `--allowOrigin` | string     | `"*"`             | Value for HTTP header Access-Control-Allow-Origin
| `--listen`      | string     | `"127.0.0.1:80"`  | HTTP listen [host]:port
| `--maxTime`     | [duration] | `"0s"` (infinite) | Max time limit
| `--poolSize`    | int        | [numcpu]          | Pool size of concurrent Fathom instances
| `--tbDir`       | string     | `"./tablebases"`  | Path to the directory containing Tablebase files

[duration]: https://pkg.go.dev/time#ParseDuration
[numcpu]: https://pkg.go.dev/runtime#NumCPU
