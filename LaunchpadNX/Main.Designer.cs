namespace LaunchpadNX
{
    partial class LaunchpadNX
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.atmsCheckbox = new System.Windows.Forms.CheckBox();
            this.lfsCheckbox = new System.Windows.Forms.CheckBox();
            this.hbmenuCheckbox = new System.Windows.Forms.CheckBox();
            this.hbmenuTitleSelect = new System.Windows.Forms.ComboBox();
            this.sigpatchesCheckbox = new System.Windows.Forms.CheckBox();
            this.tinfoilCheckbox = new System.Windows.Forms.CheckBox();
            this.startButton = new System.Windows.Forms.Button();
            this.SuspendLayout();
            // 
            // atmsCheckbox
            // 
            this.atmsCheckbox.AutoSize = true;
            this.atmsCheckbox.Checked = true;
            this.atmsCheckbox.CheckState = System.Windows.Forms.CheckState.Checked;
            this.atmsCheckbox.Cursor = System.Windows.Forms.Cursors.Default;
            this.atmsCheckbox.Enabled = false;
            this.atmsCheckbox.Location = new System.Drawing.Point(12, 12);
            this.atmsCheckbox.Name = "atmsCheckbox";
            this.atmsCheckbox.Size = new System.Drawing.Size(109, 17);
            this.atmsCheckbox.TabIndex = 0;
            this.atmsCheckbox.Text = "Atmosphere Base";
            this.atmsCheckbox.UseVisualStyleBackColor = true;
            // 
            // lfsCheckbox
            // 
            this.lfsCheckbox.AutoSize = true;
            this.lfsCheckbox.Location = new System.Drawing.Point(127, 12);
            this.lfsCheckbox.Name = "lfsCheckbox";
            this.lfsCheckbox.Size = new System.Drawing.Size(77, 17);
            this.lfsCheckbox.TabIndex = 1;
            this.lfsCheckbox.Text = "LayeredFS";
            this.lfsCheckbox.UseVisualStyleBackColor = true;
            // 
            // hbmenuCheckbox
            // 
            this.hbmenuCheckbox.AutoSize = true;
            this.hbmenuCheckbox.Location = new System.Drawing.Point(210, 12);
            this.hbmenuCheckbox.Name = "hbmenuCheckbox";
            this.hbmenuCheckbox.Size = new System.Drawing.Size(64, 17);
            this.hbmenuCheckbox.TabIndex = 2;
            this.hbmenuCheckbox.Text = "hbmenu";
            this.hbmenuCheckbox.UseVisualStyleBackColor = true;
            this.hbmenuCheckbox.CheckedChanged += new System.EventHandler(this.hbmenuCheckbox_CheckedChanged);
            // 
            // hbmenuTitleSelect
            // 
            this.hbmenuTitleSelect.Enabled = false;
            this.hbmenuTitleSelect.FormattingEnabled = true;
            this.hbmenuTitleSelect.Items.AddRange(new object[] {
            "Album",
            "Controllers Screen",
            "eShop"});
            this.hbmenuTitleSelect.Location = new System.Drawing.Point(279, 10);
            this.hbmenuTitleSelect.Name = "hbmenuTitleSelect";
            this.hbmenuTitleSelect.Size = new System.Drawing.Size(97, 21);
            this.hbmenuTitleSelect.TabIndex = 3;
            this.hbmenuTitleSelect.Text = "Album";
            // 
            // sigpatchesCheckbox
            // 
            this.sigpatchesCheckbox.AutoSize = true;
            this.sigpatchesCheckbox.Location = new System.Drawing.Point(12, 35);
            this.sigpatchesCheckbox.Name = "sigpatchesCheckbox";
            this.sigpatchesCheckbox.Size = new System.Drawing.Size(79, 17);
            this.sigpatchesCheckbox.TabIndex = 4;
            this.sigpatchesCheckbox.Text = "Sigpatches";
            this.sigpatchesCheckbox.UseVisualStyleBackColor = true;
            // 
            // tinfoilCheckbox
            // 
            this.tinfoilCheckbox.AutoSize = true;
            this.tinfoilCheckbox.Location = new System.Drawing.Point(127, 35);
            this.tinfoilCheckbox.Name = "tinfoilCheckbox";
            this.tinfoilCheckbox.Size = new System.Drawing.Size(54, 17);
            this.tinfoilCheckbox.TabIndex = 5;
            this.tinfoilCheckbox.Text = "Tinfoil";
            this.tinfoilCheckbox.UseVisualStyleBackColor = true;
            // 
            // startButton
            // 
            this.startButton.Location = new System.Drawing.Point(12, 58);
            this.startButton.Name = "startButton";
            this.startButton.Size = new System.Drawing.Size(364, 23);
            this.startButton.TabIndex = 7;
            this.startButton.Text = "Start!";
            this.startButton.UseVisualStyleBackColor = true;
            this.startButton.Click += new System.EventHandler(this.startButton_Click);
            // 
            // LaunchpadNX
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(386, 89);
            this.Controls.Add(this.startButton);
            this.Controls.Add(this.tinfoilCheckbox);
            this.Controls.Add(this.sigpatchesCheckbox);
            this.Controls.Add(this.hbmenuTitleSelect);
            this.Controls.Add(this.hbmenuCheckbox);
            this.Controls.Add(this.lfsCheckbox);
            this.Controls.Add(this.atmsCheckbox);
            this.Name = "LaunchpadNX";
            this.Text = "LaunchpadNX";
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.CheckBox atmsCheckbox;
        private System.Windows.Forms.CheckBox lfsCheckbox;
        private System.Windows.Forms.CheckBox hbmenuCheckbox;
        private System.Windows.Forms.ComboBox hbmenuTitleSelect;
        private System.Windows.Forms.CheckBox sigpatchesCheckbox;
        private System.Windows.Forms.CheckBox tinfoilCheckbox;
        private System.Windows.Forms.Button startButton;
    }
}

