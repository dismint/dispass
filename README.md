<div align="center">

<img width="60%" alt="dispass logo" src="https://github.com/user-attachments/assets/8429bccc-86c1-4eb0-8bed-689dcd1b0bd4" />

&nbsp;

[![Email](https://img.shields.io/badge/EMAIL-mintjjc%40gmail.com-4b726e?style=flat&labelColor=4d4539)](mailto:mintjjc@gmail.com)
[![Static Badge](https://img.shields.io/badge/WEBSITE-dismint.dev-77743b?style=flat&labelColor=4d4539)](https://www.dismint.dev/)
[![Go](https://github.com/dismint/dispass/actions/workflows/go.yml/badge.svg)](https://github.com/dismint/dispass/actions/workflows/go.yml)

A lightweight and comfortable CLI password manager, written in Go and powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)

</div>

---

<div align="center">
  <img width="80%" alt="dispass demo" src="https://github.com/user-attachments/assets/c2a4e7ba-d485-4bc1-bd3c-3f8d14eff3e7" />
</div>

# ⚡ Features

- 🔐 **Local-first password storage**: All credentials live in a single encrypted file.
- ⚡ **Instant search & autocomplete**: A built-in index for speedy password finding.
- 🔄 **Easy migration**: Import seamlessly from existing password managers.
- 🎨 **Fully customizable**: Tweak colors and appearance to match your terminal setup.

# 🔧 Installation

```bash
go install github.com/dismint/dispass@latest
```

# ⚙️ Configuration

You can configure `dispass` with a `dispass.toml` located either in the working directory or at `$HOME/.config/dispass`

```toml
# dispass.toml default configuration

[colors.light]
symbol          = "#4b726e"
text            = "#4b3d44"
help_key        = "#847875"
help_desc       = "#574852"
help_sep        = "#ab9b8e"
border          = "#4b726e"
message_error   = "#79444a"
message_success = "#4b726e"
message_notif   = "#8caba1"

[colors.dark]
symbol          = "#8caba1"
text            = "#d2c9a5"
help_key        = "#847875"
help_desc       = "#ab9b8e"
help_sep        = "#574852" 
border          = "#8caba1"
message_error   = "#c77b58"
message_success = "#8caba1"
message_notif   = "#4b726e"
```

# 🔨 Development

`dispass` is organized as a standard Go project and can be built as such:
```bash
# build to ./dispass
go build .
# develop
go run .
```
