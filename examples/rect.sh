#!/bin/bash
curl -X POST http://localhost:17000 -d "bgrect 50 70 50 70"
curl -X POST http://localhost:17000 -d "bgrect 10 50 30 80"
curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "update"