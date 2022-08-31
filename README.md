# Space Trouble Booking API

### Quickstart

Prerequisites:
 * Have docker and docker compose installed

```bash
$ docker compose build
$ docker compose up
```

The service is accessible on `http://localhost:8000`

### Flight schedule algorithm

The service shifts the available destinations amongst each launchpad every day. The year day is used to find out which destination is scheduled to which launchpad on a particular day. 
