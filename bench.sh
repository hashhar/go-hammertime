#!/bin/bash

set -euo pipefail

DB_HOST='localhost'
DB_NAME='godamqadb_benchmark'
DB_USER='gmdevqadb'
SLOT_NAME='benchmark'
DUMP_FILE="$HOME/code/data/godamqadb/dump.sql"

cleanup() {
    echo "Cleaning up"
    pkill -f 'java -cp ../oms-transformer/transformer/target/transformer-1.0-SNAPSHOT-shaded.jar com.delhivery.dwh.topologies.ItemTopology' || true
    curl -X DELETE localhost:8083/connectors/dwh_connector || true
    sleep 5
    containers=$(docker ps -q -f name=zookeeper -f name=kafka -f name=connect || true)
    if [[ -n "$containers" ]]; then
        docker stop $containers
        docker rm $containers
    fi

    rm -rf /tmp/kafka-streams
    psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
        -c 'SELECT pg_drop_replication_slot('"'${SLOT_NAME}'"');' || true # The true command is here since drop_replication_slot
    # will exit with an error if the replication slot does not exist and we will exit the script since it is executed with set -e
    psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
        -c 'DELETE FROM ims_item; TRUNCATE ims_item;' || true
    pkill 'sar' || true
    echo "Cleanup done"
}

# Execute the cleanup function on exit
trap cleanup EXIT
# Execute read after each command when set -x (bash -x) is used
#trap read debug

container_start() {
    docker run -d --name zookeeper --rm \
        -p 2181:2181 -p 2888:2888 -p 3888:3888 \
        debezium/zookeeper:0.9
    docker run -d --name kafka --rm \
        -e KAFKA_MESSAGE_MAX_BYTES=500000000 \
        -p 9092:9092 \
        --link zookeeper:zookeeper \
        debezium/kafka:0.9
    docker run -d --name connect --rm \
        -p 8083:8083 \
        -e GROUP_ID=1 \
        -e CONFIG_STORAGE_TOPIC=debezium-connect-configs \
        -e OFFSET_STORAGE_TOPIC=debezium-connect-offsets \
        -e STATUS_STORAGE_TOPIC=debezium-connect-status \
        -e CONNECT_STATUS_STORAGE_TOPIC=debezium-connect-status \
        -e CONNECT_KEY_CONVERTER_SCHEMAS_ENABLE=false \
        -e CONNECT_VALUE_CONVERTER_SCHEMAS_ENABLE=false \
        --link zookeeper:zookeeper \
        --link kafka:kafka \
        debezium/connect:0.9

    # Wait for some time to let Zookeeper refresh metadata so that broker can be found
    sleep 10
    docker exec -it kafka bin/kafka-topics.sh --create --zookeeper zookeeper:2181 \
        --replication-factor 1 \
        --partitions 1 \
        --config cleanup.policy=compact \
        --topic dwh_connector.public.ims_item
    docker exec -it kafka bin/kafka-topics.sh --create --zookeeper zookeeper:2181 \
        --replication-factor 1 \
        --partitions 1 \
        --config cleanup.policy=compact \
        --topic dwh_connector.public.fclkp_managementlookup
    docker exec -it kafka bin/kafka-topics.sh --create --zookeeper zookeeper:2181 \
        --replication-factor 1 \
        --partitions 1 \
        --config cleanup.policy=compact \
        --topic dwh_transformer.metrics.shipped_item_volume
}

postgres_setup() {
    psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
        -c 'SELECT pg_drop_replication_slot('"'${SLOT_NAME}'"');' || true # The true command is here since drop_replication_slot
    # will exit with an error if the replication slot does not exist and we will exit the script since it is executed with set -e
    if [[ "$1" == "--cleanup" ]]; then
        psql -h "${DB_HOST}" -U "${DB_USER}" -d godamqadb \
            -c "DROP DATABASE IF EXISTS \"${DB_NAME}\";"
        psql -h "${DB_HOST}" -U "${DB_USER}" -d godamqadb \
            -c "CREATE DATABASE \"${DB_NAME}\";"
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -f "${DUMP_FILE}"
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "D5bdb4911801895a2964bb01af245deb";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "D66ce7eb7d4d4d03844b1e58d978bee6";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "e623a161135873cd90fa05e0462b1fbe";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "ims__mgmnt_lkp_id_50066a5f171b91ff_fk_fclkp_managementlookup_id";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "ims_item_last_serial_id_f6e69ee03720e93_fk_ims_item_id";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "ims_item_layout_location_id_5acfcfd7a672e16c_fk_Location_name";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "ims_item_lot_id_6b87b25d9f6058f5_fk_lottracking_lot_sk";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "ims_item_order_line_id_216b2b32aaa5c4f8_fk_oms_orderline_sk";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "ims_item_picklist_id_2ad4d503d7d0d318_fk_outbound_picklist_id";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "ims_item_po_id_10403d6908305656_fk_po_procurementorder_id";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "ims_item_rto_po_id_11b6ffda792752a2_fk_po_procurementorder_id";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "ims_item_rtv_id_id_782c0054474564b_fk_outbound_rtvdispatch_id";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE ims_item DROP CONSTRAINT "ims_item_updated_by_id_5bca299ba2aa786e_fk_auth_user_id";'
        psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
            -c 'ALTER TABLE audit_audithistory DROP CONSTRAINT "audit_audithistory_serial_17e3aa2333593c90_fk_ims_item_serial";'
    fi
    psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
        -c 'DELETE FROM ims_item; TRUNCATE ims_item;'
}

# One time PG setup
postgres_setup "${1-}"
pushd ../oms-transformer \
        && mvn package -DskipTests
