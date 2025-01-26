[![Build + Test](https://github.com/peter-mghendi/cloakroom/actions/workflows/build.yml/badge.svg)](https://github.com/peter-mghendi/cloakroom/actions/workflows/build.yml)

# Cloakroom
> **Minimal plugin manager for Keycloak**  
> Centralize, streamline, and automate your Keycloak plugin management with Cloakroom.

---

## TL;DR

Cloakroom transforms your plugin installation process from a **long series of ADD lines** into a **single restore step**. 
Instead of cluttering your Dockerfile (or bash scripts) with multiple calls to fetch JARs from GitHub Releases, you simply define your plugins in Cloakroom’s manifest. 
Cloakroom then downloads them all in one go.

Before (manual plugin additions):

```dockerfile
FROM quay.io/keycloak/keycloak:26.1.0

ADD --chown=keycloak:keycloak https://github.com/klausbetz/apple-identity-provider-keycloak/releases/download/1.7.1/apple-identity-provider-1.7.1.jar /opt/keycloak/providers/apple-identity-provider-1.7.1.jar
ADD --chown=keycloak:keycloak https://github.com/mesutpiskin/keycloak-2fa-email-authenticator/releases/download/v0.4/keycloak-2fa-email-authenticator-v0.4-KC21.1.1.jar /opt/keycloak/providers/keycloak-2fa-email-authenticator-v0.4-KC21.1.1.jar
ADD --chown=keycloak:keycloak https://github.com/slemke/keycloak-backup-email/releases/download/v0.0.1/keycloak-backup-email.jar /opt/keycloak/providers/keycloak-backup-email.jar
ADD --chown=keycloak:keycloak https://github.com/leroyguillaume/keycloak-bcrypt/releases/download/v1.6.0/keycloak-bcrypt-1.6.0.jar /opt/keycloak/providers/keycloak-bcrypt-1.6.0.jar
ADD --chown=keycloak:keycloak https://github.com/wadahiro/keycloak-discord/releases/download/v0.5.0/keycloak-discord-0.5.0.jar /opt/keycloak/providers/keycloak-discord-0.5.0.jar
ADD --chown=keycloak:keycloak https://github.com/SnuK87/keycloak-kafka/releases/download/1.2.0/keycloak-kafka-1.2.0-jar-with-dependencies.jar /opt/keycloak/providers/keycloak-kafka-1.2.0-jar-with-dependencies.jar
ADD --chown=keycloak:keycloak https://github.com/aerogear/keycloak-metrics-spi/releases/download/7.0.0/keycloak-metrics-spi-7.0.0.jar /opt/keycloak/providers/keycloak-metrics-spi-7.0.0.jar
ADD --chown=keycloak:keycloak https://github.com/slemke/keycloak-terms-authenticator/releases/download/v0.0.1/keycloak-terms-authenticator-0.0.1.jar /opt/keycloak/providers/keycloak-terms-authenticator-0.0.1.jar

# etc etc

ENTRYPOINT ["/opt/keycloak/bin/kc.sh start --optimized"]
```

After (using Cloakroom):

```dockerfile
FROM quay.io/keycloak/keycloak:26.1.0

ADD --chmod=+x https://github.com/peter-mghendi/cloakroom/releases/download/v0.1.2/cloakroom-v0.1.2_linux_x86_64 /usr/local/bin/cloakroom
ENV CLOAKROOM_WARDROBE=/opt/keycloak/providers
RUN cloakroom restore

ENTRYPOINT ["/opt/keycloak/bin/kc.sh start --optimized"]
```

No more repetitive lines for each plugin—Cloakroom’s manifest contains those details.
A single cloakroom restore command downloads everything, keeping your Dockerfile (and your sanity) intact.

---

## Table of Contents
1. [What is Cloakroom?](#what-is-cloakroom)
2. [Why Use Cloakroom?](#why-use-cloakroom)
3. [Features](#features)
4. [Installation & Setup](#installation--setup)
5. [Configuration Overview](#configuration-overview)
6. [Usage](#usage)
   - [Commands](#commands)
   - [Examples](#examples)
7. [Configuration Examples](#configuration-examples)
8. [Handling Multiple Manifests](#handling-multiple-manifests)
9. [FAQ](#faq)
10. [Contributing](#contributing)
11. [License](#license)

---

## What is Cloakroom?
Cloakroom is a **command-line utility** for managing **Keycloak** plugins in a straightforward way. Instead of manually fetching JARs from GitHub Releases (or other endpoints) within your Dockerfiles or scripts, you define them in Cloakroom’s **manifest**. Cloakroom then downloads these artifacts into your Keycloak providers directory as specified by the `CLOAKROOM_WARDROBE` environment variable.

---

## Why Use Cloakroom?
1. **Minimal & Focused**: Designed exclusively for Keycloak plugin management—no extraneous features to distract you.
2. **Predictable Builds**: Pin plugins to specific tags and artifacts (with optional hashes) for consistent deployments.
3. **Clean Dockerfiles**: Drop all the manual fetches—let Cloakroom handle the plugin retrieval.
4. **Easy Setup**: A single environment variable, `CLOAKROOM_WARDROBE`, points to your Keycloak provider directory.

---

## Features
- **Manifest-Driven**: Centralize plugin definitions in a file like `cloakroom.json` (or TOML/INI/HCL/YAML).
- **GitHub Releases**: Download JARs using `tag` (e.g., `"v1.2.0"`) and an `artifact`.
- **Optional Hash Verification**: Provide a **SHA3-512** `hash` to verify each download’s integrity.
- **Flexible Configuration Formats**: Use JSON, TOML, INI, HCL, or YAML—whichever suits your workflow.
- **Environment-Aware**: Respects `CLOAKROOM_WARDROBE`, so you can easily switch directories across environments.

---

## Installation & Setup

1. **Download the Binary**  
   Grab the latest release from the [Cloakroom GitHub Releases](https://github.com/peter-mghendi/cloakroom/releases).

2. **Install**
   ```bash
   chmod +x cloakroom
   mv cloakroom /usr/local/bin/
   ```

3. **Verify**
   ```bash
   cloakroom --help
   ```

4. **Set `CLOAKROOM_WARDROBE`**
   ```bash
   export CLOAKROOM_WARDROBE="/opt/keycloak/providers"
   ```
   (On Windows, set it in your System Environment Variables.)

---

## Configuration Overview

### Environment Variables
- **`CLOAKROOM_WARDROBE`** (required):  
  The Keycloak provider directory where Cloakroom places the downloaded JARs.
   - Cloakroom refuses to run if not set.

### Manifest File
If the `--manifest` option is specified, cloakroom attempts to load that file, failing if it does not exist or is invalid.
Otherwise, cloakroom searches the current directory for a `cloakroom.{format}` file, where `{format}` is one of:
- [HCL](https://developer.hashicorp.com/terraform/language/syntax/configuration)
- [INI](https://en.wikipedia.org/wiki/INI_file)
- [JSON](https://json.org)
- [TOML](https://toml.io)
- [YAML](https://yaml.org)

Within your manifest, you typically define:
- **`version`**: Mnifest version. The only valid value is `"1.0"`.
- **`host`**: A host that is compatible with the GitHub releases format, e.g. `"github.com"`, or a URL to a hosted GitHub or Gitea server. 
- **`plugins`** (required): a map of `user/repo` → plugin definition.

Each plugin definition contains:
- **`tag`** (required): The title of the release e.g. `"v1.2.0"`.
- **`artifact`** (required): the name of the JAR in that release.
- **`hash`** (optional): A **SHA3-512** hash for integrity checks.

---

## Usage

### Commands

#### `init`
Generates a minimal configuration file:
```
cloakroom init
```

#### `add`
Adds a plugin to your manifest:
```
cloakroom add aerogear/keycloak-metrics-spi --tag 7.0.0 --artifact keycloak-metrics-spi-7.0.0.jar
```
Use `--fetch` to download the plugin right away.

#### `remove`
Removes a plugin from your manifest:
```
cloakroom remove aerogear/keycloak-metrics-spi
```
Use `--purge` to also delete the local JAR file.

#### `restore`
Installs or updates **all** plugins from your manifest:
```
cloakroom restore
```
- `--clean`: Empties the directory defined by `CLOAKROOM_WARDROBE` before downloading.
- `--force`: Overwrites existing JAR files if present.

#### `clean`
Completely clears the directory specified by `CLOAKROOM_WARDROBE`, without modifying your manifest:
```
cloakroom clean
```

#### `list`
Lists all plugins in the manifest, including `tag`, `artifact`, etc.:
```
cloakroom list
```

### Examples

1. **Initialize**
   ```
   cloakroom init
   # Creates cloakroom.json with minimal defaults
   ```

2. **Add a Plugin & Fetch**
   ```
   cloakroom add klausbetz/apple-identity-provider-keycloak \
       --tag 1.7.1 \
       --artifact apple-identity-provider-1.7.1.jar \
       --fetch
   ```

3. **Remove a Plugin & Purge**
   ```
   cloakroom remove klausbetz/apple-identity-provider-keycloak --purge
   ```

4. **Restore Plugins**
   ```
   cloakroom restore --clean
   # Empties CLOAKROOM_WARDROBE, then downloads everything fresh
   ```

---

## Configuration Examples

**HCL**
```hcl
version = "1.0"
host    = "github.com"

plugins "example/my-plugin" {
   tag      = "v1.2.0"
   artifact = "my-plugin-1.2.0.jar"
   hash     = "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"
}
```

**INI**
```ini
version = 1.0
host = github.com

[plugins "example/my-plugin"]
tag = v1.2.0
artifact = my-plugin-1.2.0.jar
hash = a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26
```

**JSON**:
```json
{
  "version": "1.0",
  "host": "github.com",
  "plugins": {
    "example/my-plugin": {
      "tag": "v1.2.0",
      "artifact": "my-plugin-1.2.0.jar",
      "hash": "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"
    }
  }
}
```

**TOML**:
```toml
version = "1.0"
host = "github.com"

[plugins."example/my-plugin"]
tag = "v1.2.0"
artifact = "my-plugin-1.2.0.jar"
hash = "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"
```

**YAML**
```yaml
version: "1.0"
host: "github.com"

plugins:
  "example/my-plugin":
    tag: "v1.2.0"
    artifact: "my-plugin-1.2.0.jar"
    hash: "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"
```
---

## Handling Multiple Manifests

If Cloakroom finds more than one matching manifest (e.g. `cloakroom.json` and `cloakroom.toml`) without a specific `--manifest`:
1. **Fail Immediately**:
   ```
   [ERROR] Found multiple manifests:
   - cloakroom.json
   - cloakroom.toml

   Cloakroom does not support multiple manifests at once.
   ```
2. **Merge** *(planned)*: Cloakroom may eventually merge them, but only if no collisions are detected.

---

## FAQ

**1. Where does Cloakroom store plugins?**  
In the directory defined by `CLOAKROOM_WARDROBE`. If it’s missing, Cloakroom exits with an error.

**2. Can I use private GitHub Repos or GitHub Enterprise?**  
For now, Cloakroom focuses on public GitHub releases. Other hosts (e.g. Gitea) might work if they follow a similar releases structure.

**3. What happens if the JAR already exists?**  
By default, Cloakroom skips it. Use `--force` to overwrite or `--clean` to remove existing files before fetching.

**4. Does Cloakroom handle semver ranges or advanced versioning?**  
Currently, you pin a specific tag. Advanced version logic is on the roadmap.

---

## Contributing
I'd love your contributions—whether that’s opening an issue, suggesting a feature, or sending a pull request. See [CONTRIBUTING.md](#contributing) for details.

---

## License
[MIT License](LICENSE) © 2025 Peter Mghendi
