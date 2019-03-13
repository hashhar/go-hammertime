#!/bin/bash

basedir="$(pwd)"
report="${basedir}/report.md"
true>"${report}"
for poll_interval in ./p[0-9]*; do
    pushd "${poll_interval}"
    for concurrency in ./*; do
        pushd "${concurrency}"
        for batchsize in ./*; do
            pushd "${batchsize}"
            test_id="${poll_interval//.\//}-${concurrency//.\//}-${batchsize//.\//}"
            {
                printf "## PollInterval=%s, Concurrency=%s, BatchSize=%s\n" "${poll_interval//.\/p//}" "${concurrency//.\/c//}" "${batchsize//.\/b//}"
                printf "\n"
                printf "### Debezium\n"
                printf "\n"
                printf "%s %s\n" "- **Avg. Latency:**" "$(jq .AvgLatency debezium_report.out)"
                printf "%s %s\n" "- **90 Percentile:**" "$(jq .NinetyPercentileLatency debezium_report.out)"
                printf "%s %s\n" "- **99 Percentile:**" "$(jq .NinetyNinePercentileLatency debezium_report.out)"
                printf "%s %s\n" "- **Max Latency:**" "$(jq .MaxLatency debezium_report.out)"
                printf "%s %s\n" "- **Min Latency:**" "$(jq .MinLatency debezium_report.out)"
                printf "\n"
                printf "### Kafka Streams\n"
                printf "\n"
                printf "%s %s\n" "- **Avg. Latency:**" "$(jq .AvgLatency streams_report.out)"
                printf "%s %s\n" "- **90 Percentile:**" "$(jq .NinetyPercentileLatency streams_report.out)"
                printf "%s %s\n" "- **99 Percentile:**" "$(jq .NinetyNinePercentileLatency streams_report.out)"
                printf "%s %s\n" "- **Max Latency:**" "$(jq .MaxLatency streams_report.out)"
                printf "%s %s\n" "- **Min Latency:**" "$(jq .MinLatency streams_report.out)"
                printf "\n"
                printf "### Approximate Messages per Topic\n"
                printf "%s %s\n" "- **Items:**" "$(xzcat "ims_item.${test_id}.jsonl.xz" | wc -l)"
                printf "%s %s\n" "- **Shipped Volume:**" "$(xzcat "shipped_item_volume.${test_id}.jsonl.xz" | wc -l)"
                printf "\n"
                printf "### Load Generator Logs\n"
                printf "%s\n" "\`\`\`"
                cat "load.${test_id}.log"
                printf "%s\n" "\`\`\`"
                printf "\n"
            } >> "${report}"
            popd
        done
        popd
    done
    popd
done
