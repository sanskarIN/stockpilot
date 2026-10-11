$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$RootDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$BackupDir = if ($env:BACKUP_DIR) { $env:BACKUP_DIR } else { Join-Path $RootDir "backups" }
$RetentionDays = if ($env:BACKUP_RETENTION_DAYS) { $env:BACKUP_RETENTION_DAYS } else { "14" }
$PostgresUser = if ($env:POSTGRES_USER) { $env:POSTGRES_USER } else { "stockpilot" }
$PostgresDb = if ($env:POSTGRES_DB) { $env:POSTGRES_DB } else { "stockpilot" }

$ParsedRetention = 0
if (-not [int]::TryParse($RetentionDays, [ref]$ParsedRetention) -or $ParsedRetention -lt 0) {
    throw "BACKUP_RETENTION_DAYS must be a non-negative integer."
}

$Timestamp = [DateTime]::UtcNow.ToString("yyyyMMdd'T'HHmmss'Z'")
$FinalPath = Join-Path $BackupDir "stockpilot-$Timestamp.dump"
$TempPath = "$FinalPath.tmp"

New-Item -ItemType Directory -Path $BackupDir -Force | Out-Null
Remove-Item -LiteralPath $TempPath -Force -ErrorAction SilentlyContinue

function Invoke-CheckedCommand {
    param(
        [Parameter(Mandatory = $true)][string]$Executable,
        [Parameter(Mandatory = $true)][string[]]$Arguments
    )

    & $Executable @Arguments
    if ($LASTEXITCODE -ne 0) {
        throw "$Executable exited with code $LASTEXITCODE."
    }
}

try {
    $PgDump = Get-Command pg_dump -ErrorAction SilentlyContinue
    if ($PgDump -and $env:DATABASE_URL) {
        Write-Host "Creating StockPilot backup with local pg_dump..."
        Invoke-CheckedCommand -Executable $PgDump.Source -Arguments @(
            "--format=custom",
            "--no-owner",
            "--no-privileges",
            "--file=$TempPath",
            $env:DATABASE_URL
        )
    }
    else {
        $Docker = Get-Command docker -ErrorAction SilentlyContinue
        if (-not $Docker) {
            throw "Install pg_dump and set DATABASE_URL, or install Docker and start the Compose db service."
        }

        & $Docker.Source compose version *> $null
        if ($LASTEXITCODE -ne 0) {
            throw "Docker Compose is not available."
        }

        $RunningServices = & $Docker.Source compose ps --status running --services
        if ($LASTEXITCODE -ne 0 -or $RunningServices -notcontains "db") {
            throw "The Docker Compose db service is not running."
        }

        Write-Host "Creating StockPilot backup through the Docker Compose database service..."
        $RemotePath = "/tmp/stockpilot-$Timestamp.dump"
        try {
            Invoke-CheckedCommand -Executable $Docker.Source -Arguments @(
                "compose", "exec", "-T", "db", "pg_dump",
                "-U", $PostgresUser,
                "-d", $PostgresDb,
                "--format=custom",
                "--no-owner",
                "--no-privileges",
                "--file=$RemotePath"
            )
            Invoke-CheckedCommand -Executable $Docker.Source -Arguments @(
                "compose", "cp", "db:$RemotePath", $TempPath
            )
        }
        finally {
            & $Docker.Source compose exec -T db rm -f $RemotePath *> $null
        }
    }

    if (-not (Test-Path -LiteralPath $TempPath) -or (Get-Item -LiteralPath $TempPath).Length -le 0) {
        throw "Backup command completed without producing a non-empty dump."
    }

    Move-Item -LiteralPath $TempPath -Destination $FinalPath -Force

    if ($ParsedRetention -gt 0) {
        $Cutoff = [DateTime]::UtcNow.AddDays(-$ParsedRetention)
        Get-ChildItem -LiteralPath $BackupDir -Filter "stockpilot-*.dump" -File |
            Where-Object { $_.LastWriteTimeUtc -lt $Cutoff } |
            Remove-Item -Force
    }

    Write-Host "StockPilot backup created: $FinalPath"
}
finally {
    Remove-Item -LiteralPath $TempPath -Force -ErrorAction SilentlyContinue
}
