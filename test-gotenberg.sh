<!bin/bash

# Test script para Gotenberg
# Ejecutar desde el directorio del proyecto

GOTENBERG_URL="http://10.255.254.254:3000"
TEST_DIR="/tmp/gotenberg-test"
mkdir -p "$TEST_DIR"

echo "=== Gotenberg Test Script ==="
echo ""

# Crear HTML de test
cat > "$TEST_DIR/test.html" << 'EOF'
<!DOCTYPE html>
<html>
<head><title>Test PDF</title></head>
<body><h1>Hello World - Test PDF</h1><p>Generated at: $(date)</p></body>
</html>
EOF

# Crear Markdown de test
cat > "$TEST_DIR/test.md" << 'EOF'
# Hello World

This is a **markdown** test document.

Generated at: $(date)
EOF

echo "1. Verificando versión de Gotenberg..."
curl -s "$GOTENBERG_URL/version"
echo ""
echo ""

echo "2. Test Approach 1: files=@test.html (current approach)"
curl -s -X POST "$GOTENBERG_URL/forms/chromium/convert/html" \
  -F "files=@$TEST_DIR/test.html" \
  -o "$TEST_DIR/output1.pdf"
echo "Output size: $(wc -c < "$TEST_DIR/output1.pdf") bytes"
file "$TEST_DIR/output1.pdf"
echo ""

echo "3. Test Approach 2: index.html=@test.html (old approach)"
curl -s -X POST "$GOTENBERG_URL/forms/chromium/convert/html" \
  -F "index.html=@$TEST_DIR/test.html" \
  -o "$TEST_DIR/output2.pdf"
echo "Output size: $(wc -c < "$TEST_DIR/output2.pdf") bytes"
file "$TEST_DIR/output2.pdf"
echo ""

echo "4. Test Approach 3: Markdown conversion"
curl -s -X POST "$GOTENBERG_URL/forms/chromium/convert/markdown" \
  -F "files=@$TEST_DIR/test.md" \
  -o "$TEST_DIR/output3.pdf"
echo "Output size: $(wc -c < "$TEST_DIR/output3.pdf") bytes"
file "$TEST_DIR/output3.pdf"
echo ""

echo "5. Verbose test para ver request/response"
echo "   (mirá el Content-Type y qué devuelve)"
curl -v -X POST "$GOTENBERG_URL/forms/chromium/convert/html" \
  -F "files=@$TEST_DIR/test.html" 2>&1 | head -50
echo ""

echo "=== Tests completados ==="
echo "Archivos generados en: $TEST_DIR"
ls -la "$TEST_DIR"/*.pdf 2>/dev/null || echo "No se generó ningún PDF"
