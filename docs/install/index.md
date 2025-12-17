# Installing f5xcctl

Get up and running with the F5 Distributed Cloud CLI in minutes.

## Choose Your Platform

<div class="features-grid">
  <a href="homebrew/" class="feature-card">
    <div class="feature-card-header">
      <div class="feature-icon">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 3L2 12h3v8h14v-8h3L12 3m0 2.7L18 11v8h-2v-6h-2v6h-4v-6H8v6H6v-8l6-5.3z"/></svg>
      </div>
      <h3>Homebrew</h3>
    </div>
    <p>Recommended for macOS. Simple installation and updates with brew.</p>
  </a>

  <a href="script/" class="feature-card">
    <div class="feature-card-header">
      <div class="feature-icon">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M20,19V7H4V19H20M20,3A2,2 0 0,1 22,5V19A2,2 0 0,1 20,21H4A2,2 0 0,1 2,19V5C2,3.89 2.9,3 4,3H20M13,17V15H18V17H13M9.58,13L5.57,9H8.4L11.7,12.3C12.09,12.69 12.09,13.33 11.7,13.72L8.42,17H5.59L9.58,13Z"/></svg>
      </div>
      <h3>Install Script</h3>
    </div>
    <p>Universal installer for Linux and macOS. Works with or without sudo.</p>
  </a>

  <a href="windows/" class="feature-card">
    <div class="feature-card-header">
      <div class="feature-icon">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M3,12V6.75L9,5.43V11.91L3,12M20,3V11.75L10,11.9V5.21L20,3M3,13L9,13.09V19.9L3,18.75V13M20,13.25V22L10,20.09V13.1L20,13.25Z"/></svg>
      </div>
      <h3>Windows</h3>
    </div>
    <p>Manual download for Windows. PowerShell completions included.</p>
  </a>

  <a href="source/" class="feature-card">
    <div class="feature-card-header">
      <div class="feature-icon">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12,2A3,3 0 0,1 15,5V11A3,3 0 0,1 12,14A3,3 0 0,1 9,11V5A3,3 0 0,1 12,2M19,11C19,14.53 16.39,17.44 13,17.93V21H11V17.93C7.61,17.44 5,14.53 5,11H7A5,5 0 0,0 12,16A5,5 0 0,0 17,11H19Z"/></svg>
      </div>
      <h3>Build from Source</h3>
    </div>
    <p>For development or latest unreleased features. Requires Go 1.21+.</p>
  </a>
</div>

## Quick Start

The fastest way to get started on macOS:

```bash
brew tap robinmordasiewicz/f5xcctl
brew install --cask f5xcctl
```

Or use the universal install script on any platform:

```bash
curl -fsSL https://robinmordasiewicz.github.io/f5xcctl/install.sh | sh
```

## Post-Installation Setup

<div class="features-grid">
  <a href="authentication/" class="feature-card">
    <div class="feature-card-header">
      <div class="feature-icon">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12,17A2,2 0 0,0 14,15C14,13.89 13.1,13 12,13A2,2 0 0,0 10,15A2,2 0 0,0 12,17M18,8A2,2 0 0,1 20,10V20A2,2 0 0,1 18,22H6A2,2 0 0,1 4,20V10C4,8.89 4.9,8 6,8H7V6A5,5 0 0,1 12,1A5,5 0 0,1 17,6V8H18M12,3A3,3 0 0,0 9,6V8H15V6A3,3 0 0,0 12,3Z"/></svg>
      </div>
      <h3>Authentication</h3>
    </div>
    <p>Configure API tokens or certificate-based authentication to connect to F5 XC.</p>
  </a>

  <a href="completions/" class="feature-card">
    <div class="feature-card-header">
      <div class="feature-icon">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12.89,3L14.85,3.4L11.11,21L9.15,20.6L12.89,3M19.59,12L16,8.41V5.58L22.42,12L16,18.41V15.58L19.59,12M1.58,12L8,5.58V8.41L4.41,12L8,15.58V18.41L1.58,12Z"/></svg>
      </div>
      <h3>Shell Completions</h3>
    </div>
    <p>Enable tab completion for bash, zsh, fish, or PowerShell.</p>
  </a>

  <a href="environment-variables/" class="feature-card">
    <div class="feature-card-header">
      <div class="feature-icon">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12,15.5A3.5,3.5 0 0,1 8.5,12A3.5,3.5 0 0,1 12,8.5A3.5,3.5 0 0,1 15.5,12A3.5,3.5 0 0,1 12,15.5M19.43,12.97C19.47,12.65 19.5,12.33 19.5,12C19.5,11.67 19.47,11.34 19.43,11L21.54,9.37C21.73,9.22 21.78,8.95 21.66,8.73L19.66,5.27C19.54,5.05 19.27,4.96 19.05,5.05L16.56,6.05C16.04,5.66 15.5,5.32 14.87,5.07L14.5,2.42C14.46,2.18 14.25,2 14,2H10C9.75,2 9.54,2.18 9.5,2.42L9.13,5.07C8.5,5.32 7.96,5.66 7.44,6.05L4.95,5.05C4.73,4.96 4.46,5.05 4.34,5.27L2.34,8.73C2.21,8.95 2.27,9.22 2.46,9.37L4.57,11C4.53,11.34 4.5,11.67 4.5,12C4.5,12.33 4.53,12.65 4.57,12.97L2.46,14.63C2.27,14.78 2.21,15.05 2.34,15.27L4.34,18.73C4.46,18.95 4.73,19.03 4.95,18.95L7.44,17.94C7.96,18.34 8.5,18.68 9.13,18.93L9.5,21.58C9.54,21.82 9.75,22 10,22H14C14.25,22 14.46,21.82 14.5,21.58L14.87,18.93C15.5,18.67 16.04,18.34 16.56,17.94L19.05,18.95C19.27,19.03 19.54,18.95 19.66,18.73L21.66,15.27C21.78,15.05 21.73,14.78 21.54,14.63L19.43,12.97Z"/></svg>
      </div>
      <h3>Environment Variables</h3>
    </div>
    <p>Configure f5xcctl behavior through environment variables.</p>
  </a>
</div>

## Verify Installation

After installing, verify that f5xcctl is working:

```bash
f5xcctl version
```

## Next Steps

Once installed and authenticated, explore the [Commands](../commands/index.md) documentation or follow the [Guides](../guides/index.md) to deploy your first load balancer.
