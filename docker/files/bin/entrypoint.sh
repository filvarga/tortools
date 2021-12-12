#!/bin/sh

/bin/rm -f /app/session/rtorrent.lock
/bin/rm -f /app/session/rtorrent.sock

/usr/bin/rtorrent
