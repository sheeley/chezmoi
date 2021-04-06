# chezmoi Architecture

<!--- toc --->
* [Directory structure](#directory-structure)
* [Key concepts](#key-concepts)
* [`cmd.Config`](#cmdconfig)
* [Systems](#systems)
  * [Real systems](#real-systems)
  * [Virtual systems](#virtual-systems)
* [Path handling](#path-handling)
* [Persistent state](#persistent-state)
* [Detecting changes](#detecting-changes)
* [Encryption](#encryption)
* [Run once scripts](#run-once-scripts)
* [Testing](#testing)
  * [Unit tests](#unit-tests)
  * [Command tests](#command-tests)
  * [`testscript` tests](#testscript-tests)

## Directory structure

The important directories in chezmoi are:

| Directory | Contents |
| --------- | -------- |
| `cmd/` | Code for the `chezmoi` command. `cmd/*cmd.go` contains the code for each individual command and `cmd/*templatefuncs.go` contain the template functions. |
| `docs/` | The documentation single source of truth. Help text, examples, and the [chezmoi.io](https://chezmoi.io) website are generated from the files in this directory, particularly `docs/REFERENCE.md`. |
| `internal/chezmoi/` | chezmoi's core functionality. |
| `testdata/scripts/` | High-level tests of chezmoi's commands using [`testscript`](https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript). |

## Key concepts

As described in the [reference manual](REFERENCE.md), chezmoi evalutes the
source state to compute a target state for the destination directory (typically
your home directory). It then compares the target state to the actual state of
the destination directory and performs any changes necessary to update the
destination directory to match the target state. The concepts are represented
directly in chezmoi's code.

`chezmoi/internal.SourceState` contains the code for reading a source state,

chezmoi uses the generic term *entry* to describe something that it manages.
Entries can be files, directories, symlinks, scripts, amongst other things.

## `cmd.Config`

## Systems

### Real systems

### Virtual systems

## Path handling

chezmoi uses separate types for absolute paths (`chezmoi/internal.AbsPath`) and
relative paths (`chezmoi/internal.RelPath`) to avoid errors where paths are
combined (e.g. joining two absolute paths). A futher type
`chezmoi/internal.SourceRelPath` is a relative path within the source directory
and handles file and directory attributes.

Internally, chezmoi normalizes all paths to use forward slashes with an optional
upper-cased Windows volume. Paths read from the user may include tilde (`~`) to
represent the user's home directory, use forward or backward slashes, and are
treated as external paths (`chezmoi/internal.ExtPath`). These are normalized to
absolute paths. chezmoi is case-sensitive internally and makes no attempt to
handle case-insensitive or case-preserving filesystems.

## Persistent state

Persistent state is treated as a two-level key-value store with the
pseudo-structure `map[Bucket]map[Key]Value`, where `Bucket`, `Key`, and `Value`
are all `[]byte`s. The `chezmoi/internal.PersistentState` interface defines
interaction with them. Sometimes temporary persistent states are used. For
example, in dry run mode (`--dry-run`) the actual persistent state is copied
into a temporary persistent state in memory which remembers writes but does not
persist them to disk.

## Detecting changes

## Encryption

Encryption tools (currently gpg and age) are abstracted by the
`chezmoi/internal.Encryption` interface.

## Run once scripts

## Testing

### Unit tests

### Command tests

### `testscript` tests
