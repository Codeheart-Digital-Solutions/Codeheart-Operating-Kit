# Last updated: 2026-07-10T11:35:02Z (UTC)
[CmdletBinding()]
param(
    [string]$Version = "0.1.23",
    [string]$InstallDir = "$env:LOCALAPPDATA\Codeheart\OperatingKit",
    [string]$CatalogUrl = "",
    [string]$CatalogFile = "",
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
  -Version VERSION       Release version to install. Default: 0.1.23
  -InstallDir PATH       User-level install root. Default: %LOCALAPPDATA%\Codeheart\OperatingKit
  -CatalogUrl URL        External release catalog URL.
  -CatalogFile PATH      Local external release catalog.
  -AssetUrl URL          Optional release asset URL override.
  -AssetFile PATH        Local release asset for offline validation.
  -Checksum SHA256       Optional checksum that must agree with the catalog.
  -ChecksumFile PATH     Optional checksum file that must agree with the catalog.
  -Help                  Show this help.
"@
}

function Get-Sha256([string]$Path) {
    return (Get-FileHash -Algorithm SHA256 -LiteralPath $Path).Hash.ToLowerInvariant()
}

function Assert-ExactProperties($Value, [string[]]$Expected, [string]$Label) {
    $Actual = @($Value.PSObject.Properties.Name | Sort-Object)
    $Wanted = @($Expected | Sort-Object)
    if (($Actual -join "|") -ne ($Wanted -join "|")) {
        throw "$Label has unexpected or missing fields."
    }
}

if ($Help) { Show-Usage; exit 0 }
if ($PSBoundParameters.ContainsKey("Python") -and -not [string]::IsNullOrWhiteSpace($Python)) {
    Write-Warning "-Python is deprecated and ignored; the Operating Kit installer uses a self-contained binary."
}

$ReleaseBaseUrl = "https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v$Version"
$CatalogName = "release-catalog-$Version.json"
if ([string]::IsNullOrWhiteSpace($CatalogUrl) -and [string]::IsNullOrWhiteSpace($CatalogFile)) {
    if (-not [string]::IsNullOrWhiteSpace($AssetFile)) {
        $AssetParent = Split-Path -Parent $AssetFile
        if ([string]::IsNullOrWhiteSpace($AssetParent)) { $AssetParent = (Get-Location).Path }
        $CandidateCatalog = Join-Path $AssetParent $CatalogName
        if (Test-Path -LiteralPath $CandidateCatalog) { $CatalogFile = $CandidateCatalog }
    }
    if ([string]::IsNullOrWhiteSpace($CatalogFile)) { $CatalogUrl = "$ReleaseBaseUrl/$CatalogName" }
}
if (-not [string]::IsNullOrWhiteSpace($CatalogUrl) -and -not $CatalogUrl.StartsWith("https://") -and -not $CatalogUrl.StartsWith("file://")) {
    throw "Release catalogs must use HTTPS or a local file URL."
}

