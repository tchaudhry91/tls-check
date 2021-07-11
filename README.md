# TLS-Check

This is a simple utility to get some details from a remote TLS Server.
Useful to check which servers are nearing expiry or are failing hostname verification.

This does NOT verify the certificate chains against CAs.

```bash
This is a simple command line utility to check the TLS information for a particular remote server.
You can check individual sites or batch operations by supplying a CSV

Usage:
  tls-check [url] [flags]
  tls-check [command]

Available Commands:
  batch       Run the check across a list of urls
  completion  generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -h, --help            help for tls-check
  -o, --output string   Output Format (default "json")

Use "tls-check [command] --help" for more information about a command.
```
