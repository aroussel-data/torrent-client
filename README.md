# torrent-client

Requires a test .torrent file locally to test against: 

e.g. https://cdimage.debian.org/debian-cd/current/amd64/bt-cd/debian-12.11.0-amd64-netinst.iso.torrent

```bash
go run main.go debian-12.torrent
```
BitTorrent protocol spec can be found at: 

https://www.bittorrent.org/beps/bep_0003.html

Acknowledgements:
Thank you to [Jesse Li](https://blog.jse.li/posts/torrent/) for the original guide and implementation. 

