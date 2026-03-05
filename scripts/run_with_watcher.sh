#!/bin/bash
# scripts/run_with_watcher.sh
# Runs both the PDF watcher and the Go server.
# Usage: ./scripts/run_with_watcher.sh [--server-only] [--watcher-only]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Parse arguments
WATCHER_ONLY=false
SERVER_ONLY=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --watcher-only)
            WATCHER_ONLY=true
            shift
            ;;
        --server-only)
            SERVER_ONLY=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--watcher-only] [--server-only]"
            exit 1
            ;;
    esac
done

cd "$PROJECT_ROOT"

echo -e "${BLUE}🎯 Garfex Calculator - Dev Runner${NC}"
echo ""

if [ "$SERVER_ONLY" = true ]; then
    echo -e "${GREEN}▶ Starting server only...${NC}"
    echo ""
    go run cmd/server/main.go
elif [ "$WATCHER_ONLY" = true ]; then
    echo -e "${GREEN}▶ Starting PDF watcher only...${NC}"
    echo ""
    go run cmd/pdf_watcher/main.go
else
    # Run both - use subshells to run in parallel
    echo -e "${GREEN}▶ Starting server...${NC} (Ctrl+C to stop both)"
    echo -e "${GREEN}▶ Starting PDF watcher...${NC} (Ctrl+C to stop both)"
    echo ""
    echo -e "${YELLOW}Server: http://localhost:8080${NC}"
    echo -e "${YELLOW}Watcher: watching internal/pdf/templates/${NC}"
    echo ""

    # Trap to kill both processes on Ctrl+C
    cleanup() {
        echo ""
        echo -e "${RED}🛑 Stopping...${NC}"
        kill %1 2>/dev/null || true
        kill %2 2>/dev/null || true
        exit 0
    }
    trap cleanup SIGINT

    # Run both in background
    go run cmd/server/main.go &
    SERVER_PID=$!
    
    go run cmd/pdf_watcher/main.go &
    WATCHER_PID=$!

    # Wait for either to finish
    wait -n
fi
