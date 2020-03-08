# Cloud Logger
Send log data to AWS cloud watch log.

## Description
Send log to AWS cloud watch log.
AWS access key and access key secret are load from ~/.aws/credential file or environment variables.

You can send newline separated logs as follows.
```bash
cat log | cloudlogger --group test --stream default

cloud logger --group test --stream --default "Log 1
Log 2
Log 3"
```

* Currently cloudlogger send execution time as timestamp instead of time in log.

## Usage
```
usage: clowdlogger --group=GROUP --stream=STREAM [<flags>] [<log>]

Send log data to AWS CloudWatch.

Flags:
      --help           Show context-sensitive help (also try --help-long and
                       --help-man).
  -g, --group=GROUP    Log group name
  -s, --stream=STREAM  Log stream name
      --version        Show application version.

Args:
  [<log>]  Log text.
```
