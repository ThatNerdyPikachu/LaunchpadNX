package main

import (
	"bufio"
	"fmt"
	"github.com/cobaugh/osrelease"
	"github.com/shiena/ansicolor"
	"gopkg.in/src-d/go-git.v4"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func resetTerm(w io.Writer) {
	fmt.Fprintf(w, "\x1b[0m")
}

func input(w io.Writer, prompt string) string {
	fmt.Fprintf(w, prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func wait() {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
}

func errCheck(w io.Writer, task string, err error) {
	if err != nil {
		fmt.Fprintf(w, "\x1b[91man error occured while %s:\n", task)
		panic(err)
	}
}

func inArray(array []string, item string) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}

func copy(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		return err
	}

	return nil
}

func copyFolder(src, dst string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, os.ModeDir)

	for _, f := range files {
		copy(src+"/"+f.Name(), dst+"/"+f.Name())
	}

	return nil
}

func main() {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)

	resetTerm(w)
	defer resetTerm(w)

	var osr map[string]string

	if runtime.GOOS == "windows" {
		// check for reqs
		dkpCmds := []string{"pacman", "make", "git"}
		for _, v := range dkpCmds {
			_, err := exec.LookPath(v)
			if err != nil {
				fmt.Fprintf(w, "\x1b[91msorry, but you need \x1b[21m%s\x1b[1m to continue!\n", v)
				fmt.Fprintf(w, "press any key to exit")
				resetTerm(w)
				wait()
				os.Exit(1)
			}
		}
	} else if runtime.GOOS == "linux" {
		// check for reqs
		osr, err = osrelease.Read()
		var dkpCmds []string
		if err == nil && or["NAME"] == "Arch Linux" {
			dkpCmds = []string{"pacman", "make", "git"}
		} else {
			dkpCmds = []string{"dkp-pacman", "make"}
		}
		for _, v := range dkpCmds {
			_, err := exec.LookPath(v)
			if err != nil {
				fmt.Fprintf(w, "\x1b[91msorry, but you need \x1b[21m%s\x1b[1m to continue!\n", v)
				fmt.Fprintf(w, "press any key to exit")
				resetTerm(w)
				wait()
				os.Exit(1)
			}
		}
	} else {
		fmt.Fprintf(w, "\x1b[91msorry, launchpadnx does not yet support your operating system! make sure to open a github issue!\n")
		fmt.Fprintf(w, "press any key to exit")
		resetTerm(w)
		wait()
		os.Exit(1)
	}

	fmt.Fprintf(w, "\x1b[94mwelcome to launchpadnx v2, where we go play with our devices!\n")
	fmt.Fprintf(w, "i'm your host, the magical program that lives inside your computer~\n")
	fmt.Fprintf(w, "(uhh, sorry... back to launching!)\n\n")

	fmt.Fprintf(w, "now here comes the fun part, selecting your features!\n")
	fmt.Fprintf(w, "here are your feature choices (note: as usual, atmosphere's base is selected by default):\n")
	selections := []string{
		"checkpoint (a save manager)",
		"hbmenu",
		"layeredfs (game mods)",
		"sigpatches",
		"sys-ftpd (a background ftp server)",
		"tinfoil (a title manager)",
	}
	for i, v := range selections {
		fmt.Fprintf(w, "\x1b[91m%d: %s\n", i+1, v)
	}

	var features []string

	for {
		resp := input(w, "\x1b[94mplease type the numbers of the features that you want, seperated by spaces (or type '\x1b[91mall\x1b[94m' to compile everything): ")
		features = strings.Split(resp, " ")
		if features[0] == "all" {
			features = []string{"1", "2", "3", "4", "5", "6"}
			break
		} else if features[0] == "exit" {
			resetTerm(w)
			os.Exit(0)
		} else if features[0] == "" {
			continue
		}
		nums := []int{1, 2, 3, 4, 5, 6}
		var good bool
		for _, v := range nums {
			i, err := strconv.Atoi(features[0])
			if err != nil {
				good = false
			} else {
				if i == v {
					good = true
				}
			}
		}
		if good == true {
			break
		}
	}

	var nogc bool
	for {
		resp := input(w, "do you need nogc? (y/n): ")
		if strings.ToLower(resp) == "y" {
			nogc = true
			fmt.Fprintf(w, "\n")
			break
		} else if strings.ToLower(resp) == "n" {
			nogc = false
			fmt.Fprintf(w, "\n")
			break
		}
	}

	_, err := os.Stat("sd_root")
	if err == nil {
		for {
			resp := input(w, "\x1b[91mwarning: you already have a package built. running this build will delete that one.\nare you sure you want to continue? (y/n): ")
			if strings.ToLower(resp) == "y" {
				err = os.RemoveAll("sd_root")
				errCheck(w, "deleting sd_root", err)
				fmt.Fprintf(w, "\x1b[94m\n")
				break
			} else if strings.ToLower(resp) == "n" {
				fmt.Fprintf(w, "build aborted.")
				resetTerm(w)
				os.Exit(0)
			}
		}
	}

	folders := []string{"build/atmosphere", "build/hekate", "build/checkpoint",
		"build/hbmenu", "build/sys-ftpd", "build/tinfoil"}

	ap := false

	for _, f := range folders {
		_, err := os.Stat(f)
		if err == nil {
			_, err = os.Stat(f + "/stratosphere/loader/source/ldr_config.cpp")
			if err == nil {
				cmd := exec.Command("git", "fetch", "upstream")
				cmd.Dir = f
				err = cmd.Run()
				errCheck(w, "updating the sources for "+f, err)

				cmd = exec.Command("git", "merge", "upstream/master")
				cmd.Dir = f
				err = cmd.Run()
				errCheck(w, "merging the sources for "+f, err)
				ap = true
			} else {
				r, err := git.PlainOpen(f)
				errCheck(w, "opening "+f, err)

				wt, err := r.Worktree()
				errCheck(w, "opening the worktree for "+f, err)

				err = wt.Pull(&git.PullOptions{RemoteName: "origin"})
				if err.Error() != "already up-to-date" {
					errCheck(w, "updating the sources for "+f, err)
				}
			}
		}
	}

	fmt.Fprintf(w, "running pacman -Syu...\n")
	err = exec.Command("pacman", "-Syu", "--noconfirm").Run()
	errCheck(w, "running pacman -Syu", err)

	fmt.Fprintf(w, "installing dependencies...\n")
	deps := []string{"switch-dev", "devkitARM"}

	if inArray(features, "1") {
		deps = append(deps, "switch-freetype")
	}

	if inArray(features, "2") {
		deps = append(deps, "switch-freetype", "switch-libconfig")
	}

	if inArray(features, "5") {
		deps = append(deps, "switch-mpg123")
	}

	if inArray(features, "6") {
		deps = append(deps, "switch-curl")
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("pacman", "-S", strings.Join(deps, " "), "--noconfirm", "--needed")
		err = cmd.Run()
		errCheck(w, "installing dependencies", err)
	} else if runtime.GOOS == "linux" {
		if osr["NAME"] == "Arch Linux" {
			cmd := exec.Command("sudo", "pacman", "-S", strings.Join(deps, " "), "--noconfirm", "--needed")
			err = cmd.Run()
			errCheck(w, "installing dependencies", err)
		} else {
			cmd := exec.Command("sudo", "dkp-pacman", "-S", strings.Join(deps, " "), "--noconfirm", "--needed")
			err = cmd.Run()
			errCheck(w, "installing dependencies", err)
		}
	}

	// the goddess that blessed your switch -- <3
	fmt.Fprintf(w, "cloning hekate...\n")
	_, err = git.PlainClone("build/hekate", false, &git.CloneOptions{
		URL: "https://github.com/CTCaer/hekate.git",
	})
	if err != nil && err.Error() != "repository already exists" {
		errCheck(w, "cloning hekate", err)
	}

	fmt.Fprintf(w, "building hekate...\n")
	cmd := exec.Command("make", "-j3")
	cmd.Dir = "build/hekate"
	err = cmd.Run()
	errCheck(w, "building hekate", err)

	fmt.Fprintf(w, "copying files...\n")
	err = copy("build/hekate/output/hekate.bin", "hekate.bin")
	errCheck(w, "copying the hekate payload", err)
	err = os.MkdirAll("sd_root/bootloader/sys", os.ModeDir)
	errCheck(w, "creating sd_root/bootloader/sys", err)
	err = copy("build/hekate/output/libsys_lp0.bso", "sd_root/bootloader/sys/libsys_lp0.bso")
	errCheck(w, "copying the hekate payload", err)

	if ap == false {
		if inArray(features, "2") || inArray(features, "1") || inArray(features, "6") {
			fmt.Fprintf(w, "cloning atmosphere...\n")

			cmd = exec.Command("git", "clone", "https://github.com/rajkosto/Atmosphere.git", "build/atmosphere")
			err = cmd.Run()
			errCheck(w, "cloning atmosphere", err)

			cmd = exec.Command("git", "checkout", "rajnx")
			cmd.Dir = "build/atmosphere"
			err = cmd.Run()
			errCheck(w, "checking out to rajnx", err)

			cmd = exec.Command("git", "remote", "add", "upstream", "https://github.com/Atmosphere-NX/Atmosphere.git")
			cmd.Dir = "build/atmosphere"
			err = cmd.Run()
			errCheck(w, "adding upstream", err)

			cmd = exec.Command("git", "fetch", "upstream")
			cmd.Dir = "build/atmosphere"
			err = cmd.Run()
			errCheck(w, "fetching upstream", err)

			cmd = exec.Command("git", "merge", "upstream/master")
			cmd.Dir = "build/atmosphere"
			err = cmd.Run()
			errCheck(w, "merging upstream into rajnx", err)

			fmt.Fprintf(w, "patching atmosphere...\n")

			cmd = exec.Command("git", "apply", "../../hbl.patch")
			cmd.Dir = "build/atmosphere"
			err = cmd.Run()
			errCheck(w, "patching atmosphere", err)

			fmt.Fprintf(w, "building atmosphere...\n")
			cmd = exec.Command("make", "-j3")
			cmd.Dir = "build/atmosphere/stratosphere"
			err = cmd.Run()
			errCheck(w, "building atmosphere", err)
		} else {
			fmt.Fprintf(w, "cloning atmosphere...\n")

			_, err = git.PlainClone("build/hbmenu", false, &git.CloneOptions{
				URL: "https://github.com/switchbrew/nx-hbmenu.git",
			})
			if err != nil && err.Error() != "repository already exists" {
				errCheck(w, "cloning hbmenu", err)
			}

			fmt.Fprintf(w, "building atmosphere...\n")
			cmd = exec.Command("make", "-j3")
			cmd.Dir = "build/atmosphere/atmosphere"
			err = cmd.Run()
			errCheck(w, "building atmosphere", err)
		}
	}

	fmt.Fprintf(w, "copying files...\n")
	err = os.MkdirAll("sd_root/atmosphere/titles/0100000000000036/exefs", os.ModeDir)
	errCheck(w, "creating sd_root/atmosphere/titles/0100000000000036/exefs", err)
	err = copy("build/atmosphere/stratosphere/creport/creport.npdm", "sd_root/atmosphere/titles/0100000000000036/exefs/main.npdm")
	errCheck(w, "copying creport's npdm", err)
	err = copy("build/atmosphere/stratosphere/creport/creport.nso", "sd_root/atmosphere/titles/0100000000000036/exefs/main")
	errCheck(w, "copying creport's npdm", err)
	err = os.MkdirAll("sd_root/cfw", os.ModeDir)
	errCheck(w, "creating sd_root/cfw", err)
	err = copy("build/atmosphere/stratosphere/loader/loader.kip", "sd_root/cfw/loader.kip")
	errCheck(w, "copying loader", err)
	err = copy("build/atmosphere/stratosphere/pm/pm.kip", "sd_root/cfw/pm.kip")
	errCheck(w, "copying pm", err)
	err = copy("build/atmosphere/stratosphere/sm/sm.kip", "sd_root/cfw/sm.kip")
	errCheck(w, "copying sm", err)

	var (
		hekateConfig []string
		c            []string
	)

	if nogc && inArray(features, "4") {
		hekateConfig = []string{
			"[Stock]",
			"",
			"[Stock (nogc)]",
			"kip1patch=nogc",
			"",
			"[CFW]",
		}
		c = []string{"kip1=cfw/loader.kip", "kip1=cfw/pm.kip", "kip1=cfw/sm.kip", "kip1patch=nogc,nosigchk"}
	} else if nogc {
		hekateConfig = []string{
			"[Stock]",
			"",
			"[Stock (nogc)]",
			"kip1patch=nogc",
			"",
			"[CFW]",
		}
		c = []string{"kip1=cfw/loader.kip", "kip1=cfw/pm.kip", "kip1=cfw/sm.kip", "kip1patch=nogc"}
	} else if inArray(features, "4") {
		hekateConfig = []string{
			"[Stock]",
			"",
			"[CFW]",
		}
		c = []string{"kip1=cfw/loader.kip", "kip1=cfw/pm.kip", "kip1=cfw/sm.kip", "kip1patch=nosigchk"}
	} else {
		hekateConfig = []string{
			"[Stock]",
			"",
			"[CFW]",
		}
		c = []string{"kip1=cfw/loader.kip", "kip1=cfw/pm.kip", "kip1=cfw/sm.kip"}
	}

	if inArray(features, "1") {
		fmt.Fprintf(w, "cloning checkpoint...\n")
		_, err = git.PlainClone("build/checkpoint", false, &git.CloneOptions{
			URL: "https://github.com/FlagBrew/Checkpoint.git",
		})
		if err != nil && err.Error() != "repository already exists" {
			errCheck(w, "cloning checkpoint", err)
		}

		fmt.Fprintf(w, "building checkpoint...\n")
		cmd := exec.Command("make", "-j3")
		cmd.Dir = "build/checkpoint/switch"
		err = cmd.Run()
		errCheck(w, "building checkpoint", err)

		fmt.Fprintf(w, "copying files...\n")
		err = os.MkdirAll("sd_root/switch", os.ModeDir)
		errCheck(w, "creating sd_root/switch", err)
		err = os.MkdirAll("sd_root/switch/Checkpoint", os.ModeDir)
		errCheck(w, "creating sd_root/switch/Checkpoint", err)
		err = copy("build/checkpoint/switch/out/Checkpoint.nro", "sd_root/switch/Checkpoint/Checkpoint.nro")
		errCheck(w, "copying checkpoint", err)
	}

	if inArray(features, "2") || inArray(features, "1") || inArray(features, "6") {
		fmt.Fprintf(w, "cloning hbloader...\n")
		_, err = git.PlainClone("build/hbloader", false, &git.CloneOptions{
			URL: "https://github.com/switchbrew/nx-hbloader.git",
		})
		if err != nil && err.Error() != "repository already exists" {
			errCheck(w, "cloning hbloader", err)
		}

		fmt.Fprintf(w, "building hbloader...\n")
		cmd := exec.Command("make", "-j3")
		cmd.Dir = "build/hbloader"
		err = cmd.Run()
		errCheck(w, "building hbloader", err)

		fmt.Fprintf(w, "copying files...\n")
		err = copy("build/hbloader/hbl.nsp", "sd_root/cfw/hbl.nsp")
		errCheck(w, "copying hbloader", err)

		fmt.Fprintf(w, "creating hbloader config...\n")

		f, err := os.Create("sd_root/cfw/hbl.ini")
		errCheck(w, "creating sd_root/cfw/hbl.ini", err)

		_, err = f.WriteString("[config]\n")
		errCheck(w, "writing to hbl.ini", err)
		_, err = f.WriteString("hbTitleId=010000000000100D\n")
		errCheck(w, "writing to hbl.ini", err)
		_, err = f.WriteString("hbKeyCombo=R")
		errCheck(w, "writing to hbl.ini", err)

		f.Close()

		fmt.Fprintf(w, "installing hbmenu's dependencies...\n")
		err = exec.Command("pacman", "-S", "switch-dev", "devkitARM", "switch-freetype", "switch-libconfig", "--noconfirm", "--needed").Run()
		errCheck(w, "installing hbmenu's dependencies", err)

		fmt.Fprintf(w, "cloning hbmenu...\n")
		_, err = git.PlainClone("build/hbmenu", false, &git.CloneOptions{
			URL: "https://github.com/switchbrew/nx-hbmenu.git",
		})
		if err != nil && err.Error() != "repository already exists" {
			errCheck(w, "cloning hbmenu", err)
		}

		fmt.Fprintf(w, "building hbmenu...\n")
		cmd = exec.Command("make", "nx", "-j3")
		cmd.Dir = "build/hbmenu"
		err = cmd.Run()
		errCheck(w, "building hbmenu", err)

		fmt.Fprintf(w, "copying files...\n")
		err = copy("build/hbmenu/hbmenu.nro", "sd_root/hbmenu.nro")
		errCheck(w, "copying hbmenu", err)
	}

	if inArray(features, "3") {
		fmt.Fprintf(w, "copying files...\n")
		err = copy("build/atmosphere/stratosphere/fs_mitm/fs_mitm.kip", "sd_root/cfw/fs_mitm.kip")
		errCheck(w, "copying fs_mitm (layeredfs)", err)
		c = append(c, "atmosphere=1", "kip1=cfw/fs_mitm.kip")
	}

	if inArray(features, "4") {
		err = os.MkdirAll("sd_root/atmosphere/exefs_patches", os.ModeDir)
		errCheck(w, "creating sd_root/atmosphere/exefs_patches", err)
		fmt.Fprintf(w, "copying files...\n")
		err = copyFolder("fake_tickets", "sd_root/atmosphere/exefs_patches/fake_tickets")
		errCheck(w, "copying fake_tickets (sigpatches)", err)
	}

	if inArray(features, "5") {
		fmt.Fprintf(w, "cloning sys-ftpd...\n")
		_, err = git.PlainClone("build/sys-ftpd", false, &git.CloneOptions{
			URL: "https://github.com/jakibaki/sys-ftpd.git",
		})
		if err != nil && err.Error() != "repository already exists" {
			errCheck(w, "cloning sys-ftpd", err)
		}

		fmt.Fprintf(w, "building sys-ftpd...\n")
		cmd := exec.Command("make", "-j3")
		cmd.Dir = "build/sys-ftpd"
		err = cmd.Run()
		errCheck(w, "building sys-ftpd", err)

		fmt.Fprintf(w, "copying files...\n")
		err = copyFolder("build/sys-ftpd/sd_card/ftpd", "sd_root/ftpd")
		errCheck(w, "copying sys-ftpd's sound files", err)
		err = copy("build/sys-ftpd/sys-ftpd.kip", "sd_root/cfw/sys-ftpd.kip")
		errCheck(w, "copying sys-ftpd", err)
		c = append(c, "kip1=cfw/sys-ftpd.kip")
	}

	if inArray(features, "6") {
		fmt.Fprintf(w, "cloning tinfoil...\n")
		_, err = git.PlainClone("build/tinfoil", false, &git.CloneOptions{
			URL: "https://github.com/Adubbz/Tinfoil.git",
		})
		if err != nil && err.Error() != "repository already exists" {
			errCheck(w, "cloning tinfoil", err)
		}

		fmt.Fprintf(w, "building tinfoil...\n")
		cmd := exec.Command("make", "-j3")
		cmd.Dir = "build/tinfoil"
		err = cmd.Run()
		errCheck(w, "building tinfoil", err)

		fmt.Fprintf(w, "copying files...\n")
		err = os.MkdirAll("sd_root/switch", os.ModeDir)
		errCheck(w, "creating sd_root/switch", err)
		err = copy("build/tinfoil/Tinfoil.nro", "sd_root/switch/Tinfoil.nro")
		errCheck(w, "copying tinfoil", err)
	}

	fmt.Fprintf(w, "creating hekate config...\n\n")
	sort.Strings(c)
	for _, v := range c {
		hekateConfig = append(hekateConfig, v)
	}

	f, err := os.Create("sd_root/bootloader/hekate_ipl.ini")
	errCheck(w, "creating sd_root/bootloader/hekate_ipl.ini", err)
	for _, v := range hekateConfig {
		_, err := f.WriteString(v + "\n")
		errCheck(w, "writing to hekate_ipl.ini", err)
	}
	f.Close()

	fmt.Fprintf(w, "done!")
}
