# mumble-baresip
_a proof of concept_

this project starts and configures automatically with a single CLI command both, mumble and baresip, and connect them via JACK.

The aim was to be able to interconnect a 'normal' telefon conference with a mumble room.

## Requirements
- Working go environment
- Running `jackd`
- `mumble` installed
- `baresip` installed

## Usage
```bash
$ git clone git@github.com:alangecker/mumble-baresip.git
$ go build
$ ./mumble-baresip \
  -mumble-host mumble.yourserver.com \
  -mumble-user sip_test \
  -mumble-channel "Raum Zitrone" \
  -sip-host 77.72.174.129 \
  -sip-username "USERNAME" \
  -sip-password "PASSWORD" \
  -call-number "+492114911111" \
  -dtmf "31606#87121#"
```

## Limitations
- quick and dirty - don't expect stability at all!
- first time writing _go_ ever.
- doesn't handle any responses from baresip or mumble (like call failures, disconnects, dialogs in Mumble,...)
- currently you can only run it once because baresip in the latest version 0.6.6 doesn't support different jack names yet, but it will in the next release (see [baresip#1025](https://github.com/baresip/baresip/pull/1025)).
