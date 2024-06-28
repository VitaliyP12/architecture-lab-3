#!/bin/bash

#curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "figure 300 200"
curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "figure 500 500"
curl -X POST http://localhost:17000 -d "update"