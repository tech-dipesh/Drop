# How drop works

Notes on the parts that aren't obvious from reading the code, and why they ended up this way.

## Finding devices

The receiver listens on UDP 47470. The sender sends a tiny JSON probe to the broadcast address of every network it's on, once a second, and every device that answers goes into the list.

It uses the subnet broadcast address (like `192.168.1.255`), not multicast and not `255.255.255.255`. That was the VPN problem. With a VPN on, the default route belongs to the VPN, and multicast and the all-ones broadcast tend to follow it into the tunnel, where they vanish. A subnet broadcast matches the local network's own route, which is more specific than anything the VPN adds, so the OS sends it out over Wi-Fi without any special handling.

Interfaces that are down, loopback, point-to-point (most VPN adapters) or can't broadcast are skipped.

Some networks drop broadcast but still pass normal traffic. So every third round the sender also sends the probe straight to each address in its /24. It's capped at /24 on purpose: a campus /16 would mean 65,000 packets a round.

Devices are known by a random id created on first run, never by IP. The address is taken from the reply packet every time, and the reply carries no address of its own, so a laptop can't accidentally advertise its VPN IP. An IP that changes on every boot doesn't matter.

## Pairing

Every transfer runs over TLS 1.3. The receiver creates a fresh Ed25519 certificate in memory each time it starts, nothing is written to disk. The sender doesn't check that certificate against an authority, because there isn't one on a home network. The pairing code does that job instead.

Right after the handshake:

1. The sender picks 32 random bytes and sends only their SHA-256, the commitment.
2. The receiver sends its own 32 random bytes.
3. The sender reveals its bytes and the receiver checks them against the commitment.
4. Both hash the TLS exported key material plus both sets of bytes and take the result mod 1,000,000. That's the code.

Someone in the middle has two TLS sessions, one with each side, and the exported key material differs between them, so the two screens show different codes. Without step 1 they could keep generating keys until six digits happened to line up, which takes seconds. With it, they have to commit before seeing the other side's bytes, so they get one guess in a million and a miss shows up as mismatched codes. Bluetooth pairing uses the same trick.

## Wire format

Everything after discovery is one TCP connection on port 47470. Control messages are a 4-byte big-endian length followed by JSON, capped at 1 MB so a bad peer can't make the receiver allocate gigabytes.

| Step | From | Message |
|---|---|---|
| 1 | sender | commit |
| 2 | receiver | nonce |
| 3 | sender | reveal |
| 4 | sender | offer: protocol version, sender name, files with sizes, speed test flag |
| 5 | receiver | answer: accept or decline |
| 6 | sender | raw bytes of every file, back to back |
| 7 | receiver | result: ok or the error |

There are no separators between files because the sizes are already in the offer. There's no checksum either: TLS rejects any changed or missing byte on its own.

## Saving files

Each file is written to a temporary `name.<random>.part` file and renamed only once every byte has arrived. On a failure or Ctrl+C the part file is deleted. If a name already exists the new file becomes `name (1).ext`, so nothing is ever overwritten. Names containing `/`, `\`, `:`, `..` or an absolute path are rejected before anything touches the disk.

## Speed

Most big files, videos, photos and zips, are already compressed, so drop doesn't compress again. That would only cost CPU.

The copy loop moves files in 4 MB chunks and bumps an atomic counter. The progress bar reads that counter ten times a second rather than redrawing on every write.

Encryption has a cost: before TLS, drop could hand files to the kernel with `sendfile` and never touch the bytes. Now every byte goes through AES-GCM. Go uses the CPU's AES instructions for that, and on a single machine with no network in between, an encrypted speed test ran at about 1 GB/s, far above what Wi-Fi can carry.

`drop -speed` sends 1 GB generated in memory and the receiver throws it away, so it measures the network and encryption without the disk. If the speed test is fast but real files are slow, the disk is the bottleneck.

## Tradeoffs, for now

- Same network only. There's no relay for devices on different networks.
- The receiver handles one transfer at a time.
- No folders and no resuming a half-finished transfer yet.
- The receiver has to be running to be visible.
