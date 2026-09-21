# [![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![GPL License][license-shield]][license-url]

<div align="center"> 
<a href="https://mono.tr/">
  <img src="https://r2.mono.tr/logo/Mono-Logo.svg" width="340"/>
</a>

<h2 align="center">cfcname</h2>
<b>cfcname</b> is a small Go tool that finds and bulk-rewrites <code>CNAME</code> records across every zone in a Cloudflare account. Point it at a target hostname and it lists every matching record; give it a replacement and it updates them all in place. Useful when a backend, CDN endpoint or load balancer is renamed and dozens of zones still point at the old name.
</div>

---

## Table of Contents

- [Table of Contents](#table-of-contents)
- [How it works](#how-it-works)
- [Usage](#usage)
- [Building](#building)
- [License](#license)

---

## How it works

1. Authenticates against the Cloudflare API with the given API token.
2. Lists every zone the token has access to.
3. Fetches the `CNAME` records of all zones concurrently, one goroutine per zone.
4. Selects the records whose target matches `-target-url` — either as an exact match (`-full-url`) or as a suffix of the record's target (the default).
5. Prints the matching records grouped by zone. If no `-new-url` is given, it stops here, so the tool is a read-only search by default.
6. If `-new-url` is given, it rewrites each matching record and prints its before/after state. In suffix mode only the matched suffix is replaced, so the record's own prefix is preserved: `app.old.example.com` becomes `app.new.example.com`.

---

## Usage

Build or download the `cfcname` binary, then run it with the required flags.

List every record across all zones that points at something ending in `old.example.com`:

```bash
./cfcname -token 'CF_API_TOKEN' -target-url old.example.com
```

Rewrite them to `new.example.com`, keeping each record's own prefix:

```bash
./cfcname -token 'CF_API_TOKEN' -target-url old.example.com -new-url new.example.com
```

Match and replace one exact target instead of a suffix:

```bash
./cfcname -token 'CF_API_TOKEN' -full-url -target-url lb1.old.example.com -new-url lb2.new.example.com
```

### Flags

| Flag         | Short | Default      | Description                                                        |
|--------------|-------|--------------|--------------------------------------------------------------------|
| `-token`     | `-t`  | *(required)* | Cloudflare API token                                                |
| `-target-url`| `-u`  | *(required)* | Target hostname to look for in `CNAME` records                      |
| `-new-url`   | `-n`  | *(empty)*    | New target to write. Leave empty to only list the matching records  |
| `-full-url`  | `-f`  | `false`      | Match `-target-url` as the record's full target instead of a suffix |

> **Note:** the API token needs `Zone:Read` and `DNS:Edit` permissions on every zone you want it to touch. Run the tool without `-new-url` first and check the listed records — with `-new-url` it updates everything it matched, across every zone, with no further confirmation.

---

## Building

To build the binary:

```bash
go build -o bin/cfcname .
```

The resulting binary will be in the `bin` folder.

---

## License

cfcname is licensed under GPL-3.0-only. See [LICENSE](LICENSE) file for details.

[contributors-shield]: https://img.shields.io/github/contributors/monobilisim/cfcname.svg?style=for-the-badge
[contributors-url]: https://github.com/monobilisim/cfcname/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/monobilisim/cfcname.svg?style=for-the-badge
[forks-url]: https://github.com/monobilisim/cfcname/network/members
[stars-shield]: https://img.shields.io/github/stars/monobilisim/cfcname.svg?style=for-the-badge
[stars-url]: https://github.com/monobilisim/cfcname/stargazers
[issues-shield]: https://img.shields.io/github/issues/monobilisim/cfcname.svg?style=for-the-badge
[issues-url]: https://github.com/monobilisim/cfcname/issues
[license-shield]: https://img.shields.io/github/license/monobilisim/cfcname.svg?style=for-the-badge
[license-url]: https://github.com/monobilisim/cfcname/blob/main/LICENSE
