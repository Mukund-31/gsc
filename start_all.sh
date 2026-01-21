#!/bin/bash

# Ensure we use the correct go binary
export PATH=$PWD/tools/go/bin:$PATH

# Kill any existing processes on our ports
lsof -t -i:8080 -i:8081 -i:8082 -i:8083 -i:8084 -i:8085 | xargs kill -9 2>/dev/null

echo "🚀 Starting Responsive AI Clusters (Official A2A SDK)..."

# Build all binaries
echo "📦 Building services..."
cd backend
go mod tidy
go build -o ../bin/warehouse ./cmd/warehouse
go build -o ../bin/outlet ./cmd/outlet
go build -o ../bin/coordinator ./cmd/coordinator
cd ..

# Check if builds succeeded
if [ ! -f bin/warehouse ] || [ ! -f bin/outlet ] || [ ! -f bin/coordinator ]; then
    echo "❌ Build failed. Please check errors above."
    exit 1
fi

# Start Warehouse
echo "🏭 Starting Warehouse Server on :8081"
./bin/warehouse &
WAREHOUSE_PID=$!

sleep 2

# Start Outlets
echo "🏪 Starting Outlet Servers..."
./bin/outlet -id="Outlet-1" -name="Outlet North" -port="8082" &
./bin/outlet -id="Outlet-2" -name="Outlet South" -port="8083" &
./bin/outlet -id="Outlet-3" -name="Outlet East" -port="8084" &
./bin/outlet -id="Outlet-4" -name="Outlet West" -port="8085" &

sleep 2

# Start Coordinator
echo "🧠 Starting Coordinator on :8080"
./bin/coordinator &
COORD_PID=$!

echo "✅ All systems go! Open http://localhost:8080"

# Wait
wait
