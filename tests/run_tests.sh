#!/bin/bash

# Selenium Test Runner Script for Library Management System
# This script starts the application, ChromeDriver, runs tests, and cleans up

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_PORT=8080
CHROMEDRIVER_PORT=4444
TEST_TIMEOUT=10m

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Library Management System - Test Runner${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Function to check if a port is in use
check_port() {
    lsof -i:$1 > /dev/null 2>&1
}

# Function to wait for service
wait_for_service() {
    local port=$1
    local service=$2
    local max_attempts=30
    local attempt=1
    
    echo -e "${YELLOW}Waiting for $service on port $port...${NC}"
    while ! check_port $port; do
        if [ $attempt -eq $max_attempts ]; then
            echo -e "${RED}Failed to start $service${NC}"
            return 1
        fi
        sleep 1
        attempt=$((attempt + 1))
    done
    echo -e "${GREEN}$service is ready${NC}"
}

# Function to cleanup
cleanup() {
    echo ""
    echo -e "${YELLOW}Cleaning up...${NC}"
    
    # Stop ChromeDriver
    pkill -f chromedriver 2>/dev/null || true
    
    # Stop application
    if [ ! -z "$APP_PID" ]; then
        kill $APP_PID 2>/dev/null || true
    fi
    
    echo -e "${GREEN}Cleanup completed${NC}"
}

# Set trap to cleanup on exit
trap cleanup EXIT INT TERM

# Step 1: Check prerequisites
echo -e "${YELLOW}Step 1: Checking prerequisites...${NC}"

if ! command -v chromedriver &> /dev/null; then
    echo -e "${RED}ChromeDriver not found. Please install it:${NC}"
    echo "  Ubuntu/Debian: sudo apt-get install chromium-chromedriver"
    echo "  macOS: brew install chromedriver"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo -e "${RED}Go not found. Please install Go${NC}"
    exit 1
fi

echo -e "${GREEN}Prerequisites OK${NC}"
echo ""

# Step 2: Build the application
echo -e "${YELLOW}Step 2: Building application...${NC}"
go build -o server main.go
echo -e "${GREEN}Build completed${NC}"
echo ""

# Step 3: Start the application
echo -e "${YELLOW}Step 3: Starting application...${NC}"
./server > /dev/null 2>&1 &
APP_PID=$!

if ! wait_for_service $APP_PORT "Application"; then
    echo -e "${RED}Failed to start application${NC}"
    exit 1
fi
echo ""

# Step 4: Start ChromeDriver
echo -e "${YELLOW}Step 4: Starting ChromeDriver...${NC}"
chromedriver --port=$CHROMEDRIVER_PORT > /dev/null 2>&1 &
CHROMEDRIVER_PID=$!

if ! wait_for_service $CHROMEDRIVER_PORT "ChromeDriver"; then
    echo -e "${RED}Failed to start ChromeDriver${NC}"
    exit 1
fi
echo ""

# Step 5: Run tests
echo -e "${YELLOW}Step 5: Running Selenium tests...${NC}"
echo ""

TEST_RESULT=0
go test -v -timeout $TEST_TIMEOUT ./tests/ || TEST_RESULT=$?

echo ""

# Step 6: Report results
if [ $TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}All tests passed successfully!${NC}"
    echo -e "${GREEN}========================================${NC}"
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}Some tests failed${NC}"
    echo -e "${RED}========================================${NC}"
fi

exit $TEST_RESULT
