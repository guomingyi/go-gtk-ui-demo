package main

import (
	"fmt"
	"os"
	"os/exec"
	//"path/filepath"
	"bufio"
	"io"
	"io/ioutil"

	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

func uniq(strings []string) (ret []string) {
	return
}

func authors() []string {
	if b, err := exec.Command("git", "log").Output(); err == nil {
		lines := strings.Split(string(b), "\n")

		var a []string
		r := regexp.MustCompile(`^Author:\s*([^ <]+).*$`)
		for _, e := range lines {
			ms := r.FindStringSubmatch(e)
			if ms == nil {
				continue
			}
			a = append(a, ms[1])
		}
		sort.Strings(a)
		var p string
		lines = []string{}
		for _, e := range a {
			if p == e {
				continue
			}
			lines = append(lines, e)
			p = e
		}
		return lines
	}
	return []string{"Yasuhiro Matsumoto <mattn.jp@gmail.com>"}
}

var textview_output *gtk.TextView
var btn_load *gtk.Button
var btn_start *gtk.Button
var buttons *gtk.HBox
var progressbar *gtk.ProgressBar

func ui_show() {
	gtk.Init(&os.Args)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle(" ")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		fmt.Println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	// define
	var textview_file1 *gtk.TextView
	var textview_file2 *gtk.TextView
	var textview_file3 *gtk.TextView

	var windos_w int = 800
	var windos_h int = 700

	currentPath = getCurrDir()

	//--------------------------------------------------------
	// GtkVBox
	//--------------------------------------------------------
	vbox := gtk.NewVBox(false, 1)
	//--------------------------------------------------------
	// GtkVPaned
	//--------------------------------------------------------
	vpaned := gtk.NewVPaned()

	//--------------------------------------------------------
	// GtkFrame
	//--------------------------------------------------------
	frame1 := gtk.NewFrame("")
	framebox1 := gtk.NewVBox(false, 1)

	frame2 := gtk.NewFrame("")
	framebox2 := gtk.NewVBox(false, 1)

	//-------------------------------framebox1-------start------------------------
	label := gtk.NewLabel("JusetSafe v1.0")
	label.ModifyFontEasy("DejaVu Serif 15")
	framebox1.PackStart(label, false, true, 0)

	// GtkVSeparator
	vsep := gtk.NewVSeparator()
	vsep.Activate()
	vsep.SetSizeRequest(10, 10)
	framebox1.PackStart(vsep, false, false, 0)

	//--------------------------------------------------------
	// GtkHBox
	//--------------------------------------------------------
	buttons = gtk.NewHBox(false, 1)

	//--------------------------------------------------------
	// GtkButton
	//--------------------------------------------------------
	btn_load = gtk.NewButtonWithLabel("加载固件")
	//btn_load.SetSizeRequest(50, 50)
	buttons.Add(btn_load)

	btn_start = gtk.NewButtonWithLabel("开始")
	btn_start.SetSensitive(false)
	//btn_start.SetSizeRequest(50, 50)
	buttons.Add(btn_start)

	//--------------------------------------------------------
	// GtkComboBox
	//--------------------------------------------------------
	combobox := gtk.NewComboBoxText()
	combobox.AppendText("下载")
	combobox.AppendText("测试")
	combobox.SetActive(0)
	buttons.Add(combobox)
	framebox1.PackStart(buttons, false, false, 0)

	// GtkVSeparator
	vsep0 := gtk.NewVSeparator()
	vsep0.Activate()
	vsep0.SetSizeRequest(5, 5)
	framebox1.PackStart(vsep0, false, false, 0)

	// testview
	textview_file1 = gtk.NewTextView()
	textview_file2 = gtk.NewTextView()
	textview_file3 = gtk.NewTextView()

	textview_file1.SetCanFocus(false)
	textview_file2.SetCanFocus(false)
	textview_file3.SetCanFocus(false)

	textview_file1.ModifyFontEasy("DejaVu Serif 12")
	textview_file2.ModifyFontEasy("DejaVu Serif 12")
	textview_file3.ModifyFontEasy("DejaVu Serif 12")

	framebox1.PackStart(textview_file1, false, false, 0)
	framebox1.PackStart(textview_file2, false, false, 0)
	framebox1.PackStart(textview_file3, false, false, 0)

	// GtkVSeparator
	vsep1 := gtk.NewVSeparator()
	vsep1.Activate()
	vsep1.SetSizeRequest(5, 5)
	framebox1.PackStart(vsep1, false, false, 0)

	// progressbar
	progressbarBox := gtk.NewFixed()
	progressbar = gtk.NewProgressBar()
	progressbar.SetSizeRequest(windos_w, 50)
	progressbar.SetOrientation(gtk.PROGRESS_LEFT_TO_RIGHT)
	progressbarBox.Add(progressbar)
	framebox1.Add(progressbarBox)
	//-------------------------------framebox1-----end--------------------------

	//-------------------------------framebox2----start---------------------------
	//--------------------------------------------------------
	// GtkTextView ScrolledWindow
	//--------------------------------------------------------
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetSizeRequest(500, 1000)
	swin.SetShadowType(gtk.SHADOW_ETCHED_OUT)
	textview_output = gtk.NewTextView()
	swin.Add(textview_output)
	framebox2.Add(swin)
	//-------------------------------framebox2------end-------------------------

	//-------------------------------register event -------------------------------
	btn_load.Clicked(func() {
		//fmt.Println("load clicked")
		//--------------------------------------------------------
		// GtkFileChooserDialog
		//--------------------------------------------------------
		filechooserdialog := gtk.NewFileChooserDialog(
			"Choose File...",
			btn_load.GetTopLevelAsWindow(),
			gtk.FILE_CHOOSER_ACTION_OPEN,
			gtk.STOCK_OK,
			gtk.RESPONSE_ACCEPT)

		filter := gtk.NewFileFilter()
		filter.AddPattern("*.txt")
		filechooserdialog.AddFilter(filter)
		filechooserdialog.Response(func() {
			parsefilepath(textview_file1, textview_file2, textview_file3, filechooserdialog.GetFilename())
			filechooserdialog.Destroy()
		})
		filechooserdialog.Run()
	})

	btn_start.Clicked(func() {
		fmt.Println("step 9")
		textview_output.ModifyFontEasy("DejaVu Serif 12")
		textview_output.GetBuffer().SetText(" ")
		fmt.Println("step 10")
		go stlink_download()
		buttons.SetSensitive(false)
		fmt.Println("step 11")
	})

	combobox.Connect("changed", func() {
		fmt.Println("combobox changed value:", combobox.GetActiveText())
		if combobox.GetActive() == 1 {
			btn_load.SetSensitive(false)
			textview_output.GetBuffer().SetText("33333..")
		} else {
			btn_load.SetSensitive(true)
			textview_output.GetBuffer().SetText("1111..")
		}

	})
	//-------------------------------register event-------------------------------

	frame1.Add(framebox1)
	vpaned.Pack1(frame1, false, false)

	frame2.Add(framebox2)
	vpaned.Pack2(frame2, false, false)

	vbox.Add(vpaned)
	window.Add(vbox)

	window.SetSizeRequest(windos_w, windos_h)
	window.ShowAll()
	gtk.Main()
}

func is_ftm_download() bool {
	var i int = 0
	for {
		if Binfile[i].name == "" {
			return false
		}
		if strings.Contains(Binfile[i].name, "ftm") {
			return true
		}
		i = i + 1
	}
	return false
}

func stlink_download() {
	var ret int = 0
	var cmd1, cmd2, cmd3 string

	if is_ftm_download() {
		cmd1 = "b"
		cmd2 = "ftm-m"
		cmd3 = "ftm-f"
	} else {
		cmd1 = "b"
		cmd2 = "m"
		cmd3 = "f"
	}

	fmt.Println("step 1")
	progressbar.SetFraction(0)

	exec_stlink_download("reset-stlink")
	progressbar.SetText("Programming bootloader..")
	fmt.Println("step 2")
	ret = exec_stlink_download(cmd1)
	if ret == 0 {
		progressbar.SetFraction(0.33)
	} else {
		goto exit
	}
	fmt.Println("step 3")
	progressbar.SetText("Programming metadata..")
	ret = exec_stlink_download(cmd2)
	if ret == 0 {
		progressbar.SetFraction(0.7)
	} else {
		goto exit
	}
	fmt.Println("step 4")
	progressbar.SetText("Programming firmware..")
	ret = exec_stlink_download(cmd3)
	if ret == 0 {
		progressbar.SetFraction(1)
	} else {
		goto exit
	}

	progressbar.SetText("success.")
	buttons.SetSensitive(true)
	fmt.Println("step 5")
	return

exit:
	progressbar.SetText("FAILED!")
	buttons.SetSensitive(true)
}

func exec_shell_sh(action string) string {
	var script string
	var path string
	fmt.Println("step 6")
	if runtime.GOOS == "linux" {
		script = "/" + "download.sh"
		path = strings.Replace(mBinCfgFilePath, "/map.txt", "", -1)
	} else {
		script = "\\" + "download.bat"
		path = strings.Replace(mBinCfgFilePath, "\\map.txt", "", -1)
	}

	cmd := exec.Command(currentPath+script, path, action)
	fmt.Println("cmd:", cmd.Args)
	stdout, err := cmd.StdoutPipe()
	cmd.StdoutPipe()
	cmd.Start()
	content, err := ioutil.ReadAll(stdout)
	fmt.Println("step 7")
	if err != nil {
		fmt.Println(err)
		return "err"
	}
	fmt.Println("step 8")
	return string(content)
}

func exec_stlink_download(action string) int {
	log := exec_shell_sh(action)
	fmt.Println(log)
	var start gtk.TextIter
	buffer := textview_output.GetBuffer()
	buffer.GetStartIter(&start)
	buffer.Insert(&start, log)

	var ret bool
	if log != "" {
		ret = strings.Contains(log, "Flash written and verified! jolly good!")
		if ret == true {
			return 0
		}
		ret = strings.Contains(log, "Couldn't find any ST-Link")
		if ret == true {
			textview_output.ModifyFontEasy("DejaVu Serif 20")
			buffer.Insert(&start, "Couldn't find any ST-Link")
			return -2
		}
	}

	return -1
}

func parsefilepath(f1 *gtk.TextView, f2 *gtk.TextView, f3 *gtk.TextView, f string) string {
	if f == "" {
		return ""
	}
	btn_start.SetSensitive(true)

	mBinCfgFilePath = f
	ReadLine(f)
	var path string
	if runtime.GOOS == "linux" {
		path = strings.Replace(f, "map.txt", "", -1)
	} else {
		path = strings.Replace(f, "map.txt", "", -1)
	}

	f1.GetBuffer().SetText(path + Binfile[0].name + " " + Binfile[0].to)
	f2.GetBuffer().SetText(path + Binfile[1].name + " " + Binfile[1].to)
	f3.GetBuffer().SetText(path + Binfile[2].name + " " + Binfile[2].to)
	return path
}

type Binfile_t struct {
	name string
	to   string
}

var mBinCfgFilePath string
var currentPath string

var Binfile [4]Binfile_t

func getCurrDir() string {
	d, e := filepath.Abs(filepath.Dir(os.Args[0]))
	if e != nil {
		fmt.Println("err")
	}
	return strings.Replace(d, "\\", "/", -1)
}

func substring(source string, start int, end int) string {
	var substring = ""
	var pos = 0
	for _, c := range source {
		if pos < start {
			pos++
			continue
		}
		if pos >= end {
			break
		}
		pos++
		substring += string(c)
	}
	return substring
}

func parseLine(l string, i int) {
	var s, s1 string
	idx := strings.Index(l, ":")
	s = substring(l, 0, idx)
	if idx > 0 {
		s1 = substring(l, idx+1, idx+1+10)
		Binfile[i].name = s
		Binfile[i].to = s1
	}
	fmt.Println(i, s, s1)
}

func ReadLine(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	buf := bufio.NewReader(f)
	var i int = 0

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if i < 3 {
			parseLine(line, i)
			i = i + 1
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}
