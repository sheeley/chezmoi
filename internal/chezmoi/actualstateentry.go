package chezmoi

import (
	"os"
)

// An ActualStateEntry represents the actual state of an entry in the
// filesystem.
type ActualStateEntry interface {
	EntryState() (*EntryState, error)
	Path() AbsPath
	Remove(system System) error
}

// A ActualStateAbsent represents the absence of an entry in the filesystem.
type ActualStateAbsent struct {
	absPath AbsPath
}

// A ActualStateDir represents the state of a directory in the filesystem.
type ActualStateDir struct {
	absPath AbsPath
	perm    os.FileMode
}

// A ActualStateFile represents the state of a file in the filesystem.
type ActualStateFile struct {
	absPath AbsPath
	perm    os.FileMode
	*lazyContents
}

// A ActualStateSymlink represents the state of a symlink in the filesystem.
type ActualStateSymlink struct {
	absPath AbsPath
	*lazyLinkname
}

// NewActualStateEntry returns a new ActualStateEntry populated with absPath
// from fs.
func NewActualStateEntry(system System, absPath AbsPath, info os.FileInfo, err error) (ActualStateEntry, error) {
	if info == nil {
		info, err = system.Lstat(absPath)
	}
	switch {
	case os.IsNotExist(err):
		return &ActualStateAbsent{
			absPath: absPath,
		}, nil
	case err != nil:
		return nil, err
	}
	//nolint:exhaustive
	switch info.Mode() & os.ModeType {
	case 0:
		return &ActualStateFile{
			absPath: absPath,
			perm:    info.Mode() & os.ModePerm,
			lazyContents: newLazyContentsFunc(func() ([]byte, error) {
				return system.ReadFile(absPath)
			}),
		}, nil
	case os.ModeDir:
		return &ActualStateDir{
			absPath: absPath,
			perm:    info.Mode() & os.ModePerm,
		}, nil
	case os.ModeSymlink:
		return &ActualStateSymlink{
			absPath: absPath,
			lazyLinkname: newLazyLinknameFunc(func() (string, error) {
				return system.Readlink(absPath)
			}),
		}, nil
	default:
		return nil, &errUnsupportedFileType{
			absPath: absPath,
			mode:    info.Mode(),
		}
	}
}

// EntryState returns d's entry state.
func (s *ActualStateAbsent) EntryState() (*EntryState, error) {
	return &EntryState{
		Type: EntryStateTypeRemove,
	}, nil
}

// Path returns d's path.
func (s *ActualStateAbsent) Path() AbsPath {
	return s.absPath
}

// Remove removes d.
func (s *ActualStateAbsent) Remove(system System) error {
	return nil
}

// EntryState returns d's entry state.
func (s *ActualStateDir) EntryState() (*EntryState, error) {
	return &EntryState{
		Type: EntryStateTypeDir,
		Mode: os.ModeDir | s.perm,
	}, nil
}

// Path returns d's path.
func (s *ActualStateDir) Path() AbsPath {
	return s.absPath
}

// Remove removes d.
func (s *ActualStateDir) Remove(system System) error {
	return system.RemoveAll(s.absPath)
}

// EntryState returns d's entry state.
func (s *ActualStateFile) EntryState() (*EntryState, error) {
	contents, err := s.Contents()
	if err != nil {
		return nil, err
	}
	contentsSHA256, err := s.ContentsSHA256()
	if err != nil {
		return nil, err
	}
	return &EntryState{
		Type:           EntryStateTypeFile,
		Mode:           s.perm,
		ContentsSHA256: HexBytes(contentsSHA256),
		contents:       contents,
	}, nil
}

// Path returns d's path.
func (s *ActualStateFile) Path() AbsPath {
	return s.absPath
}

// Remove removes d.
func (s *ActualStateFile) Remove(system System) error {
	return system.RemoveAll(s.absPath)
}

// EntryState returns d's entry state.
func (s *ActualStateSymlink) EntryState() (*EntryState, error) {
	linkname, err := s.Linkname()
	if err != nil {
		return nil, err
	}
	linknameSHA256, err := s.LinknameSHA256()
	if err != nil {
		return nil, err
	}
	return &EntryState{
		Type:           EntryStateTypeSymlink,
		ContentsSHA256: HexBytes(linknameSHA256),
		contents:       []byte(linkname),
	}, nil
}

// Path returns d's path.
func (s *ActualStateSymlink) Path() AbsPath {
	return s.absPath
}

// Remove removes d.
func (s *ActualStateSymlink) Remove(system System) error {
	return system.RemoveAll(s.absPath)
}
