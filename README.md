# fitbolt
Syncs from Fitbit API and stores data into a local [Bolt](https://github.com/boltdb/bolt) database

## Install

    go install github.com/sajal/fitbolt/cmd/fitbolt

## Usage

```
$ fitbolt --help
Usage of fitbolt:
  -creds string
    	Path to store/load fitbit tokens from (default "/home/sajal/.fitsyncgo")
  -db string
    	path to bolt database(gets created if not exists) (default "/home/sajal/fitsyncgo.db")
  -genisis string
    	The date to start fetching data from (default "2016-01-01")
exit status 2
```

[Register](https://dev.fitbit.com/apps/new) an app with Fitbit, make sure it is of type [Personal](https://dev.fitbit.com/docs/basics/#personal)

    FITBIT_CLIENT="<OAuth 2.0 Client ID>" FITBIT_SECRET="<Client Secret>" fitbolt

This should start fetching data from `genisis`. Due to [stringent rate-limits in Fitbit API](https://dev.fitbit.com/docs/basics/#rate-limits), you can only fetch ~30 days of data per hour. So to bootstrap the database, you would need to run fitbolt every hour until it catches up.
