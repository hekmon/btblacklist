# BT Blocklist

Scrape existing blocklists and search RIPE for ranges then merge all of them to a single blocklist in [PeerGuardian v2](https://en.wikipedia.org/wiki/PeerGuardian#Binary_formats) format available over HTTP.

## How to use

Install the package, default configuration (`/etc/btblocklist/config.json`) should be good enough to start. Edit it if needed then start the service:

```bash
systemctl start btblocklist.service
```

First start should look like this:

```raw
systemd[1]: Starting btblocklist...
btblocklist[9120]:
btblocklist[9120]:  • BT Blocklist •
btblocklist[9120]:      (￣ヘ￣)
btblocklist[9120]:
btblocklist[9120]:    INFO: [Main] Loading configuration
btblocklist[9120]: WARNING: [Updater] can't load previous state from disk: can't open 'cache.gob.gz' for reading: open cache.gob.gz: no such file or directory
btblocklist[9120]:    INFO: [Main] Starting HTTP server on 127.0.0.1:4776
systemd[1]: Started btblocklist.
btblocklist[9120]:    INFO: [Updater] ripe results changed (14 uniq results): global state will be updated
btblocklist[9120]:    INFO: [Updater] external blocklist 'bluetack_lvl1': data has changed: global state will be updated
btblocklist[9120]:    INFO: [Updater] Merging and compressing all cached results
btblocklist[9120]:    INFO: [Updater] 14 range(s) from RIPE search and 236811 line(s) from 1 external blocklist(s) compressed to 3.74 MiB in 1.836007419s
```

Then on your torrent client, let's say Transmission, configure the blocklist the be scrapped on btblocklist directly:

```json
{
    "blocklist-enabled": true,
    "blocklist-url": "http://127.0.0.1:4776",
}
```