$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString("N"))
$StagingDir = ""
New-Item -ItemType Directory -Path $TempDir | Out-Null
try {
    $CatalogPath = Join-Path $TempDir $CatalogName
    if (-not [string]::IsNullOrWhiteSpace($CatalogFile)) {
        Copy-Item -LiteralPath $CatalogFile -Destination $CatalogPath
    } elseif ($CatalogUrl.StartsWith("file://")) {
        Copy-Item -LiteralPath ([System.Uri]$CatalogUrl).LocalPath -Destination $CatalogPath
    } else {
        Invoke-WebRequest -Uri $CatalogUrl -OutFile $CatalogPath
    }
    $Catalog = Get-Content -LiteralPath $CatalogPath -Raw | ConvertFrom-Json
    Assert-ExactProperties $Catalog @("schema_version", "version", "assets") "Release catalog"
    if ($Catalog.schema_version -ne 1 -or $Catalog.version -ne $Version) {
        throw "Release catalog identity does not match the requested version."
    }
    $Matches = @($Catalog.assets | Where-Object { $_.version -eq $Version -and $_.platform -eq "windows-x64" })
    if ($Matches.Count -ne 1) { throw "Release catalog must contain exactly one windows-x64 asset for $Version." }
    $Asset = $Matches[0]
    Assert-ExactProperties $Asset @("name", "version", "platform", "url", "archive_sha256", "pack_manifest_sha256") "Release catalog asset"
    if ($Asset.name -ne "codeheart-operating-kit-$Version-windows-x64.zip") {
        throw "Release catalog asset name does not match the requested version and platform."
    }
    if ($Asset.archive_sha256 -notmatch "^[a-fA-F0-9]{64}$" -or $Asset.pack_manifest_sha256 -notmatch "^[a-fA-F0-9]{64}$") {
        throw "Release catalog contains an invalid digest."
    }

    if (-not [string]::IsNullOrWhiteSpace($ChecksumFile)) {
        $Checksum = ((Get-Content -LiteralPath $ChecksumFile -First 1) -split "\s+")[0]
    }
    if (-not [string]::IsNullOrWhiteSpace($Checksum) -and $Checksum.ToLowerInvariant() -ne $Asset.archive_sha256.ToLowerInvariant()) {
        throw "Checksum mismatch: explicit checksum disagrees with the external release catalog."
    }

    if ([string]::IsNullOrWhiteSpace($AssetUrl)) {
        $AssetUrl = $Asset.url
        if (-not [System.Uri]::IsWellFormedUriString($AssetUrl, [System.UriKind]::Absolute)) {
            if (-not [string]::IsNullOrWhiteSpace($CatalogFile)) {
                $AssetUrl = Join-Path (Split-Path -Parent $CatalogFile) $AssetUrl
            } else {
                $AssetUrl = ([System.Uri]::new([System.Uri]$CatalogUrl, $AssetUrl)).AbsoluteUri
            }
        }
    }
    $AssetPath = Join-Path $TempDir $Asset.name
    if (-not [string]::IsNullOrWhiteSpace($AssetFile)) {
        Copy-Item -LiteralPath $AssetFile -Destination $AssetPath
    } elseif ($AssetUrl.StartsWith("file://")) {
        Copy-Item -LiteralPath ([System.Uri]$AssetUrl).LocalPath -Destination $AssetPath
    } elseif ($AssetUrl.StartsWith("https://")) {
        Invoke-WebRequest -Uri $AssetUrl -OutFile $AssetPath
    } elseif ([System.Uri]::IsWellFormedUriString($AssetUrl, [System.UriKind]::Absolute)) {
        throw "Release assets must use HTTPS or a local file URL."
    } else {
        Copy-Item -LiteralPath $AssetUrl -Destination $AssetPath
    }
    if ((Get-Sha256 $AssetPath) -ne $Asset.archive_sha256.ToLowerInvariant()) {
        throw "Checksum mismatch for $($Asset.name); installation stopped."
    }

    Add-Type -AssemblyName System.IO.Compression.FileSystem
    $Zip = [System.IO.Compression.ZipFile]::OpenRead($AssetPath)
    try {
        foreach ($Entry in $Zip.Entries) {
            $Name = $Entry.FullName
            if ([string]::IsNullOrWhiteSpace($Name) -or $Name.StartsWith("/") -or $Name.Contains("\") -or $Name -match "(^|/)\.\.(/|$)") {
                throw "Release asset contains an unsafe archive path."
            }
            $UnixType = (($Entry.ExternalAttributes -shr 16) -band 0xF000)
            if ($UnixType -notin @(0, 0x4000, 0x8000)) {
                throw "Release asset contains a symbolic link or unsupported filesystem entry."
            }
        }
    } finally {
        $Zip.Dispose()
    }

    $ExtractDir = Join-Path $TempDir "extract"
    try { Expand-Archive -LiteralPath $AssetPath -DestinationPath $ExtractDir } catch {
        throw "Release asset could not be extracted; installation stopped."
    }
    $PackManifests = @(Get-ChildItem -Path $ExtractDir -Recurse -Filter "pack-manifest.json" -File)
    if ($PackManifests.Count -ne 1) { throw "Release asset must contain exactly one pack-manifest.json." }
    $PackManifestPath = $PackManifests[0].FullName
    if ((Get-Sha256 $PackManifestPath) -ne $Asset.pack_manifest_sha256.ToLowerInvariant()) {
        throw "Pack manifest checksum mismatch."
    }
    $Pack = Get-Content -LiteralPath $PackManifestPath -Raw | ConvertFrom-Json
    Assert-ExactProperties $Pack @(
        "schema_version", "version", "platform", "command", "binary_path", "binary_sha256",
        "content_manifest_path", "content_manifest_sha256", "payload_checksums_path", "payload_checksums_sha256"
    ) "Pack manifest"
    if ($Pack.schema_version -ne 1 -or $Pack.version -ne $Version -or $Pack.platform -ne "windows-x64" -or $Pack.command -ne "codeheart-operating-kit") {
        throw "Pack identity does not match the requested command, version, and platform."
    }
    if ($Pack.binary_path -ne "bin/codeheart-operating-kit.exe" -or $Pack.content_manifest_path -ne "content-manifest.yaml" -or $Pack.payload_checksums_path -ne "checksums.txt") {
        throw "Pack identity paths are invalid."
    }
    $PayloadRoot = Split-Path -Parent $PackManifestPath
    $BinaryPath = Join-Path $PayloadRoot ($Pack.binary_path -replace "/", "\")
    $ContentPath = Join-Path $PayloadRoot $Pack.content_manifest_path
    $ChecksumsPath = Join-Path $PayloadRoot $Pack.payload_checksums_path
    if (-not (Test-Path -LiteralPath $BinaryPath) -or (Get-Sha256 $BinaryPath) -ne $Pack.binary_sha256.ToLowerInvariant()) { throw "Binary checksum mismatch." }
    if (-not (Test-Path -LiteralPath $ContentPath) -or (Get-Sha256 $ContentPath) -ne $Pack.content_manifest_sha256.ToLowerInvariant()) { throw "Content manifest checksum mismatch." }
    if (-not (Test-Path -LiteralPath $ChecksumsPath) -or (Get-Sha256 $ChecksumsPath) -ne $Pack.payload_checksums_sha256.ToLowerInvariant()) { throw "Payload checksum identity mismatch." }

    $BinaryListed = $false
    foreach ($Line in Get-Content -LiteralPath $ChecksumsPath) {
        $Parts = $Line -split "  ", 2
        if ($Parts.Count -ne 2) { throw "Payload checksum line is invalid." }
        $Expected = $Parts[0].ToLowerInvariant()
        $Relative = $Parts[1]
        if ($Relative.StartsWith("/") -or $Relative.Contains("\") -or $Relative -match "(^|/)\.\.(/|$)" -or $Relative -in @("pack-manifest.json", "checksums.txt")) {
            throw "Payload checksum path is invalid."
        }
        $PayloadFile = Join-Path $PayloadRoot ($Relative -replace "/", "\")
        if (-not (Test-Path -LiteralPath $PayloadFile) -or (Get-Sha256 $PayloadFile) -ne $Expected) {
            throw "Payload checksum mismatch for $Relative"
        }
        if ($Relative -eq $Pack.binary_path) { $BinaryListed = $true }
    }
    if (-not $BinaryListed) { throw "Binary is absent from payload checksums." }

    $BinDir = Join-Path $InstallDir "bin"
    $LibDir = Join-Path $InstallDir "lib"
    $TargetBinary = Join-Path $BinDir "codeheart-operating-kit.exe"
    $Shim = Join-Path $BinDir "codeheart-operating-kit.cmd"
    New-Item -ItemType Directory -Force -Path $BinDir | Out-Null

    $LegacyFound = $false
    if (Test-Path -LiteralPath $Shim) {
        if ((Get-Content -LiteralPath $Shim -Raw) -match "python|codeheart_operating_kit") { $LegacyFound = $true }
    }
    if (Test-Path -LiteralPath $LibDir) {
        if (Get-ChildItem -Path $LibDir -Filter "codeheart_operating_kit*" -ErrorAction SilentlyContinue | Select-Object -First 1) { $LegacyFound = $true }
    }
    if ($LegacyFound) { Write-Host "Legacy Python install detected; installing the self-contained binary and preserving legacy files." }

    $StagingDir = Join-Path $InstallDir (".staging." + [System.Guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Force -Path $StagingDir | Out-Null
    $StagedBinary = Join-Path $StagingDir "codeheart-operating-kit.exe"
    Copy-Item -LiteralPath $BinaryPath -Destination $StagedBinary
    $SmokeOutput = (& $StagedBinary --version 2>&1 | Out-String).Trim()
    if ($LASTEXITCODE -ne 0 -or $SmokeOutput -ne "codeheart-operating-kit $Version") {
        throw "Staged binary validation failed; previous runnable command preserved."
    }
    & $StagedBinary __verify-content-identity --path $ContentPath --version $Version
    if ($LASTEXITCODE -ne 0) { throw "Staged content identity validation failed; previous runnable command preserved." }
    & $StagedBinary __verify-release-evidence --catalog $CatalogPath --version $Version
    if ($LASTEXITCODE -ne 0) { throw "Staged release evidence validation failed; previous runnable command preserved." }

    $BackupBinary = Join-Path $StagingDir "previous-binary.exe"
    $BackupShim = Join-Path $StagingDir "previous-shim.cmd"
    $HadBinary = Test-Path -LiteralPath $TargetBinary
    $HadShim = Test-Path -LiteralPath $Shim
    if ($HadBinary) { Move-Item -LiteralPath $TargetBinary -Destination $BackupBinary }
    if ($HadShim) { Move-Item -LiteralPath $Shim -Destination $BackupShim }
    try {
        Move-Item -LiteralPath $StagedBinary -Destination $TargetBinary
        $ShimContent = @"
@echo off
set "CODEHEART_OPERATING_KIT_CLI=1"
"%~dp0codeheart-operating-kit.exe" %*
"@
        Set-Content -LiteralPath $Shim -Value $ShimContent -Encoding ASCII
        $InstalledVersion = (& $TargetBinary --version 2>&1 | Out-String).Trim()
        if ($LASTEXITCODE -ne 0 -or $InstalledVersion -ne "codeheart-operating-kit $Version") { throw "Installed binary validation failed." }
    } catch {
        Remove-Item -LiteralPath $TargetBinary -Force -ErrorAction SilentlyContinue
        Remove-Item -LiteralPath $Shim -Force -ErrorAction SilentlyContinue
        if ($HadBinary) { Move-Item -LiteralPath $BackupBinary -Destination $TargetBinary }
        if ($HadShim) { Move-Item -LiteralPath $BackupShim -Destination $Shim }
        throw "Binary replacement failed; previous runnable command restored. $($_.Exception.Message)"
    }

    Write-Host "codeheart-operating-kit installed at $TargetBinary"
    if ($LegacyFound) { Write-Host "Legacy Python files were preserved under $InstallDir for later cleanup." }
    if (($env:PATH -split ";") -notcontains $BinDir) { Write-Host "Add this folder to PATH to run it by name: $BinDir" }
    Write-Host "Next: codeheart-operating-kit onboard"
} finally {
    if (-not [string]::IsNullOrWhiteSpace($StagingDir)) { Remove-Item -LiteralPath $StagingDir -Recurse -Force -ErrorAction SilentlyContinue }
    Remove-Item -LiteralPath $TempDir -Recurse -Force -ErrorAction SilentlyContinue
}
