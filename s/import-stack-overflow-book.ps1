#!/usr/bin/env pwsh
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
function exitIfFailed { if ($LASTEXITCODE -ne 0) { exit } }

Remove-Item -Force -ErrorAction SilentlyContinue ./importso.exe

go build -o importso.exe ./cmd/import-stack-overflow
exitIfFailed

./importso.exe $args
Remove-Item -Force -ErrorAction SilentlyContinue ./importso.exe

