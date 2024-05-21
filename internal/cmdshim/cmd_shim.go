package cmdshim

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var shebangExpr *regexp.Regexp = regexp.MustCompile(`^#!\s*(?:\/usr\/bin\/env\s+(?:-S\s+)?((?:[^ \t=]+=[^ \t=]+\s+)*))?([^ \t]+)(.*)$`)

func CmdShim(from, to string) error {
	if err := rm(to); err != nil {
		return err
	}

	if err := rm(to + ".cmd"); err != nil {
		return err
	}

	if err := rm(to + ".ps1"); err != nil {
		return err
	}

	writeShim(from, to)

	return nil
}

func rm(to string) error {
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		if err := os.Remove(to); err != nil {
			return fmt.Errorf("unable to delete existing links: %s", err.Error())
		}
	}

	return nil
}

func writeShim(from, to string) error {
	binPath := filepath.Join("node_modules", ".bin")
	if _, binExists := os.Stat(binPath); os.IsNotExist(binExists) {
		if err := os.MkdirAll(binPath, 0755); err != nil {
			return fmt.Errorf("error while creating .bin folder: %s", err.Error())
		}
	}

	file, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("error while reading file %s: %s", from, err.Error())
	}

	scanner := bufio.NewScanner(file)
	var firstLine string

	if scanner.Scan() {
		firstLine = scanner.Text()
	} else {
		return fmt.Errorf("empty executable file: %s", fmt.Sprintf("File %s is empty", to))
	}

	shebang := shebangExpr.FindStringSubmatch(firstLine)

	if shebang == nil {
		if err := write(from, to, "", "", ""); err != nil {
			return err
		}
	} else {
		vars := shebang[1]
		prog := shebang[2]
		args := shebang[3]

		if err := write(from, to, prog, vars, args); err != nil {
			return err
		}
	}

	return nil
}

func write(from, to, prog, args, variables string) error {
	params := NewParams(from, to, prog, args, variables)

	if err := CreateLink(params.To+".cmd", WriteCmd(params)); err != nil {
		return fmt.Errorf("unable to create link for %s: %s", params.To+".cmd", err.Error())
	}

	if err := CreateLink(params.To, WriteSh(params)); err != nil {
		return fmt.Errorf("unable to create link for %s: %s", params.To+".cmd", err.Error())
	}

	if err := CreateLink(params.To+".ps1", WritePwsh(params)); err != nil {
		return fmt.Errorf("unable to create link for %s: %s", params.To+".cmd", err.Error())
	}

	return nil
}

func CreateLink(to string, content string) error {
	file, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("unable to create %s: %s", to, err.Error())
	}
	defer file.Close()
	file.Write([]byte(content))

	return nil
}
