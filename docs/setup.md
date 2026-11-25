# Setup

## Install

Download the archive for the system from the [releases page](https://github.com/tech-dipesh/drop/releases), unpack it and put `drop` (`drop.exe` on Windows) somewhere on your PATH.

Or with Go 1.22 or newer:

```bash
go install github.com/tech-dipesh/drop@latest
```

## First run

Set the name other devices will see and where incoming files go. Both are saved, so you only do this once:

```bash
drop -name "Dipesh laptop" -dir ~/Downloads/drop
```

The settings live in one small file:

| System | File |
|---|---|
| Windows | `%AppData%\drop\config.json` |
| macOS | `~/Library/Application Support/drop/config.json` |
| Linux and Android | `~/.config/drop/config.json` |

```json
{
  "id": "11166ec485614f6f",
  "name": "Dipesh laptop",
  "dir": "/home/dipesh/Downloads/drop"
}
```

The id is how other devices recognise this one. Delete the file and drop creates a new identity on the next run. To save somewhere else for a single run without changing the setting, set `DROP_DIR`.

## When nobody shows up

Work down this list:

1. Both devices on the same Wi-Fi or the same phone hotspot. College, office and public Wi-Fi often stop devices from seeing each other, and a hotspot gets around that.
2. On Windows, allow drop when the firewall asks, and set the Wi-Fi network to Private in network settings. On a Public network Windows blocks the discovery replies.
3. If you manage a firewall yourself, drop needs UDP and TCP port 47470.
4. With a VPN on, local traffic has to stay outside the tunnel. Cloudflare WARP keeps private address ranges out of the tunnel by default, so it normally just works. If it doesn't, look for its local network access setting. Psiphon on Windows normally runs as a proxy that drop doesn't use, so it doesn't get in the way. For any other VPN, look for a setting like allow LAN or local network sharing.

## Android

drop runs on Android through Termux, using the Linux ARM64 build.

1. Install Termux from F-Droid.
2. On the phone, download `drop_<version>_linux_arm64.tar.gz` from the releases page.
3. In Termux:

```bash
termux-setup-storage
tar xzf ~/storage/downloads/drop_*_linux_arm64.tar.gz drop
chmod +x drop
./drop -name "phone" -dir ~/storage/downloads
```

Keep the screen on while waiting for files. Android drops broadcast packets while the phone sleeps, so the laptop won't see it.

iPhone has no terminal, so it needs a proper app. That's planned after v1.

## Build from source

```bash
git clone https://github.com/tech-dipesh/drop
cd drop
go test ./...
go build
```

Cross-compile for another system, for example Android:

```bash
GOOS=linux GOARCH=arm64 go build -o drop-arm64 .
```

Push a version tag and GitHub Actions builds and publishes every binary:

```bash
git tag v0.1.0
git push origin v0.1.0
```
