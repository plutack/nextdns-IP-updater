# Maintainer: Your Name <your.email@example.com>
pkgname=nextdns-ip
pkgver=1.0.0
pkgrel=1
pkgdesc="NextDNS IP Updater Service"
arch=('x86_64' 'aarch64' 'armv7h')
url="https://github.com/plutack/nextdns-ip"
license=('MIT')
depends=('glibc')
makedepends=('go')
source=("main.go" "go.mod" "nextdns-ip.service" "nextdns-ip.conf")
md5sums=('SKIP' 'SKIP' 'SKIP' 'SKIP')
backup=('etc/nextdns-ip.conf')

build() {
    cd "$srcdir"
    export CGO_ENABLED=0
    go build -trimpath -ldflags "-s -w" -o nextdns-ip main.go
}

package() {
    cd "$srcdir"
    install -Dm755 nextdns-ip "$pkgdir/usr/bin/nextdns-ip"
    install -Dm644 nextdns-ip.service "$pkgdir/usr/lib/systemd/system/nextdns-ip.service"
    # Install the config as a system-wide fallback/example, but don't overwrite if existing
    install -Dm600 nextdns-ip.conf "$pkgdir/etc/nextdns-ip.conf"
}
