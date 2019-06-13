# BT Blocklist

Scrape existing blocklists and search RIPE for ranges then merge all of them to a single blocklist in [PeerGuardian v2](https://en.wikipedia.org/wiki/PeerGuardian#Binary_formats) format available over HTTP.

## How to use

Install the package, default configuration (`/etc/btblocklist/config.json`) should be good enough to start. Then on your torrent client, let's say Transmission, configure the blocklist the be scrapped on btblocklist directly:

```json
{
    "blocklist-enabled": true,
    "blocklist-url": "http://127.0.0.1:4776",
}
```
