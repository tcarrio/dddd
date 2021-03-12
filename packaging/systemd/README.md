# systemd

Various files related to systemd service files for `dddd`.

There are two implementations of systemd unit files:

1. native*
2. docker

These both make use of systemd timers for scheduled jobs, referencing `dddd.timer` which will update every 10-15 minutes.

_* Native assumes that the native binary is installed at /usr/bin/dddd_
