$AudioDir = "internal\quote\data\audio"

$EnvFile = Join-Path $PSScriptRoot ".env"
if (Test-Path $EnvFile) {
    Get-Content $EnvFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            [System.Environment]::SetEnvironmentVariable($Matches[1].Trim(), $Matches[2].Trim(), 'Process')
        }
    }
}

$ZipSource = $env:VOICE_ZIP_URL
if (-not $ZipSource) {
    Write-Error "VOICE_ZIP_URL is not set. Create a .env file with VOICE_ZIP_URL=<url or path>"
    exit 1
}

if (Test-Path $AudioDir) {
    Write-Output "Audio directory already exists at $AudioDir, skipping download."
    exit 0
}

$TmpDir = "$env:TEMP\voice"

if (Test-Path $ZipSource) {
    Write-Output "Extracting from local file: $ZipSource"
    New-Item -ItemType Directory -Force -Path "internal\quote\data" | Out-Null
    Expand-Archive -Path $ZipSource -DestinationPath $TmpDir
    Move-Item -Path "$TmpDir\voice" -Destination $AudioDir
    Remove-Item -Recurse -Force $TmpDir
} else {
    $TmpZip = "$env:TEMP\voice.zip"
    Write-Output "Downloading voice files..."
    Invoke-WebRequest -Uri $ZipSource -OutFile $TmpZip
    Write-Output "Extracting..."
    New-Item -ItemType Directory -Force -Path "internal\quote\data" | Out-Null
    Expand-Archive -Path $TmpZip -DestinationPath $TmpDir
    Move-Item -Path "$TmpDir\voice" -Destination $AudioDir
    Remove-Item -Recurse -Force $TmpZip, $TmpDir
}

Write-Output "Done. Audio files extracted to $AudioDir"
