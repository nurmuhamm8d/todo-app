# Stop any running processes
Write-Host "Stopping any running instances..." -ForegroundColor Yellow
Get-Process -Name "todo-app*" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue

# Clean up
Write-Host "Cleaning up..." -ForegroundColor Yellow
if (Test-Path "frontend/dist") {
    Remove-Item -Recurse -Force frontend/dist
}
if (Test-Path "frontend/node_modules") {
    Remove-Item -Recurse -Force frontend/node_modules
}
if (Test-Path "frontend/.wails") {
    Remove-Item -Recurse -Force frontend/.wails
}
if (Test-Path "frontend/package-lock.json") {
    Remove-Item -Force frontend/package-lock.json
}
if (Test-Path "frontend/.svelte-kit") {
    Remove-Item -Recurse -Force frontend/.svelte-kit
}
if (Test-Path "frontend/.vite") {
    Remove-Item -Recurse -Force frontend/.vite
}

# Install dependencies
Write-Host "Installing dependencies..." -ForegroundColor Yellow
Push-Location frontend

try {
    # Clear npm cache
    Write-Host "Clearing npm cache..." -ForegroundColor Yellow
    npm cache clean --force

    # Install Node.js version if nvm is available
    if (Get-Command nvm -ErrorAction SilentlyContinue) {
        Write-Host "Setting Node.js version..." -ForegroundColor Yellow
        nvm use
    }

    # Install dependencies
    Write-Host "Running npm install..." -ForegroundColor Yellow
    npm install --legacy-peer-deps --force
    if ($LASTEXITCODE -ne 0) { 
        throw "Failed to install dependencies"
    }

    # Build frontend
    Write-Host "Building frontend..." -ForegroundColor Yellow
    npm run build
    if ($LASTEXITCODE -ne 0) { 
        throw "Frontend build failed"
    }
}
catch {
    Write-Host "Error: $_" -ForegroundColor Red
    Pop-Location
    exit 1
}

Pop-Location

# Clean and rebuild Wails
Write-Host "Cleaning Wails build..." -ForegroundColor Yellow
wails clean

# Install Wails dependencies
Write-Host "Installing Wails dependencies..." -ForegroundColor Yellow
wails install

# Run Wails
Write-Host "Starting Wails development server..." -ForegroundColor Green
wails dev -v
