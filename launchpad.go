/*
   launchpad.go -- launchpadnx's one and only file
   written by pika, licensed under the gnu gpl
   you can grab a copy at https://www.gnu.org/licenses/gpl-3.0.en.html
*/

package main

import (
	"bufio"
	"fmt"
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

func resetTerm(w *io.Writer) {
	fmt.Fprintf(*w, "\x1b[0m")
}

func input(w *io.Writer, prompt string) string {
	fmt.Fprintf(*w, prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func wait() {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
}

func errCheck(w *io.Writer, task string, err error) {
	if err != nil {
		fmt.Fprintf(*w, "\x1b[91man error occured while %s:\n", task)
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

func copyFile(src, dst string) error {
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

	err = os.MkdirAll(dst, 0700)

	for _, f := range files {
		copyFile(src+"/"+f.Name(), dst+"/"+f.Name())
	}

	return nil
}

func main() {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)

	resetTerm(&w)
	defer resetTerm(&w)

	if runtime.GOOS == "windows" {
		// check for reqs
		dkpCmds := []string{"pacman", "make"}
		for _, v := range dkpCmds {
			_, err := exec.LookPath(v)
			if err != nil {
				fmt.Fprintf(w, "\x1b[91msorry, but you need \x1b[21m%s\x1b[1m to continue!\n", v)
				fmt.Fprintf(w, "press any key to exit")
				resetTerm(&w)
				wait()
				os.Exit(1)
			}
		}
	} else if runtime.GOOS == "linux" {
		// check for reqs
		dkpCmds := []string{}
		_, err := exec.LookPath("pacman")
		if err == nil {
			dkpCmds = []string{"pacman", "make"}
		} else {
			dkpCmds = []string{"dkp-pacman", "make"}
		}
		for _, v := range dkpCmds {
			_, err := exec.LookPath(v)
			if err != nil {
				fmt.Fprintf(w, "\x1b[91msorry, but you need \x1b[21m%s\x1b[1m to continue!\n", v)
				fmt.Fprintf(w, "press any key to exit")
				resetTerm(&w)
				wait()
				os.Exit(1)
			}
		}
	} else {
		fmt.Fprintf(w, "\x1b[91msorry, launchpadnx does not yet support your operating system! make sure to open a github issue!\n")
		fmt.Fprintf(w, "press any key to exit")
		resetTerm(&w)
		wait()
		os.Exit(1)
	}

	fmt.Fprintf(w, "\x1b[94mwelcome to launchpadnx v2, where we go play with our devices!\n")
	fmt.Fprintf(w, "i'm your host, the magical program that lives inside your computer~\n")
	fmt.Fprintf(w, "(uhh, sorry... back to launching!)\n\n")

	fmt.Fprintf(w, "now here comes the fun part, selecting your features!\n")
	fmt.Fprintf(w, "here are your feature choices (note: as usual, atmosphere's base is selected by default):\n")
	selections := []string{
		"hbmenu",
		"layeredfs (game mods)",
		"sigpatches",
		"sys-ftpd (a background ftp server)",
		"tinfoil (a title manager)",
	}
	for i, v := range selections {
		fmt.Fprintf(w, "\x1b[91m%d: %s\n", i+1, v)
	}

	features := []string{}

	for {
		resp := input(&w, "\x1b[94mplease type the numbers of the features that you want, seperated by spaces (or type '\x1b[91mall\x1b[94m' to compile everything): ")
		features = strings.Split(resp, " ")
		if features[0] == "all" && !inArray(features, "no-hbl") {
			features = []string{"1", "2", "3", "4", "5"}
			break
		} else if features[0] == "all" && inArray(features, "no-hbl") {
			features = []string{"1", "2", "3", "4", "5", "no-hbl"}
			break
		} else if features[0] == "exit" {
			resetTerm(&w)
			os.Exit(0)
		} else if features[0] == "" {
			continue
		}
		nums := []int{1, 2, 3, 4, 5}
		good := false
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

	nogc := false
	for {
		resp := input(&w, "do you need nogc? (y/n): ")
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
			resp := input(&w, "\x1b[91mwarning: you already have a package built. running this build will delete that one.\nare you sure you want to continue? (y/n): ")
			if strings.ToLower(resp) == "y" {
				errCheck(&w, "deleting sd_root", os.RemoveAll("sd_root"))
				fmt.Fprintf(w, "\x1b[94m\n")
				break
			} else if strings.ToLower(resp) == "n" {
				fmt.Fprintf(w, "build aborted.\n")
				resetTerm(&w)
				os.Exit(0)
			}
		}
	}

	folders := []string{"build/atmosphere", "build/atmosphere/common/include/boost", "build/hekate",
		"build/hbmenu", "build/sys-ftpd", "build/tinfoil"}

	for _, f := range folders {
		_, err := os.Stat(f)
		if err == nil {
			r, err := git.PlainOpen(f)
			errCheck(&w, "opening "+f, err)

			wt, err := r.Worktree()
			errCheck(&w, "opening the worktree for "+f, err)

			err = wt.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil && err.Error() != "already up-to-date" {
				errCheck(&w, "updating the sources for "+f, err)
			}

			submodules, err := wt.Submodules()
			errCheck(&w, "getting the submodules for "+f, err)

			err = submodules.Update(&git.SubmoduleUpdateOptions{})
			errCheck(&w, "updating the submodules for "+f, err)
		}
	}

	fmt.Fprintf(w, "running pacman -Syu...\n")
	if runtime.GOOS == "windows" {
		err = exec.Command("pacman", "-Syu", "--noconfirm").Run()
	} else if runtime.GOOS == "linux" {
		_, err = exec.LookPath("pacman")
		if err == nil {
			err = exec.Command("sudo", "pacman", "-Syu", "--noconfirm").Run()
		} else {
			err = exec.Command("sudo", "dkp-pacman", "-Syu", "--noconfirm").Run()
		}
	}
	errCheck(&w, "running pacman -Syu", err)

	fmt.Fprintf(w, "installing dependencies...\n")

	args := []string{}

	if runtime.GOOS == "windows" {
		args = []string{"-S", "--noconfirm", "--needed", "switch-dev", "devkitARM"}
	} else if runtime.GOOS == "linux" {
		_, err = exec.LookPath("pacman")
		if err == nil {
			args = []string{"pacman", "-S", "--noconfirm", "--needed", "switch-dev", "devkitARM"}
		} else {
			args = []string{"dkp-pacman", "-S", "--noconfirm", "--needed", "switch-dev", "devkitARM"}
		}
	}

	// LIBS	:= `freetype-config --libs` -lconfig -lturbojpeg
	if inArray(features, "1") {
		args = append(args, "switch-freetype", "switch-libconfig", "switch-libjpeg-turbo")
	}

	// LIBS	:= -lnx -lmpg123 -lm
	if inArray(features, "4") {
		args = append(args, "switch-mpg123")
	}

	// LIBS	:= `freetype-config --libs` -lcurl -lz -lnx
	if inArray(features, "5") {
		args = append(args, "switch-freetype", "switch-curl", "switch-zlib")
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("pacman", args...)
		errCheck(&w, "installing dependencies", cmd.Run())
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("sudo", args...)
		errCheck(&w, "installing dependencies", cmd.Run())
	}

	// the goddess that blessed your switch -- <3
	fmt.Fprintf(w, "cloning hekate...\n")
	_, err = git.PlainClone("build/hekate", false, &git.CloneOptions{
		URL: "https://github.com/CTCaer/hekate.git",
	})
	if err != nil && err.Error() != "repository already exists" {
		errCheck(&w, "cloning hekate", err)
	}

	fmt.Fprintf(w, "building hekate...\n")
	cmd := exec.Command("make", "-j")
	cmd.Dir = "build/hekate"
	errCheck(&w, "building hekate", cmd.Run())

	fmt.Fprintf(w, "copying files...\n")
	errCheck(&w, "copying the hekate payload", copyFile("build/hekate/output/hekate.bin", "hekate.bin"))
	errCheck(&w, "creating sd_root/bootloader/sys", os.MkdirAll("sd_root/bootloader/sys", 0700))
	errCheck(&w, "copying the hekate payload", copyFile("build/hekate/output/libsys_lp0.bso",
		"sd_root/bootloader/sys/libsys_lp0.bso"))

	fmt.Fprintf(w, "cloning atmosphere...\n")

	_, err = git.PlainClone("build/atmosphere", false, &git.CloneOptions{
		URL: "https://github.com/Atmosphere-NX/Atmosphere.git",
	})
	if err != nil && err.Error() != "repository already exists" {
		errCheck(&w, "cloning atmosphere", err)
	}

	_, err = git.PlainClone("build/atmosphere/common/include/boost", false, &git.CloneOptions{
		URL: "https://github.com/Atmosphere-NX/ext-boost.git",
	})
	if err != nil && err.Error() != "repository already exists" {
		errCheck(&w, "cloning ext-boost", err)
	}

	fmt.Fprintf(w, "building exosphere...\n")
	cmd = exec.Command("make", "-j")
	cmd.Dir = "build/atmosphere/exosphere"
	errCheck(&w, "building exosphere", cmd.Run())

	fmt.Fprintf(w, "building stratosphere...\n")
	cmd = exec.Command("make", "-j")
	cmd.Dir = "build/atmosphere/stratosphere"
	errCheck(&w, "building stratosphere", cmd.Run())

	fmt.Fprintf(w, "copying files...\n")
	errCheck(&w, "creating sd_root/atmosphere/titles/0100000000000036",
		os.MkdirAll("sd_root/atmosphere/titles/0100000000000036", 0700))

	errCheck(&w, "copying creport's exefs", copyFile("build/atmosphere/stratosphere/creport/creport.nsp",
		"sd_root/atmosphere/titles/0100000000000036/exefs.nsp"))
	errCheck(&w, "creating sd_root/cfw", os.MkdirAll("sd_root/cfw", 0700))

	errCheck(&w, "copying exosphere", copyFile("build/atmosphere/exosphere/exosphere.bin",
		"sd_root/cfw/exosphere.bin"))

	errCheck(&w, "copying loader", copyFile("build/atmosphere/stratosphere/loader/loader.kip",
		"sd_root/cfw/loader.kip"))

	errCheck(&w, "copying pm", copyFile("build/atmosphere/stratosphere/pm/pm.kip", "sd_root/cfw/pm.kip"))

	errCheck(&w, "creating sd_root/atmosphere/titles/0100000000000032",
		os.MkdirAll("sd_root/atmosphere/titles/0100000000000032", 0700))

	errCheck(&w, "copying set.mitm's exefs", copyFile("build/atmosphere/stratosphere/set_mitm/set_mitm.nsp",
		"sd_root/atmosphere/titles/0100000000000032/exefs.nsp"))

	errCheck(&w, "copying sm", copyFile("build/atmosphere/stratosphere/sm/sm.kip", "sd_root/cfw/sm.kip"))

	hekateConfig := []string{}
	c := []string{}

	if nogc && inArray(features, "4") {
		hekateConfig = []string{
			"[Stock]",
			"",
			"[Stock (nogc)]",
			"kip1patch=nogc",
			"",
			"[CFW]",
		}
		c = []string{"kip1=cfw/*", "kip1patch=nogc,nosigchk", "secmon=cfw/exosphere.bin"}
	} else if nogc {
		hekateConfig = []string{
			"[Stock]",
			"",
			"[Stock (nogc)]",
			"kip1patch=nogc",
			"",
			"[CFW]",
		}
		c = []string{"kip1=cfw/*", "kip1patch=nogc", "secmon=cfw/exosphere.bin"}
	} else if inArray(features, "4") {
		hekateConfig = []string{
			"[Stock]",
			"",
			"[CFW]",
		}
		c = []string{"kip1=cfw/*", "kip1patch=nosigchk", "secmon=cfw/exosphere.bin"}
	} else {
		hekateConfig = []string{
			"[Stock]",
			"",
			"[CFW]",
		}
		c = []string{"kip1=cfw/*", "secmon=cfw/exosphere.bin"}
	}

	if inArray(features, "1") || inArray(features, "4") {
		if !inArray(features, "no-hbl") {
			fmt.Fprintf(w, "cloning hbloader...\n")
			_, err = git.PlainClone("build/hbloader", false, &git.CloneOptions{
				URL: "https://github.com/switchbrew/nx-hbloader.git",
			})
			if err != nil && err.Error() != "repository already exists" {
				errCheck(&w, "cloning hbloader", err)
			}

			fmt.Fprintf(w, "building hbloader...\n")
			cmd := exec.Command("make", "-j")
			cmd.Dir = "build/hbloader"
			errCheck(&w, "building hbloader", cmd.Run())

			fmt.Fprintf(w, "copying files...\n")
			errCheck(&w, "copying hbloader", copyFile("build/hbloader/hbl.nsp", "sd_root/atmosphere/hbl.nsp"))
		}

		fmt.Fprintf(w, "cloning hbmenu...\n")
		_, err = git.PlainClone("build/hbmenu", false, &git.CloneOptions{
			URL: "https://github.com/switchbrew/nx-hbmenu.git",
		})
		if err != nil && err.Error() != "repository already exists" {
			errCheck(&w, "cloning hbmenu", err)
		}

		fmt.Fprintf(w, "building hbmenu...\n")
		cmd = exec.Command("make", "nx", "-j")
		cmd.Dir = "build/hbmenu"
		errCheck(&w, "building hbmenu", cmd.Run())

		fmt.Fprintf(w, "copying files...\n")
		errCheck(&w, "copying hbmenu", copyFile("build/hbmenu/hbmenu.nro", "sd_root/hbmenu.nro"))
	}

	if inArray(features, "2") {
		fmt.Fprintf(w, "copying files...\n")
		errCheck(&w, "copying fs_mitm (layeredfs)", copyFile("build/atmosphere/stratosphere/fs_mitm/fs_mitm.kip",
			"sd_root/cfw/fs_mitm.kip"))
		c = append(c, "atmosphere=1")
	}

	if inArray(features, "3") {
		errCheck(&w, "creating sd_root/atmosphere/exefs_patches", os.MkdirAll("sd_root/atmosphere/exefs_patches", 0700))
		fmt.Fprintf(w, "copying files...\n")
		errCheck(&w, "copying fake_tickets (sigpatches)", copyFolder("fake_tickets",
			"sd_root/atmosphere/exefs_patches/fake_tickets"))
	}

	if inArray(features, "4") {
		fmt.Fprintf(w, "cloning sys-ftpd...\n")
		_, err = git.PlainClone("build/sys-ftpd", false, &git.CloneOptions{
			URL: "https://github.com/jakibaki/sys-ftpd.git",
		})
		if err != nil && err.Error() != "repository already exists" {
			errCheck(&w, "cloning sys-ftpd", err)
		}

		fmt.Fprintf(w, "building sys-ftpd...\n")
		cmd := exec.Command("make", "-j")
		cmd.Dir = "build/sys-ftpd"
		errCheck(&w, "building sys-ftpd", cmd.Run())

		fmt.Fprintf(w, "copying files...\n")
		errCheck(&w, "copying sys-ftpd's sound files", copyFolder("build/sys-ftpd/sd_card/ftpd", "sd_root/ftpd"))
		errCheck(&w, "copying sys-ftpd", copyFile("build/sys-ftpd/sys-ftpd.kip", "sd_root/cfw/sys-ftpd.kip"))
	}

	if inArray(features, "5") {
		fmt.Fprintf(w, "cloning tinfoil...\n")
		_, err = git.PlainClone("build/tinfoil", false, &git.CloneOptions{
			URL: "https://github.com/XorTroll/Tinfoil.git",
		})
		if err != nil && err.Error() != "repository already exists" {
			errCheck(&w, "cloning tinfoil", err)
		}

		fmt.Fprintf(w, "building tinfoil...\n")
		cmd := exec.Command("make", "-j")
		cmd.Dir = "build/tinfoil"
		errCheck(&w, "building tinfoil", cmd.Run())

		fmt.Fprintf(w, "copying files...\n")
		errCheck(&w, "creating sd_root/switch", os.MkdirAll("sd_root/switch", 0700))
		errCheck(&w, "copying tinfoil", copyFile("build/tinfoil/tinfoil.nro", "sd_root/switch/Tinfoil.nro"))
	}

	fmt.Fprintf(w, "creating hekate config...\n\n")
	sort.Strings(c)
	for _, v := range c {
		hekateConfig = append(hekateConfig, v)
	}

	f, err := os.Create("sd_root/bootloader/hekate_ipl.ini")
	errCheck(&w, "creating sd_root/bootloader/hekate_ipl.ini", err)
	for i, v := range hekateConfig {
		if i+1 == len(hekateConfig) {
			_, err = f.WriteString(v)
		} else {
			_, err = f.WriteString(v + "\n")
		}
		errCheck(&w, "writing to hekate_ipl.ini", err)
	}
	f.Close()

	fmt.Fprintf(w, "done!")
}
