# Last updated: 2026-07-04T21:42:59Z (UTC)
[CmdletBinding()]
param(
    [string]$Version = "0.1.19",
    [string]$InstallDir = "$env:LOCALAPPDATA\Codeheart\OperatingKit",
    [string]$AssetUrl = "",
    [string]$AssetFile = "",
    [string]$Checksum = "",
    [string]$ChecksumFile = "",
    [string]$Python = "",
    [switch]$Help
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Show-Usage {
    @"
Install or repair codeheart-operating-kit for the current Windows user.

Options:
  -Version VERSION       Release version to install. Default: 0.1.19
  -InstallDir PATH       User-level install root. Default: %LOCALAPPDATA%\Codeheart\OperatingKit
  -AssetUrl URL          Release asset URL. Defaults to the GitHub release asset.
  -AssetFile PATH        Local release asset path for validation or offline repair.
  -Checksum SHA256       Expected asset SHA-256.
  -ChecksumFile PATH     File containing the expected SHA-256.
  -Help                  Show this help.
"@
}

if ($Help) {
    Show-Usage
    exit 0
}

if ($PSBoundParameters.ContainsKey("Python") -and -not [string]::IsNullOrWhiteSpace($Python)) {
    Write-Warning "-Python is deprecated and ignored; the Operating Kit installer uses a self-contained binary."
}

$AssetName = "codeheart-operating-kit-$Version-windows-x64.zip"
if ([string]::IsNullOrWhiteSpace($AssetUrl)) {
    $AssetUrl = "https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v$Version/$AssetName"
}

$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
$StagingDir = ""
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
    try {
        Expand-Archive -LiteralPath $AssetPath -DestinationPath $ExtractDir
    } catch {
        throw "Release asset could not be extracted; installation stopped."
    }

    $Binary = Get-ChildItem -Path $ExtractDir -Recurse -Filter "codeheart-operating-kit.exe" |
        Where-Object { $_.FullName -match "[\\/]+bin[\\/]+codeheart-operating-kit\.exe$" } |
        Select-Object -First 1
    if ($null -eq $Binary) {
        throw "Release asset did not contain bin/codeheart-operating-kit.exe."
    }

    $BinDir = Join-Path $InstallDir "bin"
    $LibDir = Join-Path $InstallDir "lib"
    $TargetBinary = Join-Path $BinDir "codeheart-operating-kit.exe"
    $Shim = Join-Path $BinDir "codeheart-operating-kit.cmd"
    New-Item -ItemType Directory -Force -Path $BinDir | Out-Null

    $LegacyFound = $false
    if (Test-Path -LiteralPath $Shim) {
        $ShimText = Get-Content -LiteralPath $Shim -Raw
        if ($ShimText -match "python|codeheart_operating_kit") {
            $LegacyFound = $true
        }
    }
    if (Test-Path -LiteralPath $LibDir) {
        $LegacyPayload = Get-ChildItem -Path $LibDir -Filter "codeheart_operating_kit*" -ErrorAction SilentlyContinue | Select-Object -First 1
        if ($null -ne $LegacyPayload) {
            $LegacyFound = $true
        }
    }
    if ($LegacyFound) {
        Write-Host "Legacy Python install detected; installing the self-contained binary and preserving legacy files."
    }

    $StagingDir = Join-Path $InstallDir (".staging." + [System.Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Force -Path $StagingDir | Out-Null
    $StagedBinary = Join-Path $StagingDir "codeheart-operating-kit.exe"
    Copy-Item -LiteralPath $Binary.FullName -Destination $StagedBinary

    $SmokeOutput = & $StagedBinary --version 2>&1
    if ($LASTEXITCODE -ne 0) {
        $SmokeOutput | ForEach-Object { Write-Error $_ }
        throw "Staged binary validation failed; previous runnable command preserved."
    }

    Move-Item -LiteralPath $StagedBinary -Destination $TargetBinary -Force
$ShimContent = @"
@echo off
set "CODEHEART_OPERATING_KIT_CLI=1"
"%~dp0codeheart-operating-kit.exe" %*
"@
    Set-Content -LiteralPath $Shim -Value $ShimContent -Encoding ASCII

    Write-Host "codeheart-operating-kit installed at $TargetBinary"
    if ($LegacyFound) {
        Write-Host "Legacy Python files were preserved under $InstallDir for later cleanup."
    }
    $PathEntries = ($env:PATH -split ";") | Where-Object { $_ -ne "" }
    if ($PathEntries -notcontains $BinDir) {
        Write-Host "Add this folder to PATH to run it by name: $BinDir"
    }
    Write-Host "Next: codeheart-operating-kit onboard"
} finally {
    if (-not [string]::IsNullOrWhiteSpace($StagingDir)) {
        Remove-Item -LiteralPath $StagingDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    Remove-Item -LiteralPath $TempDir -Recurse -Force -ErrorAction SilentlyContinue
}
