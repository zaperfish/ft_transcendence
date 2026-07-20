#!/usr/bin/env sh
for i in $(seq 1 100); do
    curl -sk -o /dev/null -w "%{http_code}\n" https://localhost:7443/api/health
done | sort | uniq -c
