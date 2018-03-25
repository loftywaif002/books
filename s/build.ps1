#!/usr/bin/env pwsh

# you can pass additional args like:
# -update-go-deps
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
function exitIfFailed { if ($LASTEXITCODE -ne 0) { exit } }

Remove-Item -Force -ErrorAction SilentlyContinue ./cmd/gen-books/preview.exe

Set-Location -Path cmd/gen-books
go build -o preview.exe
Set-Location -Path ../..
exitIfFailed

Remove-Item -Force -ErrorAction SilentlyContinue ./cmd/gen-books/preview.exe