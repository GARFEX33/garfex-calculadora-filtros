#!/bin/bash
# Test script for PDF generation endpoint
# Usage: ./memoria.sh [server_url]
# Default server: http://localhost:8080

SERVER_URL="${1:-http://localhost:8080}"

echo ""
echo "========================================"
echo "Testing PDF Generation Endpoint"
echo "========================================"
echo "Server: $SERVER_URL"
echo "Endpoint: $SERVER_URL/api/v1/pdf/memoria"
echo "Input: test_memoria.json"
echo "Output: memoria.pdf"
echo ""

# Check if test_memoria.json exists
if [ ! -f "test_memoria.json" ]; then
    echo "ERROR: test_memoria.json not found in current directory"
    exit 1
fi

echo "Making POST request to generate PDF..."
echo ""

# Make the request using curl
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -d @test_memoria.json \
    -o memoria.pdf \
    "$SERVER_URL/api/v1/pdf/memoria" 2>&1)

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "200" ]; then
    echo "SUCCESS: PDF generated successfully"
    echo "Status Code: $HTTP_CODE"
    FILE_SIZE=$(stat -f%z "memoria.pdf" 2>/dev/null || stat -c%s "memoria.pdf" 2>/dev/null)
    echo "File size: $FILE_SIZE bytes"
else
    echo "ERROR: HTTP Status Code: $HTTP_CODE"
    echo "Response:"
    echo "$RESPONSE" | head -n -1
    rm -f memoria.pdf
    exit 1
fi

echo ""
echo "Opening PDF..."

# Detect OS and open PDF
if [[ "$OSTYPE" == "darwin"* ]]; then
    open memoria.pdf
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    xdg-open memoria.pdf
else
    echo "Unknown OS. PDF saved at: memoria.pdf"
fi

echo ""
echo "DONE! Check memoria.pdf"
