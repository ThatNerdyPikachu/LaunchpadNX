using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Windows.Forms;

namespace LaunchpadNX
{
    public partial class LaunchpadNX : Form
    {
        public LaunchpadNX()
        {
            InitializeComponent();
            this.MaximizeBox = false;
        }

        private void hbmenuCheckbox_CheckedChanged(object sender, EventArgs e)
        {
            if (hbmenuCheckbox.Checked)
            {
                hbmenuTitleSelect.Enabled = true;
                tinfoilCheckbox.Enabled = true;
                checkpointCheckbox.Enabled = true;
                hbasCheckbox.Enabled = true;
                switchidentCheckbox.Enabled = true;
            } else
            {
                hbmenuTitleSelect.Enabled = false;
                tinfoilCheckbox.Enabled = false;
                checkpointCheckbox.Enabled = false;
                hbasCheckbox.Enabled = false;
                switchidentCheckbox.Enabled = false;
            }
        }

        private void RunCommand(string command)
        {
            ProcessStartInfo proc = new ProcessStartInfo("cmd.exe")
            {
                Arguments = "/C " + command
            };
            Process cmd = Process.Start(proc);
            cmd.WaitForExit();
        }

        // thanks stack overflow
        private static void CopyDirectory(string sourceDirName, string destDirName, bool copySubDirs)
        {
            DirectoryInfo dir = new DirectoryInfo(sourceDirName);
            DirectoryInfo[] dirs = dir.GetDirectories();

            // If the source directory does not exist, throw an exception.
            if (!dir.Exists)
            {
                throw new DirectoryNotFoundException(
                    "Source directory does not exist or could not be found: "
                    + sourceDirName);
            }

            // If the destination directory does not exist, create it.
            if (!Directory.Exists(destDirName))
            {
                Directory.CreateDirectory(destDirName);
            }


            // Get the file contents of the directory to copy.
            FileInfo[] files = dir.GetFiles();

            foreach (FileInfo file in files)
            {
                // Create the path to the new copy of the file.
                string temppath = Path.Combine(destDirName, file.Name);

                // Copy the file.
                file.CopyTo(temppath, false);
            }

            // If copySubDirs is true, copy the subdirectories.
            if (copySubDirs)
            {

                foreach (DirectoryInfo subdir in dirs)
                {
                    // Create the subdirectory.
                    string temppath = Path.Combine(destDirName, subdir.Name);

                    // Copy the subdirectories.
                    CopyDirectory(subdir.FullName, temppath, copySubDirs);
                }
            }
        }