popd

# Debezium Measure Throughput #######################################################################################################
# Debezium Poll Intervals
poll_interval_list=( 10 25 50 100 150 200 500 1000 )
# Number of database connections
concurrency_list=( 1 4 8 16 32 )
# Batchsize for item creation
batchsize_list=( 1 100 1000 10000 )
for poll_interval in "${poll_interval_list[@]}"; do
    for concurrency in "${concurrency_list[@]}"; do
        for batchsize in "${batchsize_list[@]}"; do
            echo "TEST RUN WITH poll_interval=${poll_interval} concurrency=${concurrency} batchsize=${batchsize}"
            test_id="p${poll_interval}-c${concurrency}-b${batchsize}"
            test_dir="${test_id//-/\/}"
            if [[ -d "${test_dir}" ]]; then
                    echo "SKIPPING SINCE DONE"
                    continue
            fi
            mkdir -p "${test_dir}"
            # Start monitoring
            sar -o sar.bin -BbdpH -n DEV -r ALL -u ALL 1 >/dev/null 2>&1 &
            pid_sar=$!
            # Number of items
            # con * batch   rate    10k     50k     100k    200k
            # 1     1       1       3m30s
            # 4     1       4       1m7s

            # 8     1       8               3m12s
            # 16    1       16              1m58s
            # 32    1       32              1m29s

            # 1     100     100                     2m45s
            # 4     100     400                     1m43s
            # 8     100     800                     1m12s
            # 1     1000    1000                    1m17s

            # 16    100     1600                            2m32s
            # 32    100     3200                            2m17s
            # 4     1000    4000                            1m59s
            # 8     1000    8000                            2m4s
            # 1     10000   10000                           2m32s
            # 16    1000    16000                           1m34s
            # 32    1000    32000                           2m18s
            # 4     10000   40000                           2m0s
            # 8     10000   80000                           3m23s
            # 16    10000   160000                          3m7s
            # 32    10000   320000                          2m26s
            if (( (concurrency * batchsize) < 5 )); then
                items=10000
            elif (( (concurrency * batchsize) < 100 )); then
                items=50000
            elif (( (concurrency * batchsize) < 1500 )); then
                items=100000
            else
                items=200000
            fi
            # Start Docker containers
            container_start

            # Edit Debezium configuration for selected run and deploy
            sed -i -e 's/\("poll.interval.ms":\).*$/\1 '"${poll_interval}"'/g' debezium.json \
                && sed -i -e 's/\("slot.name":\).*$/\1 "'"${SLOT_NAME}"'",/g' debezium.json \
                && sed -i -e 's/\("database.hostname":\).*$/\1 "'"${DB_HOST}"'",/g' debezium.json \
                && curl -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d @debezium.json
            # Wait for some time to allow Debezium to start fully
            sleep 10

            # Generate load using our loadgen
            go run cmd/load/main.go -host "${DB_HOST}" -database "${DB_NAME}" -concurrency "${concurrency}" -batch-size "${batchsize}" -items "${items}" -progress 30 | tee "${test_dir}/load.${test_id}.log"

            # After our load has finished dump the data from topics into files and
            # then get the needed metrics for PG--->Debezium
            # Get data from topic and exit if no data seen for a minute
            docker exec -it kafka bin/kafka-console-consumer.sh --bootstrap-server 0.0.0.0:9092 \
                    --from-beginning \
                    --timeout-ms 60000 \
                    --topic dwh_connector.public.ims_item > "${test_dir}/ims_item.${test_id}.jsonl"

            # Wait for the WALSender to catch up
            while (( $(cat "${test_dir}/ims_item.${test_id}.jsonl" | wc -l) < items )); do
                sleep 10
            done
            # Start the Kafka Streams application in background and kill when
            # expected number of records have been pushed
            java -cp ../oms-transformer/transformer/target/transformer-1.0-SNAPSHOT-shaded.jar com.delhivery.dwh.topologies.ItemTopology --reset &
            pid_streams=$!

            sleep 10

            # Get the data from Kafka Streams output topics into files and then get
            # the needed metrics for Kafka Streams
            # Wait until the expected number of messages have been generated
            expected_count=$(psql -A -t -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" -c "SELECT COUNT(1) FROM ims_item WHERE status='SHP';")
            docker exec -it kafka bin/kafka-console-consumer.sh --bootstrap-server 0.0.0.0:9092 \
                    --from-beginning \
                    --max-messages "${expected_count}" \
                    --topic dwh_transformer.metrics.shipped_item_volume > "${test_dir}/shipped_item_volume.${test_id}.jsonl"
            kill "${pid_streams}"
            # Dump current composition of the ims_item table
            psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" \
                -c 'SELECT status, (COUNT(status)/(SELECT COUNT(1) FROM ims_item)::FLOAT)*100 FROM ims_item GROUP BY status ORDER BY 1 ASC;' | tee "${test_dir}/table_composition.${test_id}.out"
            echo "TEST RUN DONE"
            echo "STARTING ANALYSIS"
            go run cmd/analyze/analyze.go -logfile "${test_dir}/ims_item.${test_id}.jsonl" -target debezium | tee "${test_dir}/debezium_report.out"
            go run cmd/analyze/analyze.go -logfile "${test_dir}/shipped_item_volume.${test_id}.jsonl" -target streams | tee "${test_dir}/streams_report.out"
            echo "ANALYSIS DONE"
            # Cleanup for the run
            cleanup
            # Compress data
            xz -v -T 0 -9 "${test_dir}"/*.jsonl "${test_dir}/sar.bin"
            # Stop monitoring
            kill "${pid_sar}" || true
        done
    done
done
