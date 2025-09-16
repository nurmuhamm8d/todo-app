# Stop any running processes
Write-Host "Stopping any running processes..."
taskkill /F /IM todo-app-dev.exe 2>$null

# Clean up previous build
Write-Host "Cleaning up previous build..."
if (Test-Path -Path "frontend\dist") {
    Remove-Item -Recurse -Force "frontend\dist"
}

# Install frontend dependencies
Write-Host "Installing frontend dependencies..."
Set-Location frontend
npm install

# Build frontend
Write-Host "Building frontend..."
npm run build

# Go back to backend and run
Set-Location ..\backend
Write-Host "Starting Wails dev server..."
wails dev
