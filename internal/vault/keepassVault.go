package vault

import (
	"errors"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
)

// https://github.com/tobischo/gokeepasslib

type Vault interface {
	GetEntry(key string) (string, error)
	SetEntry(key string, value string) error
	SetDescription(key string, description string) error
	Close() error
}

type KeepassVault struct {
	filePath string
	db       *gokeepasslib.Database
	group    *gokeepasslib.Group
}

func NewKeepassVault(filePath string, password string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	db := gokeepasslib.NewDatabase(
		gokeepasslib.WithDatabaseKDBXVersion4(),
	)
	db.Content.Meta.DatabaseName = "confy vault"
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)

	// Create root group and "confy" subgroup
	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "confy"
	db.Content.Root.Groups = []gokeepasslib.Group{rootGroup}
	rootGroup.Entries = []gokeepasslib.Entry{}

	db.LockProtectedEntries()

	encoder := gokeepasslib.NewEncoder(file)
	if err := encoder.Encode(db); err != nil {
		return err
	}

	return nil
}

func OpenKeepassVault(filePath string, password string) (Vault, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)

	if err := gokeepasslib.NewDecoder(file).Decode(db); err != nil {
		return nil, err
	}

	// Unlock protected entries to allow modifications
	if err := db.UnlockProtectedEntries(); err != nil {
		return nil, err
	}

	var mainGroup *gokeepasslib.Group
	for i := range db.Content.Root.Groups {
		if db.Content.Root.Groups[i].Name == "confy" {
			mainGroup = &db.Content.Root.Groups[i]
			break
		}
	}

	if mainGroup == nil {
		newGroup := gokeepasslib.NewGroup()
		newGroup.Name = "confy"
		db.Content.Root.Groups = append(db.Content.Root.Groups, newGroup)
		mainGroup = &db.Content.Root.Groups[len(db.Content.Root.Groups)-1]
	}

	return &KeepassVault{
		filePath: filePath,
		db:       db,
		group:    mainGroup,
	}, nil
}

func (v *KeepassVault) GetEntry(key string) (string, error) {
	for _, entry := range v.group.Entries {
		var title string
		var value string

		for _, val := range entry.Values {
			switch val.Key {
			case "Title":
				title = val.Value.Content
			case "Password":
				value = val.Value.Content
			}
		}

		if title == key {
			return value, nil
		}
	}

	return "", errors.New("entry not found")
}

func (v *KeepassVault) SetEntry(key string, value string) error {
	// If entry already exists: update
	for i := range v.group.Entries {
		var titleIndex = -1
		var passwordIndex = -1
		var title string

		for j, val := range v.group.Entries[i].Values {
			switch val.Key {
			case "Title":
				title = val.Value.Content
				titleIndex = j
			case "Password":
				passwordIndex = j
			}
		}

		if title == key {
			if passwordIndex >= 0 {
				v.group.Entries[i].Values[passwordIndex].Value.Content = value
			} else {
				v.group.Entries[i].Values = append(v.group.Entries[i].Values, mkValue("Password", value))
			}

			if titleIndex == -1 {
				v.group.Entries[i].Values = append(v.group.Entries[i].Values, mkValue("Title", key))
			}

			return nil
		}
	}

	// Otherwise, create a new entry
	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values, mkValue("Title", key))
	entry.Values = append(entry.Values, mkValue("Password", value))

	v.group.Entries = append(v.group.Entries, entry)
	return nil
}

func (v *KeepassVault) SetDescription(key string, description string) error {
	for i := range v.group.Entries {
		var title string
		var notesIndex = -1

		for j, val := range v.group.Entries[i].Values {
			switch val.Key {
			case "Title":
				title = val.Value.Content
			case "Notes":
				notesIndex = j
			}
		}

		if title == key {
			if notesIndex >= 0 {
				v.group.Entries[i].Values[notesIndex].Value.Content = description
			} else {
				v.group.Entries[i].Values = append(v.group.Entries[i].Values, mkValue("Notes", description))
			}
			return nil
		}
	}

	return errors.New("entry not found")
}

func (v *KeepassVault) Close() error {
	file, err := os.Create(v.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	v.db.LockProtectedEntries()

	encoder := gokeepasslib.NewEncoder(file)
	if err := encoder.Encode(v.db); err != nil {
		return err
	}

	return nil
}

func mkValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key:   key,
		Value: gokeepasslib.V{Content: value},
	}
}
