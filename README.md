# NextDNS IP Updater

A lightweight, cross-platform Go utility designed to automatically link your ISP-issued dynamic IP address to your NextDNS profile.

## Motivation

This project was created to address a specific limitation: Huawei routers often lack native Dynamic DNS (DDNS) support, making it difficult to maintain a consistent connection to NextDNS with a dynamic IP.

I considered using a simple `dnsmasq` configuration, but I wasn't sure how that would play with other tools I use like **HBlock** and **Tailscale**. They both seem to rewrite `dnsmasq` config files whenever they feel like it, so I wanted to avoid potential conflicts. I also figured a standalone Go binary would be cleaner and easier to use cross-platform.

By implementing this as a lightweight Go binary, the updater remains isolated, reliable, and portable across different operating systems and architectures.

## Features

- **Reliable Updates**: Automatically calls the NextDNS link API at a configurable interval.
- **Conflict-Free**: Runs independently of `dnsmasq`, `systemd-resolved`, or other DNS services.
- **Cross-Platform**: Built with Go, making it easy to compile for Linux (x86/ARM), macOS, or Windows.
- **Systemd Integration**: Includes a service file for background operation on Linux.
- **XDG Config Support**: Reads configuration from `~/.config/nextdns-ip/config` (if running as user) or `/etc/nextdns-ip.conf`.
- **Arch Linux Ready**: Includes a `PKGBUILD` for easy installation via `makepkg` or `paru`.

## Installation

### Arch Linux (PKGBUILD)

This project includes a `PKGBUILD` for easy installation on Arch-based systems.

**Important:** Arch Linux packages do **not** enable services automatically. You must run `systemctl enable --now nextdns-ip` after installation.

1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/nextdns-ip.git
    cd nextdns-ip
    ```

2.  **Configure (Optional)**:
    You can modify `nextdns-ip.conf.example` and rename it to `nextdns-ip.conf` before building if you want to bake the config into `/etc`.
    
    *Alternatively*, skip this and create your own user config later at `~/.config/nextdns-ip/config`.

3.  Build and install:
    ```bash
    makepkg -si
    ```

4.  Enable and start the service:
    ```bash
    sudo systemctl enable --now nextdns-ip.service
    ```

### Manual Installation (Generic Linux)

1.  **Build the binary**:
    ```bash
    go build -o nextdns-ip main.go
    ```

2.  **Move to system path**:
    ```bash
    sudo mv nextdns-ip /usr/local/bin/
    ```

3.  **Install Service**:
    Copy `nextdns-ip.service` to `/etc/systemd/system/`.

4.  **Start**:
    ```bash
    sudo systemctl daemon-reload
    sudo systemctl enable --now nextdns-ip.service
    ```

## Configuration

The updater looks for configuration in the following order:

1.  **Environment Variables**: `NEXTDNS_IP` and `UPDATE_INTERVAL` (highest priority).
2.  **User Config**: `~/.config/nextdns-ip/config` (follows XDG spec).
3.  **System Config**: `/etc/nextdns-ip.conf`.

### Configuration File Format

The file is a simple key-value pair, similar to an environment file or `.env`:

```ini
NEXTDNS_IP="https://link-ip.nextdns.io/YOUR_ID/YOUR_KEY"
UPDATE_INTERVAL="10m"
```

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXTDNS_IP` | **Required**. The full Link IP URL from NextDNS. | *None* |
| `UPDATE_INTERVAL` | How often to update the IP (e.g., `5m`, `1h`). | `10m` |

## Todo

- [ ] **Release Workflow**: Implement a GitHub Action (or similar) to automatically build binaries for all supported architectures (x86, ARM, ARM64) upon tagging a release.
- [ ] **Dynamic PKGBUILD**: Update the PKGBUILD to automatically fetch the latest release tag and calculate checksums, making version updates seamless.
