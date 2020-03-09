# Cloud Logger
Send log data to AWS cloud watch log.

## Description
Send log to AWS cloud watch log.
AWS access key and access key secret are load from ~/.aws/credential file or environment variables.

if you specified __--log-regexp__ and __--log-format__ options, cloudlogger use time in log string instead of current time.

Also you can send newline separated logs as follows.
```bash
cat log | cloudlogger --group test --stream default

cloud logger --group test --stream --default "Log 1
Log 2
Log 3"
```

## Usage
```
usage: clowdlogger --group=GROUP --stream=STREAM [<flags>] [<log>]

Send log data to AWS CloudWatch.

Flags:
      --help            Show context-sensitive help (also try --help-long and
                        --help-man).
  -g, --group=GROUP     Log group name
  -s, --stream=STREAM   Log stream name
  -t, --time-regexp=""  Time regexp
  -f, --time-format="2006-01-02T15:04:05Z07:00"  
                        Time format
      --version         Show application version.

Args:
  [<log>]  Log text.
```

## Example

```
cloudlogger --group test --stream default --time-regexp "\[(.+?)\] " --time-format "2006-01-02T15:04:05Z07:00" "[2020-03-09T07:18:27.188Z] TEST LOG."
```
