$ErrorActionPreference = "Stop"

$moduleRoot = Split-Path $PSScriptRoot -Parent
$source = Join-Path (Split-Path $moduleRoot -Parent | Split-Path -Parent) "agentics"
$versionSource = Join-Path (Split-Path $moduleRoot -Parent | Split-Path -Parent) "VERSION"
$destination = Join-Path $moduleRoot "internal\source\embedded\agentics"
$versionDestination = Join-Path $moduleRoot "internal\source\embedded\VERSION"

if (-not (Test-Path $source -PathType Container)) {
    throw "Methodology source not found: $source"
}
if (-not (Test-Path $versionSource -PathType Leaf)) {
    throw "Version source not found: $versionSource"
}

Remove-Item $destination -Recurse -Force -ErrorAction SilentlyContinue
New-Item (Split-Path $destination -Parent) -ItemType Directory -Force | Out-Null
Copy-Item $source $destination -Recurse
Copy-Item $versionSource $versionDestination -Force