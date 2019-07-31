# yadig
Your AWS `dig`. Well, really more like `dig -x` but...

Searches for IPs across multiple configured AWS profiles.

## Install

```bash
go get -u github.com/belminf/yadig
```

## Usage
By default, `yadig` will search all your configured profiles (sourced from `~/.aws/config`) against the default regions:

```bash
$ yadig 10.1.3.51

[personal/us-east2] i-111111111 (bfworks-k8s-control-plane)
```

If you want to customize which profiles and regions `yadig` searches, you could drop a configuration file in `~/.config/yadig/config.yaml`. E.g.:

```yaml
search:
  - alias: personal-alpha
    profile: personal
    region: us-east-1

  - profile: personal
    region: us-east-2

  - alias: work-west
    profile: work
    region: us-west-2
```

## TODO
* Add unit testing
* Add better error handling

## License
[GNU GPLv3](LICENSE.md)
