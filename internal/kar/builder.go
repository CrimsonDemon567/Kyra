package kar

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

// ---------------------------
// KAR format
// ---------------------------
//
// Header:
//   K A R 1
//
// Layout:
//   uint32 fileCount
//
//   For each file:
//       uint32 nameLen
//       bytes  name
//       uint32 dataLen
//       bytes  data
//
// ---------------------------

type FileEntry struct {
	Name string
	Data []byte
}

type Archive struct {
	Files []FileEntry
}

func New() *Archive {
	return &Archive{
		Files: []FileEntry{},
	}
}

func (a *Archive) AddFile(name string, data []byte) {
	a.Files = append(a.Files, FileEntry{
		Name: name,
		Data: data,
	})
}

func (a *Archive) AddFileFromDisk(path string, nameInArchive string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	a.AddFile(nameInArchive, data)
	return nil
}

// ---------------------------
// Encode KAR archive
// ---------------------------

func (a *Archive) Encode() []byte {
	buf := &bytes.Buffer{}

	// Header
	buf.Write([]byte{'K', 'A', 'R', 1})

	// File count
	binary.Write(buf, binary.LittleEndian, uint32(len(a.Files)))

	for _, f := range a.Files {
		// Name
		binary.Write(buf, binary.LittleEndian, uint32(len(f.Name)))
		buf.Write([]byte(f.Name))

		// Data
		binary.Write(buf, binary.LittleEndian, uint32(len(f.Data)))
		buf.Write(f.Data)
	}

	return buf.Bytes()
}

// ---------------------------
// Save to disk
// ---------------------------

func (a *Archive) Save(path string) error {
	data := a.Encode()
	return os.WriteFile(path, data, 0644)
}

// ---------------------------
// Build a KAR from a folder
// ---------------------------

func BuildFromFolder(folder string) (*Archive, error) {
	arc := New()

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(folder, path)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		arc.AddFile(rel, data)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return arc, nil
}

// ---------------------------
// Debug print
// ---------------------------

func (a *Archive) String() string {
	out := "KAR Archive:\n"
	for _, f := range a.Files {
		out += fmt.Sprintf(" - %s (%d bytes)\n", f.Name, len(f.Data))
	}
	return out
}
