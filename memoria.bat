@echo off
REM Test script for PDF generation endpoint
REM Usage: memoria.bat [server_url]
REM Default server: http://localhost:8080

set SERVER_URL=%1
if "%SERVER_URL%"=="" set SERVER_URL=http://localhost:8080

echo.
echo ========================================
echo Testing PDF Generation Endpoint
echo ========================================
echo Server: %SERVER_URL%
echo Endpoint: %SERVER_URL%/api/v1/pdf/memoria
echo Input: test_memoria.json
echo Output: memoria.pdf
echo.

REM Check if test_memoria.json exists
if not exist "test_memoria.json" (
    echo ERROR: test_memoria.json not found in current directory
    exit /b 1
)

echo Making POST request to generate PDF...
echo.

REM Use PowerShell to make the request (available on Windows 10+)
powershell -NoProfile -Command "
$headers = @{
    'Content-Type' = 'application/json'
}

try {
    $response = Invoke-WebRequest -Uri '%SERVER_URL%/api/v1/pdf/memoria' -Method POST -Headers $headers -InFile 'test_memoria.json' -OutFile 'memoria.pdf' -ErrorAction Stop
    Write-Host 'SUCCESS: PDF generated successfully'
    Write-Host 'Status Code:' $response.StatusCode
    Write-Host 'File size:' (Get-Item 'memoria.pdf').Length 'bytes'
} catch {
    Write-Host 'ERROR:' $_.Exception.Message
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode
        Write-Host 'Status Code:' $statusCode
    }
    exit /b 1
}
"

if %ERRORLEVEL% neq 0 (
    echo.
    echo FAILED: Could not generate PDF
    exit /b 1
)

echo.
echo Opening PDF...
start "" "memoria.pdf"

echo.
echo DONE! Check memoria.pdf
