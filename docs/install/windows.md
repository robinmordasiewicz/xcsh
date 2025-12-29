# Windows

Install xcsh on Windows by downloading the binary directly.

## Download

1. Go to the [GitHub Releases](https://github.com/robinmordasiewicz/f5xc-xcsh/releases) page
2. Download `xcsh_windows_amd64.zip` (or `xcsh_windows_arm64.zip` for ARM)
3. Extract the archive

## Installation

### Option 1: Add to User PATH

1. Create a directory for the binary:

   ```powershell
   mkdir $env:USERPROFILE\bin
   ```

2. Move the extracted `xcsh.exe` to this directory

3. Add to your PATH:

   ```powershell
   $env:PATH += ";$env:USERPROFILE\bin"
   [Environment]::SetEnvironmentVariable("PATH", $env:PATH, [EnvironmentVariableTarget]::User)
   ```

### Option 2: System-Wide Installation

Move `xcsh.exe` to a directory already in your system PATH, such as `C:\Windows\System32`.

## Verify Installation

```powershell
xcsh version
```

## PowerShell Completions

Enable tab completion in PowerShell:

```powershell
xcsh completion powershell | Out-String | Invoke-Expression
```

To make completions permanent, add the above command to your PowerShell profile:

```powershell
# Open your profile for editing
notepad $PROFILE

# Add this line:
xcsh completion powershell | Out-String | Invoke-Expression
```

## WSL Alternative

If you prefer a Unix-like environment, you can use Windows Subsystem for Linux (WSL) and follow the [Script](script.md) installation method:

```bash
curl -fsSL https://robinmordasiewicz.github.io/f5xc-xcsh/install.sh | sh
```

## Troubleshooting

### Command Not Found

Ensure the directory containing `xcsh.exe` is in your PATH:

```powershell
echo $env:PATH
```

### Execution Policy

If PowerShell blocks the script, adjust the execution policy:

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```
