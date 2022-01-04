package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	InfoColor    = "\033[1;34m"
	NoticeColor  = "\033[1;36m"
	WarningColor = "\033[1;33m"
	ErrorColor   = "\033[1;31m"
	DebugColor   = "\033[0;36m"
	RESET        = "\033[0m"
)

func fatal(error string) {
	fmt.Printf(ErrorColor)
	log.Fatal(error, RESET)
}
func info(info string) {
	fmt.Printf(InfoColor)
	log.Println(info, RESET)
}
func debug(debug string) {
	fmt.Printf(DebugColor)
	log.Println(debug, RESET)
}
func success(success string) {
	fmt.Printf(NoticeColor)
	log.Println(success, RESET)
}
func fatalerr(err error) {
	if err != nil {
		fatal(err.Error())
	}
	success("done")
}
func chroot_exec(chroot_dir, arg string) {
	val, err := exec.Command("su", "-c", "chroot "+chroot_dir+" "+arg).CombinedOutput()
	fmt.Println("su", "-c", "chroot "+chroot_dir+" "+arg)
	fmt.Println(string(val))
	fatalerr(err)
}
func prompt(s string) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		return strings.TrimSpace(response)
	}
}

func main() {
	fmt.Println("Linux for Phone")
	distro := flag.String("distro", "blackarch", "Write distribution name to install. arch, blackarch, debian and kali-linux are supported")
	arch := flag.String("arch", "aarch64", "['aarch64','aarch32','x86','x64'] write arch. Use only compatable one")
	kernel := flag.String("kernel", "important", "Liux kernel for phone")
	linux := flag.String("linux", "important", "linux directory made on phone")
	ndk := flag.String("ndk", "important", "android-ndk directory made on phone")
	firmware := flag.String("firmware", "", "firmware directory for installing it on the system")
	deconfig := flag.String("deconfig", "", "deconfig name for make")
	offer_deps := flag.String("offer_dependencies", "n", "deconfig name for make")
	flag.Parse()
	if *kernel == "important" || *linux == "important" || *ndk == "important" {
		fatal("kernel and linux is required")
	}
	fmt.Println(*linux, *kernel)
	info("Building Kernel")
	if *distro == "blackarch" || *distro == "arch" {
		if *offer_deps != "n" {
			debug("Installing Dependencies")
			fmt.Println("/bin/sh", "-c", "yay -S qemu-user-static mkbootimg libselinux-static android-tools android-sdk-platform-tools --noconfirm")
			fmt.Println("/bin/sh", "-c", "sudo pacman -S qemu-headless android-ndk  --noconfirm --needed")
			_, err := exec.Command("sudo", "ln", "/usr/lib/libselinux.so.1", "/usr/lib/libselinux.so.0").CombinedOutput()
			fatalerr(err)
		}
		debug("Setting up permissions")
		_, err := exec.Command("/bin/sh", "-c", "sudo chmod -R 0755 "+*linux+"/usr/bin/*").CombinedOutput()
		fatalerr(err)
		_, err = exec.Command("/bin/sh", "-c", "sudo chown -R root:root "+*linux+"/usr/bin/*").CombinedOutput()
		fatalerr(err)
		_, err = exec.Command("/bin/sh", "-c", "sudo chmod -R 0744 "+*linux+"/etc/*").CombinedOutput()
		_, err = exec.Command("/bin/sh", "-c", "sudo chown -R root:root "+*linux+"/etc/*").CombinedOutput()
		_, err = exec.Command("/bin/sh", "-c", "sudo chmod -R 0755 "+*linux+"/usr/lib/*").CombinedOutput()
		fatalerr(err)
		_, err = exec.Command("/bin/sh", "-c", "sudo chown -R root:root "+*linux+"/usr/lib/*").CombinedOutput()
		fatalerr(err)
		debug("Setting up emulation")
		_, err = exec.Command("sudo", "cp", "/usr/bin/qemu-arm-static", *linux+"/usr/bin/").CombinedOutput()
		fatalerr(err)
		debug("Setting Eviroment")
		err = os.Setenv("LC_ALL", "C")
		fatalerr(err)
		err = os.Setenv("LANGUAGE", "C")
		fatalerr(err)
		err = os.Setenv("LANG", "C")
		fatalerr(err)
		debug("Updatinng repos")
		chroot_exec(*linux, "pacman -Sy")
		debug("Installing X")
		chroot_exec(*linux, "pacman -S xserver–xorg–video–fbdev xserver–xorg–input–evdev initramfs–tools wpa-supplicant android-tools --noconfirm --needed")
		if *firmware != "" {
			debug("Setting up fimware")
			chroot_exec(*linux, "addgroup ––gid 3003 inet")
			err = os.MkdirAll(*linux+"/lib/firmware", 0777)
			fatalerr(err)
			_, err = exec.Command("sudo", "cp", *firmware+"/*", *linux+"/lib/firmware").CombinedOutput()
		}
		debug("Setting Eviroment For Build")
		if *arch == "aarch64" {
			err = os.Setenv("CROSS_COMPILE", *ndk+"/toolchains/aarch64-linux-android-4.9/prebuilt/linux-x86_64/bin/aarch64-linux-android-")
			fatalerr(err)
			err = os.Setenv("ARCH", "arm64")
			fatalerr(err)
		} else if *arch == "aarch32" {
			err = os.Setenv("CROSS_COMPILE", *ndk+"/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/bin/arm-linux-androideabi-")
			fatalerr(err)
			err = os.Setenv("ARCH", "arm")
			fatalerr(err)
		} else if *arch == "x86" {
			err = os.Setenv("CROSS_COMPILE", *ndk+"/toolchains/x86-4.9/prebuilt/linux-x86_64/bin/i686-linux-android-")
			fatalerr(err)
			err = os.Setenv("ARCH", "i386")
			fatalerr(err)
		} else if *arch == "x64" {
			err = os.Setenv("CROSS_COMPILE", *ndk+"/toolchains/x86_64-4.9/prebuilt/linux-x86_64/bin/x86_64-linux-android-")
			fatalerr(err)
			err = os.Setenv("ARCH", "x86_64")
			fatalerr(err)
		}
		err = os.Setenv("INSTALL_PATH", *linux+"/boot")
		fatalerr(err)
		err = os.Setenv("INSTALL_MOD_PATH", *linux)
		fatalerr(err)
		debug("Building kernel")
		kernel_build := exec.Command("/bin/bash", "-c", "cd "+*kernel+" && make clean && make mrproper")
		kernel_build.Stdin = os.Stdin
		kernel_build.Stdout = os.Stdout
		kernel_build.Stderr = os.Stderr
		err = kernel_build.Run()
		if *deconfig != "" {
			kernel_build = exec.Command("/bin/bash", "-c", "cd "+*kernel+" && make "+*deconfig)
			kernel_build.Stdin = os.Stdin
			kernel_build.Stdout = os.Stdout
			kernel_build.Stderr = os.Stderr
			err = kernel_build.Run()
			fatalerr(err)
		}
		fmt.Println(`go to Device Drivers > Character Devices > and enable "Virtual Terminal"`, "\n", `go to Device Drivers > Graphics Support > Console Display Driver support > enable "Framebuffer Console Support"`)
		s := prompt("Have you read instructions [y/n/q]: ")
		for s != "y" {
			if s == "q" {
				os.Exit(1)
			}
			fmt.Println(`go to Device Drivers > Character Devices > and enable "Virtual Terminal"`, "\n", `go to Device Drivers > Graphics Support > Console Display Driver support > enable "Framebuffer Console Support"`)
			s = prompt("Have you read instructions [y/n/q]: ")
		}
		nproc, err := exec.Command("nproc", "--all").CombinedOutput()
		fatalerr(err)
		kernel_build = exec.Command("/bin/bash", "-c", "cd "+*kernel+" && make menuconfig && make –j"+string(nproc)+" && sudo make modules_install")
		kernel_build.Stdin = os.Stdin
		kernel_build.Stdout = os.Stdout
		kernel_build.Stderr = os.Stderr
		err = kernel_build.Run()
		fatalerr(err)
		debug("Getting zImage and System.map")
		_, err = exec.Command("/bin/bash", "-c", "cd "+*kernel+" && sudo cp arch/arm/boot/zImage "+*linux+"/boot").CombinedOutput()
		fatalerr(err)
		_, err = exec.Command("/bin/bash", "-c", "cd "+*kernel+" && cp System.map "+*linux+"/boot").CombinedOutput()
		fatalerr(err)
		_, err = exec.Command("sudo", "cp", "./conf/com.conf", *linux+"/linux/root/").CombinedOutput()
		fatalerr(err)
		debug("Generating initrd.img")
		chroot_exec(*linux, "mkinitcpio --generate /boot/initrd.img.gz --kernel `ls /lib/modules`")
		debug("Generating system.img")
		_, err = exec.Command("sudo", "./bin/make_ext4fs", "-s", "-l", "3096M", "s.img "+*linux).CombinedOutput()
		fatalerr(err)
		_, err = exec.Command("simg2img", "s.img", "system.img").CombinedOutput()
		fatalerr(err)
		debug("Generating boot.img")
		_, err = exec.Command("/bin/sh", "-c", "mkbootimg --kernel "+*linux+"/boot/zImage --ramdisk "+*linux+`/boot/initrd.img.gz --base "`+"`cat out/boot.img–base`"+`" --cmdline "root=/dev/mmcblk0p2 console=tty0" –o boot.img`).CombinedOutput()
		fatalerr(err)
		debug("flashing system.img")
		_, err = exec.Command("fastboot", "boot", "system", "system.img").CombinedOutput()
		fatalerr(err)
		debug("booting boot.img")
		_, err = exec.Command("fastboot", "boot", "boot", "boot.img").CombinedOutput()
		fatalerr(err)
	}
}