        private void startButton_Click(object sender, EventArgs e)
        {
            // do a bit of cleanup
            if (Directory.Exists("SD Root"))
            {
                RunCommand("rmdir /S /Q \"SD Root\"");
            }

            if (Directory.Exists("temp"))
            {
                RunCommand("rmdir /S /Q temp");
            }

            if (File.Exists("CFW.bin"))
            {
                File.Delete("CFW.bin");
            }

            // create hekate config list
            List<string> hekateConfig = new List<string>()
            {
                "[Stock]",
                "",
                "[CFW]",
                "atmosphere=1",
                "secmon=cfw/exosphere.bin",
                "kip1=cfw/loader.kip",
                "kip1=cfw/pm.kip",
                "kip1=cfw/sm.kip"
            };

            // create "SD Root"
            Directory.CreateDirectory("SD Root");

            // install switch dev tools
            RunCommand("pacman -S switch-dev devkitARM --noconfirm --needed");
            
            // Hekate
            //
            // clone it
            RunCommand("git clone https://github.com/CTCaer/hekate.git temp\\hekate");
            
            // remove sample module
            RunCommand("rmdir /S /Q temp\\hekate\\modules\\simple_sample");
            
            // build it
            RunCommand("cd temp\\hekate && make -j");

            // create needed directories
            Directory.CreateDirectory("SD Root\\bootloader");
            Directory.CreateDirectory("SD Root\\bootloader\\ini");
            Directory.CreateDirectory("SD Root\\bootloader\\payloads");
            Directory.CreateDirectory("SD Root\\bootloader\\sys");

            // copy files
            File.Copy("temp\\hekate\\output\\hekate.bin", "CFW.bin", true);
            File.Copy("temp\\hekate\\output\\libsys_lp0.bso", "SD Root\\bootloader\\libsys_lp0.bso");

            // Atmosphere base (always ran)
            //
            // clone it
            RunCommand("git clone https://github.com/Atmosphere-NX/Atmosphere.git temp\\Atmosphere");

            if (hbmenuCheckbox.Checked)
            {
                if (hbmenuTitleSelect.Text == "Album")
                {
                    RunCommand("cd temp\\Atmosphere && git apply ../../files/hbl-album.patch");
                } else if (hbmenuTitleSelect.Text == "Controllers Screen")
                {
                    RunCommand("cd temp\\Atmosphere && git apply ../../files/hbl-controllers.patch");
                } else
                {
                    RunCommand("cd temp\\Atmosphere && git apply ../../files/hbl-album.patch");
                }
            }

            // build it
            RunCommand("cd temp\\Atmosphere\\exosphere && make -j && cd ../stratosphere && make -j");

            // create needed directories
            Directory.CreateDirectory("SD Root\\cfw");
            Directory.CreateDirectory("SD Root\\atmosphere\\titles\\0100000000000036\\exefs");

            // copy files
            File.Copy("temp\\atmosphere\\exosphere\\exosphere.bin", "SD Root\\cfw\\exosphere.bin");
            File.Copy("temp\\Atmosphere\\stratosphere\\creport\\creport.nso", "SD Root\\atmosphere\\titles\\0100000000000036\\exefs\\main");
            File.Copy("temp\\Atmosphere\\stratosphere\\creport\\creport.npdm", "SD Root\\atmosphere\\titles\\0100000000000036\\exefs\\main.npdm");
            File.Copy("temp\\Atmosphere\\stratosphere\\loader\\loader.kip", "SD Root\\cfw\\loader.kip");
            File.Copy("temp\\Atmosphere\\stratosphere\\pm\\pm.kip", "SD Root\\cfw\\pm.kip");
            File.Copy("temp\\Atmosphere\\stratosphere\\sm\\sm.kip", "SD Root\\cfw\\sm.kip");

            // LayeredFS
            if (lfsCheckbox.Checked)
            {
                File.Copy("temp\\Atmosphere\\stratosphere\\fs_mitm\\fs_mitm.kip", "SD Root\\cfw\\fs_mitm.kip");
                hekateConfig.Add("kip1=cfw/fs_mitm.kip");
            }

            // hbmenu!
            if (hbmenuCheckbox.Checked)
            {        
                // hbloader
                //
                // clone it
                RunCommand("git clone https://github.com/switchbrew/nx-hbloader.git temp\\nx-hbloader");

                // copy the right json (defaults to album unless otherwise stated)
                string tid = "010000000000100D";
                if (hbmenuTitleSelect.Text == "Album")
                {
                    File.Copy("files\\hbl-album.json", "temp\\nx-hbloader\\hbl.json", true);
                    tid = "010000000000100D";
                } else if (hbmenuTitleSelect.Text == "Controllers Screen")
                {
                    File.Copy("files\\hbl-controllers.json", "temp\\nx-hbloader\\hbl.json", true);
                    tid = "0100000000001003";
                } else
                {
                    File.Copy("files\\hbl-album.json", "temp\\nx-hbloader\\hbl.json", true);
                }

                // build it
                RunCommand("cd temp\\nx-hbloader && make -j");

                // create needed directories
                Directory.CreateDirectory("SD Root\\atmosphere\\titles\\" + tid + "\\exefs");

                // copy files
                File.Copy("temp\\nx-hbloader\\hbl.nso", "SD Root\\atmosphere\\titles\\" + tid + "\\exefs\\main");
                File.Copy("temp\\nx-hbloader\\hbl.npdm", "SD Root\\atmosphere\\titles\\" + tid + "\\exefs\\main.npdm");
                RunCommand("cd SD Root\\atmosphere\\titles\\" + tid + "\\exefs && touch rtld.stub");


                // actual hbmenu
                //
                // install freetype
                RunCommand("pacman -S switch-freetype --noconfirm --needed");

                // clone it
                RunCommand("git clone https://github.com/switchbrew/nx-hbmenu.git temp\\nx-hbmenu");

                // build it
                RunCommand("cd temp\\nx-hbmenu && make nx -j");

                // copy file
                File.Copy("temp\\nx-hbmenu\\nx-hbmenu.nro", "SD Root\\hbmenu.nro");
            }

            // sigpatches!
            if (sigpatchesCheckbox.Checked)
            {
                // create needed directory
                Directory.CreateDirectory("SD Root\\atmosphere\\exefs_patches");

                // copy directory
                CopyDirectory("files\\Fake Tickets", "SD Root\\atmosphere\\exefs_patches\\Fake Tickets", false);

                hekateConfig.Add("kip1patch=nosigchk");
            }

            // tinfoil!
            if (tinfoilCheckbox.Checked)
            {
                // install curl
                RunCommand("pacman -S switch-curl --noconfirm --needed");

                // clone it
                RunCommand("git clone https://github.com/Adubbz/Tinfoil.git temp\\Tinfoil");

                // build it
                RunCommand("cd temp\\Tinfoil && make -j");

                // create needed directories
                Directory.CreateDirectory("SD Root\\switch");
                Directory.CreateDirectory("SD Root\\tinfoil\\extracted");
                Directory.CreateDirectory("SD Root\\tinfoil\\nsp");
                Directory.CreateDirectory("SD Root\\tinfoil\\ticket");

                // copy file
                File.Copy("temp\\Tinfoil\\Tinfoil.nro", "SD Root\\switch\\Tinfoil.nro");
            }

            // checkpoint!
            if (checkpointCheckbox.Checked)
            {
                // install freetype
                RunCommand("pacman -S switch-freetype --noconfirm --needed");

                // clone it
                RunCommand("git clone https://github.com/BernardoGiordano/Checkpoint.git temp\\Checkpoint");

                // build it
                RunCommand("cd temp\\Checkpoint\\switch && make -j");

                // create needed directory
                if (!Directory.Exists("SD Root\\switch\\Checkpoint"))
                {
                    Directory.CreateDirectory("SD Root\\switch\\Checkpoint");
                }

                // copy file
                File.Copy("temp\\Checkpoint\\switch\\out\\Checkpoint.nro", "SD Root\\switch\\Checkpoint\\Checkpoint.nro");
            }

            // sys-ftpd!
            if (ftpdCheckbox.Checked)
            {
                // install mpg123
                RunCommand("pacman -S switch-mpg123 --noconfirm --needed");

                // clone it
                RunCommand("git clone https://github.com/jakibaki/sys-ftpd.git temp\\sys-ftpd");

                // build it
                RunCommand("cd temp\\sys-ftpd && make -j");

                // copy files
                File.Copy("temp\\sys-ftpd\\sys-ftpd.kip", "SD Root\\cfw\\sys-ftpd.kip");
                CopyDirectory("temp\\sys-ftpd\\sd_card", "SD Root", true);

                hekateConfig.Add("kip1=cfw/sys-ftpd.kip");
            }

            // hbas
            if (hbasCheckbox.Checked)
            {
                // install deps
                RunCommand("pacman -S switch-curl switch-bzip2 switch-freetype switch-libjpeg-turbo switch-sdl2 " +
                    "switch-sdl2_gfx switch-sdl2_image switch-sdl2_ttf switch-zlib " +
                    "switch-libpng --noconfirm --needed");

                // clone it
                RunCommand("git clone https://github.com/vgmoose/appstorenx.git temp\\appstorenx --recursive");

                // build it
                RunCommand("cd temp\\appstorenx && make -j");

                // create needed directory
                if (!Directory.Exists("SD Root\\switch\\appstore"))
                {
                    Directory.CreateDirectory("SD Root\\switch\\appstore");
                }

                // copy files
                File.Copy("temp\\appstorenx\\appstore.nro", "SD Root\\switch\\appstore\\appstore.nro");
                CopyDirectory("temp\\appstorenx\\res", "SD Root\\switch\\appstore\\res", false);
            }

            // switchident
            if (switchidentCheckbox.Checked)
            {

                // clone it
                RunCommand("git clone https://github.com/joel16/SwitchIdent.git temp\\SwitchIdent");

                // build it
                RunCommand("cd temp\\SwitchIdent\\console && make -j");

                // create needed directory
                if (!Directory.Exists("SD Root\\switch"))
                {
                    Directory.CreateDirectory("SD Root\\switch");
                }

                // copy file
                File.Copy("temp\\SwitchIdent\\console\\console.nro", "SD Root\\switch\\SwitchIdent.nro");
            }

            // hekate config (THANKS C#!!)
            System.IO.File.WriteAllLines("SD Root\\bootloader\\hekate_ipl.ini", hekateConfig);

            // cleanup
            RunCommand("rmdir /S /Q temp");

            string done = "Done! Copy the contents of the \"SD Root\" folder to your SD Card, then launch the CFW.bin payload using your favorite method!";
            if (hbmenuCheckbox.Checked)
            {
                done = done + " (Since you selected hbmenu, you can load it by holding R while launching the " + hbmenuTitleSelect.Text + " from the HOME Menu)";
            }

            // we are done!
            MessageBox.Show(done, 
                "LaunchpadNX", MessageBoxButtons.OK, MessageBoxIcon.Information);
        }
    }
}
