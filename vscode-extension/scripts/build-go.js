const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

const platforms = [
    { os: 'linux', arch: 'amd64' },
    { os: 'linux', arch: 'arm64' },
    { os: 'windows', arch: 'amd64', ext: '.exe' },
    { os: 'darwin', arch: 'amd64' },
    { os: 'darwin', arch: 'arm64' }
];

const builtDir = path.join(__dirname, '..', 'go-cli', 'bin');
const goCliDir = path.join(__dirname, '..', '..', 'go-cli');
const mainFile = path.join(goCliDir, 'main.go');

// Ensure the output directory exists
if (fs.existsSync(builtDir)) {
    fs.rmSync(builtDir, { recursive: true, force: true });
}
fs.mkdirSync(builtDir, { recursive: true });

console.log('Building Go binaries...');

platforms.forEach(platform => {
    const { os, arch, ext } = platform;
    const outputName = `llamit-${os}-${arch}${ext || ''}`;
    const outputPath = path.join(builtDir, outputName);

    console.log(`Building for ${os}/${arch}...`);

    try {
        const env = { ...process.env, GOOS: os, GOARCH: arch };
        // Clean environment is important for cross-compilation if CGO is not needed
        // Assuming CGO_ENABLED=0 to allow easy cross-compilation
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
