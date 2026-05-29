#!/usr/bin/env node
// Copyright (c) 2026 Amr Elsagaei. All rights reserved.
//
// Thin launcher: resolves the platform-specific @updock/<os>-<arch> package
// (installed as an optional dependency) and execs the bundled binary.

"use strict";

const { execFileSync } = require("node:child_process");

// Map Node's platform/arch to our package naming.
const PLATFORM_PACKAGES = {
  "linux-x64": "@updock/linux-x64",
  "linux-arm64": "@updock/linux-arm64",
  "darwin-x64": "@updock/darwin-x64",
  "darwin-arm64": "@updock/darwin-arm64",
  "win32-x64": "@updock/win32-x64",
};

function binaryPath() {
  const key = `${process.platform}-${process.arch}`;
  const pkg = PLATFORM_PACKAGES[key];

  if (!pkg) {
    throw new Error(
      `updock: unsupported platform ${key}. ` +
        `Download a binary from https://github.com/amrelsagaei/updock/releases`,
    );
  }

  const binName = process.platform === "win32" ? "updock.exe" : "updock";
  try {
    return require.resolve(`${pkg}/bin/${binName}`);
  } catch {
    throw new Error(
      `updock: the platform package ${pkg} is not installed. ` +
        `Reinstall with 'npm install updock', or grab a binary from ` +
        `https://github.com/amrelsagaei/updock/releases`,
    );
  }
}

function main() {
  let bin;
  try {
    bin = binaryPath();
  } catch (err) {
    process.stderr.write(`${err.message}\n`);
    process.exit(1);
  }

  try {
    execFileSync(bin, process.argv.slice(2), { stdio: "inherit" });
  } catch (err) {
    // Propagate the child's exit code; execFileSync throws on non-zero exit.
    process.exit(typeof err.status === "number" ? err.status : 1);
  }
}

main();
