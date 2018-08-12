using System;
using System.IO;
using System.Windows.Forms;

namespace LaunchpadNX
{
    public partial class LaunchpadNX : Form
    {
        private int optionsChecked;

        public LaunchpadNX()
        {
            InitializeComponent();
            optionsChecked = 1;
        }

        private void lfsCheckbox_CheckedChanged(object sender, EventArgs e)
        {
            if (lfsCheckbox.Checked)
            {
                optionsChecked += 1;
            } else
            {
                optionsChecked -= 1;
            }
        }

        private void hbmenuCheckbox_CheckedChanged(object sender, EventArgs e)
        {
            if (hbmenuCheckbox.Checked)
            {
                optionsChecked += 1;
                hbmenuTitleSelect.Enabled = true;
            } else
            {
                optionsChecked -= 1;
                hbmenuTitleSelect.Enabled = false;
            }
        }

        private void sigpatchesCheckbox_CheckedChanged(object sender, EventArgs e)
        {
            if (sigpatchesCheckbox.Checked)
            {
                optionsChecked += 1;
            } else
            {
                optionsChecked -= 1;
            }
        }

        private void tinfoilCheckbox_CheckedChanged(object sender, EventArgs e)
        {
            if (tinfoilCheckbox.Checked)
            {
                optionsChecked += 1;
            } else
            {
                optionsChecked -= 1;
            }
        }
    }
}
