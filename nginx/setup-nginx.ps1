# setup-nginx-correct.ps1 - Updated for your Nginx path
param(
    [string]$Action = "install"
)

$ProjectRoot = "D:\DTC-Workspace\netplay"
$NginxPath = "C:\nginx\nginx-1.30.2"  # Your  Nginx path

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "NetPlay - Nginx Setup for Windows" -ForegroundColor Green
Write-Host "Nginx Path: $NginxPath" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan

function Install-NginxConfig {
    Write-Host "`nInstalling Nginx configuration..." -ForegroundColor Yellow
    
    # Check if Nginx exists
    if (-not (Test-Path "$NginxPath\nginx.exe")) {
        Write-Host "Nginx not found at $NginxPath" -ForegroundColor Red
        Write-Host "Please check your Nginx installation path" -ForegroundColor Yellow
        return $false
    }
    
    # Backup original config
    if (Test-Path "$NginxPath\conf\nginx.conf") {
        Copy-Item "$NginxPath\conf\nginx.conf" "$NginxPath\conf\nginx.conf.backup" -Force
        Write-Host "Original config backed up" -ForegroundColor Green
    }
    
    # Copy our config
    $configSource = "$ProjectRoot\nginx\nginx-windows.conf"
    $configDest = "$NginxPath\conf\nginx.conf"
    
    if (-not (Test-Path $configSource)) {
        Write-Host "Config template not found at $configSource" -ForegroundColor Red
        return $false
    }
    
    # Read and update the config
    $configContent = Get-Content $configSource -Raw
    # Update the root path to your frontend
    $frontendPath = "$ProjectRoot\frontend".Replace('\', '/')
    $configContent = $configContent -replace "D:/DTC-Workspace/netplay/frontend", $frontendPath
    
    # Save the config
    $configContent | Out-File -FilePath $configDest -Encoding ASCII
    Write-Host "Nginx configuration installed" -ForegroundColor Green
    Write-Host "Config location: $configDest" -ForegroundColor Gray
    
    return $true
}

function Start-Nginx {
    Write-Host "`nStarting Nginx..." -ForegroundColor Yellow
    
    # Check if already running
    $nginxProcess = Get-Process -Name "nginx" -ErrorAction SilentlyContinue
    if ($nginxProcess) {
        Write-Host "Nginx is already running" -ForegroundColor Yellow
        return $true
    }
    
    # Start Nginx
    cd $NginxPath
    Start-Process -FilePath "$NginxPath\nginx.exe" -WorkingDirectory $NginxPath -WindowStyle Normal
    
    Start-Sleep -Seconds 2
    
    # Verify it's running
    $nginxProcess = Get-Process -Name "nginx" -ErrorAction SilentlyContinue
    if ($nginxProcess) {
        Write-Host "Nginx started successfully" -ForegroundColor Green
        Write-Host "   PID: $($nginxProcess.Id)" -ForegroundColor Gray
        return $true
    } else {
        Write-Host "Failed to start Nginx" -ForegroundColor Red
        Write-Host "Check error log: $NginxPath\logs\error.log" -ForegroundColor Yellow
        return $false
    }
}

function Stop-Nginx {
    Write-Host "`nStopping Nginx..." -ForegroundColor Yellow
    
    if (Test-Path "$NginxPath\nginx.exe") {
        cd $NginxPath
        .\nginx.exe -s stop 2>$null
        Start-Sleep -Seconds 1
        Write-Host "Nginx stopped" -ForegroundColor Green
    }
}

function Test-Nginx {
    Write-Host "`nTesting Nginx configuration..." -ForegroundColor Yellow
    
    cd $NginxPath
    .\nginx.exe -t
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Configuration test passed" -ForegroundColor Green
        return $true
    } else {
        Write-Host "Configuration test failed" -ForegroundColor Red
        return $false
    }
}

# Main execution
switch ($Action.ToLower()) {
    "install" { 
        Install-NginxConfig
        Test-Nginx
    }
    "start" { 
        Test-Nginx
        if ($?) { Start-Nginx }
    }
    "stop" { Stop-Nginx }
    "restart" { 
        Stop-Nginx
        Start-Sleep -Seconds 2
        Start-Nginx
    }
    "test" { Test-Nginx }
    default {
        Write-Host "Usage: .\setup-nginx-correct.ps1 -Action [install|start|stop|restart|test]" -ForegroundColor Yellow
    }
}