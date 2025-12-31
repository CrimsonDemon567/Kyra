$version = "1.0.0"
$installDir = "$env:ProgramFiles\Kyra"

Write-Host "Installing Kyra SDK $version..."

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "x86" }

# Create install directory
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

# Copy binaries
Copy-Item "kyra-$version-windows-$arch.exe" "$installDir\kyra.exe" -Force
Copy-Item "kyrac-$version-windows-$arch.exe" "$installDir\kyrac.exe" -Force

# Add to PATH if missing
$path = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($path -notlike "*$installDir*") {
    Write-Host "Adding Kyra to PATH..."
    $newPath = "$path;$installDir"
    setx /M PATH "$newPath" | Out-Null
}

Write-Host "Kyra SDK installed successfully."
Write-Host "You can now run:"
Write-Host "  kyra -kbc file.kbc"
Write-Host "  kyrac -kbc file.kyra"
