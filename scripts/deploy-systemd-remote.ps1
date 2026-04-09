param(
  [Parameter(Mandatory = $true)]
  [Alias("Host")]
  [string]$RemoteHost,

  [string]$User = "root",

  [string]$Password,

  [string]$RemoteRoot = "/opt/opsnexus-remote-test",

  [string]$RemoteEnvDir = "/etc/opsnexus",

  [switch]$IncludeWeb,

  [switch]$ConfigureNginx,

  [switch]$UploadExamples
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$apiBinary = Join-Path $repoRoot "api/opsnexus-api-linux-amd64"
$webDist = Join-Path $repoRoot "web/dist"
$serviceFile = Join-Path $repoRoot "deploy/systemd/opsnexus-api.service"
$configExample = Join-Path $repoRoot "deploy/systemd/opsnexus-api.config.yaml.example"
$envExample = Join-Path $repoRoot "deploy/systemd/opsnexus-api.env.example"
$remoteNginxConfig = Join-Path $repoRoot "opsnexus-remote-nginx.conf"
$usePoshSsh = -not [string]::IsNullOrWhiteSpace($Password)
$sshCredential = $null
$webArchive = Join-Path $repoRoot "web/dist-upload.tar.gz"

if (-not (Test-Path $apiBinary)) {
  throw "Missing backend binary: $apiBinary. Build the Linux binary in the api directory first."
}

if ($IncludeWeb -and -not (Test-Path $webDist)) {
  throw "Missing frontend dist directory: $webDist. Run npm run build in the web directory first."
}

if ($ConfigureNginx -and -not $IncludeWeb) {
  throw "ConfigureNginx requires IncludeWeb so the frontend assets are uploaded together."
}

if ($ConfigureNginx -and -not (Test-Path $remoteNginxConfig)) {
  throw "Missing nginx config template: $remoteNginxConfig"
}

if ($usePoshSsh) {
  Import-Module Posh-SSH -ErrorAction Stop
  $securePassword = ConvertTo-SecureString $Password -AsPlainText -Force
  $sshCredential = New-Object System.Management.Automation.PSCredential($User, $securePassword)
}

$remote = "$User@$RemoteHost"
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"

function Copy-ToRemote {
  param(
    [Parameter(Mandatory = $true)]
    [string]$LocalPath,

    [Parameter(Mandatory = $true)]
    [string]$RemoteDirectory,

    [Parameter(Mandatory = $true)]
    [string]$RemoteName
  )

  if ($usePoshSsh) {
    Set-SCPItem `
      -ComputerName $RemoteHost `
      -Credential $sshCredential `
      -AcceptKey `
      -Path $LocalPath `
      -Destination $RemoteDirectory `
      -NewName $RemoteName `
      -Force | Out-Null
    return
  }

  & scp $LocalPath "${remote}:$RemoteDirectory/$RemoteName"
  if ($LASTEXITCODE -ne 0) {
    throw "scp upload failed: $LocalPath -> $RemoteDirectory/$RemoteName"
  }
}

function Invoke-RemoteScript {
  param(
    [Parameter(Mandatory = $true)]
    [string]$Command
  )

  if ($usePoshSsh) {
    $session = $null
    try {
      $session = New-SSHSession -ComputerName $RemoteHost -Credential $sshCredential -AcceptKey -ConnectionTimeout 15
      $result = Invoke-SSHCommand -SessionId $session.SessionId -Command $Command -TimeOut 600

      if ($result.Output) {
        $result.Output | ForEach-Object { Write-Host $_ }
      }

      if ($result.Error) {
        $result.Error | ForEach-Object { Write-Host $_ -ForegroundColor Yellow }
      }

      if ($result.ExitStatus -ne 0) {
        throw "Remote script failed with exit code: $($result.ExitStatus)"
      }
    }
    finally {
      if ($session) {
        Remove-SSHSession -SessionId $session.SessionId | Out-Null
      }
    }
    return
  }

  & ssh $remote $Command
  if ($LASTEXITCODE -ne 0) {
    throw "Remote script failed with exit code: $LASTEXITCODE"
  }
}

Write-Host "Uploading backend artifacts to $remote ..." -ForegroundColor Cyan

Copy-ToRemote -LocalPath $apiBinary -RemoteDirectory "/tmp" -RemoteName "opsnexus-api-linux-amd64.new"
Copy-ToRemote -LocalPath $serviceFile -RemoteDirectory "/tmp" -RemoteName "opsnexus-api.service"

if ($UploadExamples) {
  Copy-ToRemote -LocalPath $configExample -RemoteDirectory "/tmp" -RemoteName "opsnexus-api.config.yaml.example"
  Copy-ToRemote -LocalPath $envExample -RemoteDirectory "/tmp" -RemoteName "opsnexus-api.env.example"
}

if ($IncludeWeb) {
  if (Test-Path $webArchive) {
    Remove-Item $webArchive -Force
  }
  & tar -czf $webArchive -C $webDist .
  if ($LASTEXITCODE -ne 0) {
    throw "Failed to create frontend archive: $webArchive"
  }
  Copy-ToRemote -LocalPath $webArchive -RemoteDirectory "/tmp" -RemoteName "opsnexus-web-dist.tar.gz"

  if ($ConfigureNginx) {
    Copy-ToRemote -LocalPath $remoteNginxConfig -RemoteDirectory "/tmp" -RemoteName "opsnexus-remote-nginx.conf"
  }
}

$remoteScript = @"
set -e
mkdir -p $RemoteRoot
mkdir -p $RemoteEnvDir

cp $RemoteRoot/opsnexus-api-linux-amd64 $RemoteRoot/opsnexus-api-linux-amd64.bak-$timestamp 2>/dev/null || true
install -m 755 /tmp/opsnexus-api-linux-amd64.new $RemoteRoot/opsnexus-api-linux-amd64
install -m 644 /tmp/opsnexus-api.service /etc/systemd/system/opsnexus-api.service

if [ -f /tmp/opsnexus-api.config.yaml.example ]; then
  install -m 640 /tmp/opsnexus-api.config.yaml.example $RemoteRoot/config.yaml.example
fi

if [ -f /tmp/opsnexus-api.env.example ]; then
  install -m 640 /tmp/opsnexus-api.env.example $RemoteEnvDir/opsnexus-api.env.example
fi

if [ -f /tmp/opsnexus-web-dist.tar.gz ]; then
  rm -rf $RemoteRoot/web-dist.deploy-$timestamp
  mkdir -p $RemoteRoot/web-dist.deploy-$timestamp
  tar -xzf /tmp/opsnexus-web-dist.tar.gz -C $RemoteRoot/web-dist.deploy-$timestamp
  if [ ! -f $RemoteRoot/web-dist.deploy-$timestamp/index.html ]; then
    echo "Frontend archive extraction failed"
    exit 1
  fi
  mv $RemoteRoot/web-dist $RemoteRoot/web-dist.bak-$timestamp 2>/dev/null || true
  mv $RemoteRoot/web-dist.deploy-$timestamp $RemoteRoot/web-dist
  rm -f /tmp/opsnexus-web-dist.tar.gz
fi

if [ -f /tmp/opsnexus-remote-nginx.conf ]; then
  export DEBIAN_FRONTEND=noninteractive
  apt-get update
  apt-get install -y nginx
  rm -rf /usr/share/nginx/html/*
  cp -a $RemoteRoot/web-dist/. /usr/share/nginx/html/
  install -m 644 /tmp/opsnexus-remote-nginx.conf /etc/nginx/sites-available/opsnexus-remote.conf
  ln -sfn /etc/nginx/sites-available/opsnexus-remote.conf /etc/nginx/sites-enabled/opsnexus-remote.conf
  rm -f /etc/nginx/sites-enabled/default
  nginx -t
  systemctl enable --now nginx
  systemctl reload nginx
fi

systemctl daemon-reload
systemctl enable --now opsnexus-api.service
systemctl status opsnexus-api.service --no-pager -l
"@

Write-Host "Running remote deployment script ..." -ForegroundColor Cyan
Invoke-RemoteScript -Command $remoteScript

Write-Host ""
Write-Host "Done. Verify these items on the remote host if needed:" -ForegroundColor Green
Write-Host "1. $RemoteRoot/config.yaml contains the target environment settings"
Write-Host "2. $RemoteEnvDir/opsnexus-api.env exists and contains real secrets"
Write-Host "3. systemctl status opsnexus-api.service"
if ($ConfigureNginx) {
  Write-Host "4. http://${RemoteHost}:8080/"
  Write-Host "5. http://${RemoteHost}:8080/api/v1/captcha"
  Write-Host "6. systemctl status nginx"
}

if (Test-Path $webArchive) {
  Remove-Item $webArchive -Force
}
