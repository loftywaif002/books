#!/usr/bin/env pwsh

# you can pass additional args like:
# -update-go-deps
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
function exitIfFailed { if ($LASTEXITCODE -ne 0) { exit } }

go build -o ./gen.exe ./cmd/gen-books
exitIfFailed

Remove-Item -Force -ErrorAction SilentlyContinue ./gen.exe
