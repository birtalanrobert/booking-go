#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings -dbname=bookings -dbuser=robertbirtalan -cache=false -production=false
