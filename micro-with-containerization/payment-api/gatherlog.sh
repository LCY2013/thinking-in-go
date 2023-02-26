#!/bin/bash
docker run --rm -v $(pwd):/datalog -v $(pwd)/filebeat.yml:/usr/share/filebeat/filebeat.yml elastic/filebeat:8.6.2
