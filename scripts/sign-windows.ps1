param(
    [Parameter(Mandatory=$true)][string]$File
)
$ErrorActionPreference = 'Stop'
# mygo build runs this as windows.signCommand on each Windows executable it
# ships, before the installer packs them: the app, the bundled Agari engine,
# the installer and its uninstaller. The app gets its icon first, since
# changing resources afterwards would break a signature.
$path = (Resolve-Path -LiteralPath $File).Path
if ([IO.Path]::GetFileName($path) -eq 'Riichi Studio.exe') {
    & "$PSScriptRoot/apply-icon.ps1" -Executable $path
}

# Signing is optional: without a certificate, builds are unsigned and
# SmartScreen warns on first run.
if (-not $env:WINDOWS_CERTIFICATE_FILE) { exit 0 }
$signtool = Get-Command signtool.exe -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source
if (-not $signtool) {
    $signtool = Get-ChildItem "${env:ProgramFiles(x86)}\Windows Kits\10\bin\*\x64\signtool.exe" -ErrorAction SilentlyContinue |
        Sort-Object FullName -Descending | Select-Object -First 1 -ExpandProperty FullName
}
if (-not $signtool) { throw 'Signing needs signtool from the Windows SDK.' }
$arguments = @('sign', '/f', $env:WINDOWS_CERTIFICATE_FILE, '/fd', 'sha256', '/tr', 'http://timestamp.digicert.com', '/td', 'sha256')
if ($env:WINDOWS_CERTIFICATE_PASSWORD) { $arguments += @('/p', $env:WINDOWS_CERTIFICATE_PASSWORD) }
& $signtool @arguments $path
if ($LASTEXITCODE -ne 0) { throw "signtool failed with exit code $LASTEXITCODE" }
Write-Output "Signed $path"
