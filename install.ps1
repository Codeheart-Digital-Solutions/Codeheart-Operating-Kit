# Last updated: 2026-06-14T00:13:08Z (UTC)
[CmdletBinding()]
param(
    [string]$Version = "0.1.0",
    [string]$InstallDir = "$env:LOCALAPPDATA\Codeheart\OperatingKit",
    [string]$AssetUrl = "",
    [string]$AssetFile = "",
    [string]$Checksum = "",
    [string]$ChecksumFile = "",
    [string]$Python = "python",
    [switch]$Help
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Show-Usage {
    @"
Install or repair codeheart-operating-kit for the current Windows user.

Options:
  -Version VERSION       Release version to install. Default: 0.1.0
  -InstallDir PATH       User-level install root. Default: %LOCALAPPDATA%\Codeheart\OperatingKit
  -AssetUrl URL          Release asset URL. Defaults to the GitHub release asset.
  -AssetFile PATH        Local release asset path for validation or offline repair.
  -Checksum SHA256       Expected asset SHA-256.
  -ChecksumFile PATH     File containing the expected SHA-256.
  -Python PATH           Python executable to use. Default: python
  -Help                  Show this help.
"@
}

if ($Help) {
    Show-Usage
    exit 0
}

$AssetName = "codeheart-operating-kit-$Version-windows.zip"
if ([string]::IsNullOrWhiteSpace($AssetUrl)) {
    $AssetUrl = "https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v$Version/$AssetName"
}

$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
New-Item -ItemType Directory -Path $TempDir | Out-Null
try {
    $AssetPath = Join-Path $TempDir $AssetName
    if (-not [string]::IsNullOrWhiteSpace($AssetFile)) {
        Copy-Item -LiteralPath $AssetFile -Destination $AssetPath
    } elseif ($AssetUrl.StartsWith("file://")) {
        $LocalAssetPath = ([System.Uri]$AssetUrl).LocalPath
        Copy-Item -LiteralPath $LocalAssetPath -Destination $AssetPath
    } else {
        Invoke-WebRequest -Uri $AssetUrl -OutFile $AssetPath
    }

    if ([string]::IsNullOrWhiteSpace($Checksum)) {
        if (-not [string]::IsNullOrWhiteSpace($ChecksumFile)) {
            $Checksum = ((Get-Content -LiteralPath $ChecksumFile -First 1) -split "\s+")[0]
        } elseif ($AssetUrl.StartsWith("file://")) {
            $ChecksumPath = ([System.Uri]("$AssetUrl.sha256")).LocalPath
            $Checksum = ((Get-Content -LiteralPath $ChecksumPath -First 1) -split "\s+")[0]
        } else {
            $ChecksumPath = Join-Path $TempDir "$AssetName.sha256"
            Invoke-WebRequest -Uri "$AssetUrl.sha256" -OutFile $ChecksumPath
            $Checksum = ((Get-Content -LiteralPath $ChecksumPath -First 1) -split "\s+")[0]
        }
    }

    if ([string]::IsNullOrWhiteSpace($Checksum)) {
        throw "Expected checksum is required; installation stopped."
    }

    $ActualChecksum = (Get-FileHash -Algorithm SHA256 -LiteralPath $AssetPath).Hash.ToLowerInvariant()
    if ($ActualChecksum -ne $Checksum.ToLowerInvariant()) {
        throw "Checksum mismatch for $AssetName; installation stopped."
    }

    $ExtractDir = Join-Path $TempDir "extract"
    Expand-Archive -LiteralPath $AssetPath -DestinationPath $ExtractDir
    $Wheel = Get-ChildItem -Path $ExtractDir -Recurse -Filter "codeheart_operating_kit-*.whl" | Select-Object -First 1
    if ($null -eq $Wheel) {
        throw "Release asset did not contain a codeheart-operating-kit wheel."
    }

    $BinDir = Join-Path $InstallDir "bin"
    $LibDir = Join-Path $InstallDir "lib"
    New-Item -ItemType Directory -Force -Path $BinDir, $LibDir | Out-Null
    & $Python -m pip install --upgrade --target $LibDir $Wheel.FullName
    if ($LASTEXITCODE -ne 0) {
        throw "pip install failed."
    }

    $Wrapper = Join-Path $BinDir "codeheart-operating-kit.cmd"
    $WrapperContent = @"
@echo off
set "PYTHONPATH=$LibDir;%PYTHONPATH%"
"$Python" -m codeheart_operating_kit.cli %*
"@
    Set-Content -LiteralPath $Wrapper -Value $WrapperContent -Encoding ASCII

    Write-Host "codeheart-operating-kit installed at $Wrapper"
    $PathEntries = ($env:PATH -split ";") | Where-Object { $_ -ne "" }
    if ($PathEntries -notcontains $BinDir) {
        Write-Host "Add this folder to PATH to run it by name: $BinDir"
    }
    Write-Host "Next: codeheart-operating-kit onboard"
} finally {
    Remove-Item -LiteralPath $TempDir -Recurse -Force -ErrorAction SilentlyContinue
}
