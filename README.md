# Go Hammertime

This is a simple benchmark driver I wrote to test performance of
[Debezium](http://debezium.io/). As such, the code isn't very extensible or
general purpose. Having said that, you should be able to get away by making
required changes in `models` and `constants`. The `models` package houses the
database model and the `constants` package has the queries that are executed.
`generator.go` decides the distribution the queries.

I'm working to re-architecture this to something more general but that isn't
a priority at the moment.

**NOTE: As of now, the `bench.sh` script and `generate.sh` script won't be of
use to people since it's too tied into the database structure and the code that
it is benchmarking.**

## Steps

1. Create a PostgreSQL server (local or hosted).
2. Take a snapshot of an existing database using `pgdump`.
3. Update the path to that dump in `DUMP_FILE` variable in `bench.sh`.
4. Install `docker` and configure rootless docker.
   ```bash
   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
   sudo add-apt-repository \
     "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
     $(lsb_release -cs) \
     stable"
   sudo apt update
   sudo apt install docker-ce
   sudo groupadd docker
   sudo usermod -aG docker $USER
   ```
5. Install and enable `sysstat`.
   ```bash
   sudo apt install sysstat
   sudo vi /etc/default/sysstat
   # Change ENABLED="false" to ENABLED="true"
   sudo systemctl enable sysstat
   sudo systemctl restart sysstat
   ```
6. Install golang 1.11 or later.
7. Open `bench.sh` and ensure that the following variables are set correctly:
  - `DB_HOST`: Hostname of the database server
  - `DB_NAME`: Database to create for benchmarking
  - `DB_USER`: User to use for benchmarking
  - `SLOT_NAME`: Name to give the postgresql replication slot
  - `DUMP_FILE`: Absolute path to the database dump to restore
8. Open `debezium.json` and ensure that the configuration matches your
  PostgreSQL configuration.
9. Enable logical replication on your RDS or PostgreSQL instance and increase
  `wal_keep_segments` to 512 or more. For RDS you can create a parameter group
  and set `rds.logical_replication=1`. For a non-managed postgreSQL instance
  consult the documentation.
10. Start the benchmark process by executing `bash -x bench.sh --cleanup`.
11. This will create something like the following directory structure as it
  goes along.
    ```
    p10
    ├── c1
    │   ├── b1
    │   ├── b100
    │   ├── b1000
    │   └── b10000
    ├── c16
    │   ├── b1
    │   ├── b100
    │   ├── b1000
    │   └── b10000
    ├── c32
    │   ├── b1
    │   ├── b100
    │   ├── b1000
    │   └── b10000
    ├── c4
    │   ├── b1
    │   ├── b100
    │   ├── b1000
    │   └── b10000
    └── c8
        ├── b1
        ├── b100
        ├── b1000
        └── b10000
    ```
12. If the benchmark fails or seems to stuck due to any reason you can simply
  remove the latest created directory and rerun the benchmark. It will continue
  from where it left off.
13. Once the benchmark is done you can generate a markdown report of the runs
  by executing `generate.sh`. Please install `jq` before doing so as it depends
  on `jq`.
