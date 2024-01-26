package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gofrs/uuid"
	"golang.design/x/clipboard"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
	Commandline       string `json:"commandline"`
	Guid              string `json:"guid"`
	Hidden            bool   `json:"hidden"`
	Name              string `json:"name"`
	Icon              string `json:"icon"`
	StartingDirectory string `json:"startingDirectory,omitempty"`
}

func main() {
	var (
		clp = flag.Bool("c", false, "output to clipboard")
		shl = flag.String("s", "bash", "shell (bash/zsh)")
		ist = flag.String("i", "C:/msys64", "MSYS2 install path")
		typ = flag.String("t", "msys2,ucrt64,mingw32,mingw64,clang64,clang32,clangarm64", "types of msys2")
		sdy = flag.String("d", "C:\\msys64\\home\\%USERNAME%", "starting directory")
		gfw = flag.Bool("gfw", false, "adding a entry of Git for Windows")
	)

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage of %s:
    %s

Options:
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	instgt := "C:/msys64"
	if *ist != "" {
		instgt = *ist
	}

	shopt := "bash"
	if *shl != "" {
		shopt = *shl
	}

	sdopt := "C:/msys64/home/%USERNAME%"
	if *sdy != "" {
		sdopt = *sdy
	}

	profiles, err := msys2profile(*typ, instgt, shopt, sdopt)
	if err != nil {
		slog.Error("fail to msys2profile()", slog.Any("error", err))
		os.Exit(1)
	}

	if *gfw {
		gfwp, err := gitforwindows()
		if err != nil {
			slog.Error("fail to gitforwindows()", slog.Any("error", err))
			os.Exit(1)
		}
		profiles = append(profiles, gfwp)
	}

	v, err := json.MarshalIndent(profiles, "        ", "    ")
	if err != nil {
		slog.Error("fail to json.MarshalIndent()", slog.Any("error", err))
	}

	if *clp {
		outtocb(v)
	} else {
		printProfiles(v)
	}
}

func msys2profile(typ, instgt, shopt, sdopt string) ([]Profile, error) {
	cmdpath := filepath.Join(instgt, "msys2_shell.cmd")
	commonopt := " -defterm -here -no-start "

	profiles := []Profile{}

	rtyp := strings.ReplaceAll(typ, " ", "")
	ts := strings.Split(rtyp, ",")

	for _, t := range ts {
		v4, err := uuid.NewV4()
		if err != nil {
			return profiles, err
		}

		profiles = append(profiles,
			Profile{
				Commandline:       cmdpath + commonopt + "-" + t + " -shell " + shopt,
				Guid:              "{" + v4.String() + "}",
				Name:              "MSYS2: " + t,
				Icon:              filepath.Join(instgt, t+".ico"),
				StartingDirectory: sdopt,
			},
		)
	}
	return profiles, nil
}

func gitforwindows() (Profile, error) {
	v4, err := uuid.NewV4()
	if err != nil {
		return Profile{}, err
	}

	gfwpath := "C:/Program Files/Git"

	return Profile{
		Commandline:       filepath.Join(gfwpath, "bin/bash.exe") + " -i -l",
		Guid:              "{" + v4.String() + "}",
		Name:              "Git Bash",
		Icon:              filepath.Join(gfwpath, "mingw64/share/git/git-for-windows.ico"),
		StartingDirectory: "%USERPROFILE%",
	}, nil
}

func printProfiles(p []byte) {
	scanner := bufio.NewScanner(strings.NewReader(string(p)))
	n := 0
	for scanner.Scan() {
		if n == 0 {
			n++
			continue
		}
		if scanner.Text() == "        ]" {
			break
		}
		if n == 1 {
			fmt.Println(",")
			n++
		}
		fmt.Println(scanner.Text())
	}
}

func outtocb(p []byte) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	b := []byte{}
	scanner := bufio.NewScanner(strings.NewReader(string(p)))
	n := 0
	for scanner.Scan() {
		if n == 0 {
			n++
			continue
		}
		if scanner.Text() == "        ]" {
			break
		}
		if n == 1 {
			b = append(b, ',', '\n')
			n++
		}
		b = append(b, scanner.Bytes()...)
		b = append(b, '\n')
	}

	clipboard.Write(clipboard.FmtText, b)
	return nil
}

// https://stackoverflow.com/questions/71045716/adding-msys-to-windows-terminal
