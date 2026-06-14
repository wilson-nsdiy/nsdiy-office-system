#Requires -Version 5.1
param(
    [Parameter(Position=0)]
    [string]$Command
)

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$ProjectRoot = Resolve-Path (Join-Path $ScriptDir "..")

$Version = try { git -C $ProjectRoot describe --tags --always --dirty 2>$null } catch { "dev" }

function Show-Usage {
    Write-Host "Usage: $($MyInvocation.ScriptName) <command>" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  build           - Build frontend + backend (without embedding frontend)"
    Write-Host "  build-embed     - Build frontend and embed into backend binary (production)"
    Write-Host "  build-backend   - Build Go backend (without embed)"
    Write-Host "  build-frontend  - Build Next.js frontend (static export)"
    Write-Host "  test            - Run all tests"
    Write-Host "  test-backend    - Run backend tests"
    Write-Host "  test-frontend   - Run frontend tests"
    Write-Host "  lint            - Run linters"
    Write-Host "  dev-backend     - Run backend in development mode"
    Write-Host "  dev-frontend    - Run frontend in development mode"
    Write-Host "  docker-up       - Start all services with Docker Compose"
    Write-Host "  docker-down     - Stop all services"
    Write-Host "  clean           - Clean build artifacts"
    Write-Host "  generate        - Generate Ent code"
    Write-Host "  migrate-new     - Create a new migration file"
}

function Invoke-BuildBackend {
    Push-Location "$ProjectRoot\backend"
    try {
        $env:CGO_ENABLED = "0"
        go build -ldflags="-s -w -X main.Version=$Version" -trimpath -o bin\server.exe .\cmd\server
    } finally {
        Pop-Location
        Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
    }
}

function Invoke-BuildFrontend {
    Push-Location "$ProjectRoot\frontend"
    try {
        npm run build
    } finally {
        Pop-Location
    }
}

function Invoke-Build {
    Invoke-BuildFrontend
    Invoke-BuildBackend
}

function Invoke-BuildEmbed {
    Invoke-BuildFrontend
    Write-Host "Copying frontend output to backend embed directory..." -ForegroundColor Cyan
    $distDir = "$ProjectRoot\backend\internal\web\dist"
    Get-ChildItem "$distDir\*" -ErrorAction SilentlyContinue | Where-Object { $_.Name -ne ".gitkeep" } | Remove-Item -Recurse -Force
    Copy-Item "$ProjectRoot\frontend\out\*" "$distDir\" -Recurse -Force
    Invoke-BuildBackend
    Write-Host "Embedded build complete." -ForegroundColor Green
}

function Invoke-TestBackend {
    Push-Location "$ProjectRoot\backend"
    try {
        go test ./...
    } finally {
        Pop-Location
    }
}

function Invoke-TestFrontend {
    Push-Location "$ProjectRoot\frontend"
    try {
        npm run lint:check
        npx tsc --noEmit
    } finally {
        Pop-Location
    }
}

function Invoke-Test {
    Invoke-TestBackend
    Invoke-TestFrontend
}

function Invoke-Lint {
    Push-Location "$ProjectRoot\backend"
    try {
        golangci-lint run ./...
    } finally {
        Pop-Location
    }
}

function Invoke-DevBackend {
    Push-Location "$ProjectRoot\backend"
    try {
        go run .\cmd\server
    } finally {
        Pop-Location
    }
}

function Invoke-DevFrontend {
    Push-Location "$ProjectRoot\frontend"
    try {
        npm run dev
    } finally {
        Pop-Location
    }
}

function Invoke-DockerUp {
    Push-Location $ScriptDir
    try {
        docker-compose up -d
    } finally {
        Pop-Location
    }
}

function Invoke-DockerDown {
    Push-Location $ScriptDir
    try {
        docker-compose down
    } finally {
        Pop-Location
    }
}

function Invoke-Clean {
    $binDir = "$ProjectRoot\backend\bin"
    if (Test-Path $binDir) { Remove-Item $binDir -Recurse -Force }

    $nextDir = "$ProjectRoot\frontend\.next"
    if (Test-Path $nextDir) { Remove-Item $nextDir -Recurse -Force }
}

function Invoke-Generate {
    Push-Location "$ProjectRoot\backend"
    try {
        go generate .\ent
        go generate .\cmd\server
    } finally {
        Pop-Location
    }
}

function Invoke-MigrateNew {
    $name = Read-Host "Enter migration name"
    $ts = Get-Date -Format "yyyyMMddHHmmss"
    $file = "$ProjectRoot\backend\migrations\${ts}_${name}.sql"
    New-Item -ItemType File -Path $file -Force | Out-Null
    Write-Host "Created $file" -ForegroundColor Green
}

# Main
if (-not $Command) {
    Show-Usage
    exit 1
}

switch ($Command) {
    "build"         { Invoke-Build }
    "build-embed"   { Invoke-BuildEmbed }
    "build-backend" { Invoke-BuildBackend }
    "build-frontend"{ Invoke-BuildFrontend }
    "test"          { Invoke-Test }
    "test-backend"  { Invoke-TestBackend }
    "test-frontend" { Invoke-TestFrontend }
    "lint"          { Invoke-Lint }
    "dev-backend"   { Invoke-DevBackend }
    "dev-frontend"  { Invoke-DevFrontend }
    "docker-up"     { Invoke-DockerUp }
    "docker-down"   { Invoke-DockerDown }
    "clean"         { Invoke-Clean }
    "generate"      { Invoke-Generate }
    "migrate-new"   { Invoke-MigrateNew }
    { $_ -in "help", "--help", "-h" } { Show-Usage }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Show-Usage
        exit 1
    }
}
