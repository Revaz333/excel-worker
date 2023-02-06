#!/bin/bash
while true
do 
    /usr/bin/docker exec -i golang go run ./cmd/build_queue
    /usr/bin/docker exec -i golang go run ./cmd/import
    /usr/bin/docker exec -i golang go run ./cmd/export_table    
    sleep 10
done
