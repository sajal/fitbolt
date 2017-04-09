# fitbolt
Syncs from Fitbit API and stores data into a local [Bolt](https://github.com/boltdb/bolt) database

## Install

    go install github.com/sajal/fitbolt/cmd/fitbolt

## Usage

```
$ fitbolt --help
NAME:
   fitbolt - sync and query fitbit data locally

USAGE:
   fitbolt [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     sync     Sync with fitbit, needs FITBIT_CLIENT and FITBIT_SECRET to be set
     steps    list steps by day
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dbpath value   path to bolt database(gets created if not exists) (default: "/home/sajal/fitsyncgo.db")
   --creds value    path to bolt database(gets created if not exists) (default: "/home/sajal/.fitsyncgo")
   --genesis value  The date to start fetching data from (default: "2016-01-01")
   --help, -h       show help
   --version, -v    print the version
```

### sync

[Register](https://dev.fitbit.com/apps/new) an app with Fitbit, make sure it is of type [Personal](https://dev.fitbit.com/docs/basics/#personal)

    FITBIT_CLIENT="<OAuth 2.0 Client ID>" FITBIT_SECRET="<Client Secret>" fitbolt sync

This should start fetching data from `genesis`. Due to [stringent rate-limits in Fitbit API](https://dev.fitbit.com/docs/basics/#rate-limits), you can only fetch ~30 days of data per hour. So to bootstrap the database, you would need to run fitbolt every hour until it catches up.

### Query sub commands

TODO
