$env:GOARCH = "amd64"
$env:GOOS = "windows"
$env:CGO_ENABLED = "0"

$ldflags = "-s -w -extldflags=-static"

Write-Host "Building for Windows amd64..." -ForegroundColor Green

go build -trimpath -ldflags="$ldflags" -o bin/sauron.exe ./cmd/sauron/

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful! Output: bin/sauron.exe" -ForegroundColor Green
    Get-Item bin/sauron.exe | Format-Table Name, Length, LastWriteTime
} else {
    Write-Host "Build failed!" -ForegroundColor Red
}