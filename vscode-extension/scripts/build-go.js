const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

// Map VSCE_TARGET to platform configuration
function getPlatformForTarget(target) {
    const targetMap = {
        'win32-x64': { os: 'windows', arch: 'amd64', ext: '.exe' },
        'win32-arm64': { os: 'windows', arch: 'arm64', ext: '.exe' },
        'linux-x64': { os: 'linux', arch: 'amd64' },
        'linux-arm64': { os: 'linux', arch: 'arm64' },
        'darwin-x64': { os: 'darwin', arch: 'amd64' },
        'darwin-arm64': { os: 'darwin', arch: 'arm64' }
    };
    return targetMap[target];
}

// Detect if we're building for a specific platform target
const vsceTarget = process.env.VSCE_TARGET;

let platformsToBuild;
if (vsceTarget) {
    const platform = getPlatformForTarget(vsceTarget);
    if (!platform) {
        console.error(`Unknown VSCE_TARGET: ${vsceTarget}`);
        process.exit(1);
    }
    platformsToBuild = [platform];
    console.log(`Building for target: ${vsceTarget}`);
} else {
    // Build all platforms when not targeting a specific one
    platformsToBuild = [
        { os: 'linux', arch: 'amd64' },
        { os: 'linux', arch: 'arm64' },
        { os: 'windows', arch: 'amd64', ext: '.exe' },
        { os: 'darwin', arch: 'amd64' },
        { os: 'darwin', arch: 'arm64' }
    ];
}

const builtDir = path.join(__dirname, '..', 'go-cli', 'bin');
const goCliDir = path.join(__dirname, '..', '..', 'go-cli');
const mainFile = path.join(goCliDir, 'main.go');

// Clean and recreate the output directory
if (fs.existsSync(builtDir)) {
    fs.rmSync(builtDir, { recursive: true, force: true });
}
fs.mkdirSync(builtDir, { recursive: true });

console.log('Building Go binaries...');

platformsToBuild.forEach(platform => {
    const { os, arch, ext } = platform;
    const outputName = `llamit-${os}-${arch}${ext || ''}`;
    const outputPath = path.join(builtDir, outputName);

    console.log(`Building for ${os}/${arch}...`);

    try {
        const env = { ...process.env, GOOS: os, GOARCH: arch };
        env.CGO_ENABLED = '0';

        execSync(`go build -ldflags="-s -w" -o "${outputPath}" "${mainFile}"`, {
            env,
            cwd: goCliDir,
            stdio: 'inherit'
        });
        console.log(`✓ Built ${outputName}`);
    } catch (error) {
        console.error(`✗ Failed to build for ${os}/${arch}:`, error.message);
        process.exit(1);
    }
});

console.log('All binaries built successfully.');
