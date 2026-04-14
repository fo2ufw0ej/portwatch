# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with default settings:

```bash
portwatch start
```

Specify a custom polling interval and alert on any new or closed ports:

```bash
portwatch start --interval 30s --notify
```

Take a snapshot of the current port state to use as a baseline:

```bash
portwatch snapshot --output baseline.json
```

Watch against an existing baseline:

```bash
portwatch start --baseline baseline.json
```

### Example Output

```
[2024-01-15 10:32:01] INFO  Watching 12 open ports...
[2024-01-15 10:32:31] ALERT New port detected: TCP 0.0.0.0:8080
[2024-01-15 10:33:01] ALERT Port closed: TCP 0.0.0.0:3000
```

---

## Configuration

`portwatch` can be configured via a `portwatch.yaml` file in the working directory:

```yaml
interval: 30s
notify: true
baseline: baseline.json
```

---

## License

MIT © [yourusername](https://github.com/yourusername)