# mcdepth-tools

A lightweight Swiss Army Knife CLI for common dev tasks. Built in Go, runs everywhere.

## Install

```bash
go install github.com/MaplesMcDepth/mcdepth-tools/cmd/mcdepth@latest
```

Or build from source:
```bash
git clone https://github.com/MaplesMcDepth/mcdepth-tools.git
cd mcdepth-tools
go build -o mcdepth ./cmd/mcdepth/
```

## Commands

### `fmt` — Pretty-print JSON
```bash
mcdepth fmt '{"a":1,"b":[1,2,3]}'
echo '{"x":true}' | mcdepth fmt
```

### `b64` — Base64 encode/decode
```bash
mcdepth b64 encode "hello world"
mcdepth b64 decode "aGVsbG8gd29ybGQ="
```

### `jwt` — Decode JWT (no verification)
```bash
mcdepth jwt decode eyJhbGciOiJIUzI1NiIs...
```

### `pass` — Generate password
```bash
mcdepth pass -l 20       # 20 chars, alphanumeric
mcdepth pass -l 20 -s    # include special chars
```

### `url` — URL encode/decode
```bash
mcdepth url encode "hello world!"
mcdepth url decode "hello+world%21"
```

### `time` — Timestamp conversions
```bash
mcdepth time now                    # current time in multiple formats
mcdepth time unix 1778509295        # unix → RFC3339
mcdepth time parse "2026-05-12T00:00:00Z"  # RFC3339 → unix
```

### `uuid` — Generate UUID v4
```bash
mcdepth uuid
# → ab9b5f4a-210a-4a2a-bef2-6156e7ca9873
```

### `hash` — MD5/SHA256
```bash
mcdepth hash md5 "hello"
mcdepth hash sha256 file.txt
```

## License

MIT
