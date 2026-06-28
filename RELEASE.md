# Release Guide

### 1. Publish the Helm chart

The release script updates `Chart.yaml` (`version`, `appVersion`) and `values.yaml` (`app.image.tag`) to the new version, then packages and pushes to GHCR.

```bash
./charts/core/release.sh \
  -u <github-username> \
  -t <github-pat> \
  -v <version>
```

The GitHub PAT requires the `write:packages` scope.

### 2. Publish the npm package

The release script updates `package.json` to the new version, builds the SDK, then publishes to npm.

```bash
./sdks/node/release.sh \
  -t <npm-token> \
  -v <version>
```

If your npm account has 2FA enabled, pass the one-time code:

```bash
./sdks/node/release.sh \
  -t <npm-token> \
  -v <version> \
  --otp <code>
```

Alternatively, set `NPM_TOKEN` in your environment to avoid passing `-t` on every run:

```bash
export NPM_TOKEN=<npm-token>
./sdks/node/release.sh -v <version>
```

### 3. Commit the version bumps

Both scripts modify tracked files. Commit them before tagging:

```bash
git add .
git commit -m "chore: bump version to <version>"
```

### 4. Push the Git tag

Pushing a tag triggers all GitHub Actions workflows, which build and push the Docker images tagged with the version.

```bash
git tag <version>
git push origin <version>
```

---

## Version checklist

Before tagging, confirm these are all set to the same version:

- [ ] `charts/core/Chart.yaml` -- `version` and `appVersion`
- [ ] `charts/core/values.yaml` -- `app.image.tag`
- [ ] `sdks/node/package.json` -- `version`

The release scripts in steps 1 and 2 handle these automatically.
