# tsbench
A command line tool that can be used to benchmark SELECT query performance across multiple workers/clients against a PostgreSQL instance. The tool takes as its input a CSV file containing parameters for the SQL query, SQL query itself, number of concurrent workers and different database connection related flags. After executing all queries, it outputs a summary containing stats of how many queries were executed, total processing time across all queries, minimum, median, average and maximum query times, and standard deviation.

The input file should be CSV formatted and should contain parameters in the following form:

```
hostname,start_time,end_time
host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
host_000001,2017-01-02 13:02:02,2017-01-02 14:02:02
host_000008,2017-01-02 18:50:28,2017-01-02 19:50:28
...
```

The first line will always be ignored and besides naming columns, defines as well a number of columns each further row is expected to have.

Here is an example of a SQL query that can be benchmarked:

```sql
SELECT time_bucket('1 minutes', ts) bucket, COUNT(*) cnt
FROM cpu_usage
WHERE host = \$1 AND ts >= \$2 AND ts <= \$3
GROUP BY bucket
ORDER BY bucket;
```

Notice the `$1`, `$2` and `$3` parameters in it. To successfully run SQL script with parameters from the CSV file, the number of parameters in SQL script and number of fields in each line of the CSV file must match.

# Installation
```bash
> cd $GOPATH
> go get -u github.com/r3dfor3st/tsbench
> tsbench -p file_with_params.csv
```

From your `$GOPATH` folder you can descent now into `github.com/r3dfor3st/tsbench` and play with the tool a bit more. Running all unit tests is easy:

```bash
> go test
```

# Available options
Running `tsbench` with the usual -h option reveals a number of available arguments:

```
> tsbench -h
Usage of tsbench:
  -database string
        database name (default "homework")
  -host string
        database server host or socket directory (default "localhost")
  -params string
        file containing parameters for the query (default uses input stream)
  -password string
        connect using a specified password (default none)
  -port uint
        database server port number (default 5432)
  -sql string
        SQL query to execute (uses provided params)
  -user string
        connect as specified database user (default "postgres")
  -workers uint
        number of concurrent workers (default 1)
```

# Usage

If no SQL script is provided, a default one will be used (check [cliparams.go](https://github.com/r3dfor3st/tsbench/blob/master/cliparams.go) for the details). If input params argument is omitted, params will be read from the STDIN.

Here's an example output of running the tool:

```
> tsbench -params query_params.csv -workers 6
total db queries: 200
total time: 909.727776ms
min time: 2.350601ms
max time: 20.490282ms
avg time: 4.548638ms
median time: 4.50262ms
standard deviation: 2.558207ms
```

If you prefer, you can feed parameters using the pipe operator:
```bash
cat query_params.csv | tsbench -workers 6
```
