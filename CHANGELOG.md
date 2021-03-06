# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `pack` rule to Makefile.
- YouTube thumbnail as Album Art fallback.
- AUR name in the README.
- Terminal resize listener (used to improve album art feature).
- Discord Rich Presence integration (`/discord` in the search bar).
- `play` and `simple-play` modes.
- GitHub action that compiles the code.
- A install script for debian based distros.

### Changed
- README screenshot.

### Removed
- `io.ReadAll` calls to keep Tuner compatible with older GoLang versions.
- Unused libmpv files.

## [0.0.2] 2021-02-07

### Added
- SoundCloud search support.
- Shuffle playlist (keybinded to R).
- Simple migration for `data.json`.
- Menu to delete playlist and remove the current song from the playlist.

### Removed
- MPRIS dependency, use libmpv instead.
- Progress bar update from renderPlayer().

## [0.0.1] 2020-12-17

### Added
- Lyric fetcher.
- Album art fetcher.
- Playlist suppport.
- ...
