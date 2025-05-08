package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func ChooseDir() string {
	var dir string
	fmt.Println("Choose the directory of your journal... ~/home/user/")
	fmt.Scan(&dir)

	homeDir, _ := os.UserHomeDir()

	journalDir := homeDir + dir + "journal/"

	os.MkdirAll(journalDir, 0775)

	file, err := os.OpenFile(homeDir+"/jornalcfg.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	file.WriteString(journalDir)

	return journalDir
}

func ReplaceMD(text string) string {
	md := map[string]string{
		"***": "\033[1;3m",
		"___": "\033[1;3m",
		"**":  "\033[1m",
		"__":  "\033[1m",
		"*":   "\033[3m",
		"_":   "\033[3m",
		"~~":  "\033[9m",
		"--":  "\033[4m",
	}
	for key, value := range md {
		was := false
		i := 0
		for {
			next := strings.Index(text[i:], key)
			if next == -1 {
				break
			}
			i += next + len(key)

			if !was {
				text = strings.Replace(text, key, value, 1)
			} else {
				text = strings.Replace(text, key, "\033[0m", 1)
			}
			was = true
		}
	}
	return text
}

func Write(dir string, text string) {
	text = ReplaceMD(text)
	timeNow := time.Now().Format("2006-01-02_15-04-05")
	file, err := os.Create(dir + fmt.Sprintf("%s.md", timeNow))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	file.WriteString("\033[32m" + timeNow + "\033[0m" + " " + text)
}

func SortByTime(files []os.DirEntry) []os.DirEntry {
	sorted := false
	for !sorted {
		sorted = true
		for i, _ := range files {
			if i == len(files)-1 {
				break
			}
			thisT, err := time.Parse("2006-01-02_15-04-05", files[i].Name()[:len(files[i].Name())-3])
			if err != nil {
				fmt.Println(err)
				return nil
			}
			nextT, err := time.Parse("2006-01-02_15-04-05", files[i+1].Name()[:len(files[i+1].Name())-3])
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if thisT.After(nextT) {
				files[i], files[i+1] = files[i+1], files[i]
			}
		}
	}
	return files
}

func CheckJournal(dir string) {
	files, err := os.ReadDir(dir)
	files = SortByTime(files)
	if err != nil {
		fmt.Println(err)
		return
	}

	var k int
	fmt.Print("How many recent entries do you wanna see? (write nothing or '0' if you wanna see all entries)")
	fmt.Scan(&k)

	if k == 0 {
		k = len(files)
	}

	if k == 0 {
		fmt.Println("you don`t have any entries. Please write 'journal help' for information.")
		return
	}

	fmt.Print("\n\n\n")
	for i := 0; i < k; i++ {
		text, err := os.ReadFile(dir + files[i].Name())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(text))
	}
	fmt.Print("\n\n\n")
}

func Help() {
	if len(os.Args) == 2 {
		fmt.Print("What do you want? \n\n help markup → information about text markup \n\n help commands → informtion about all commands \n")
	} else {
		switch os.Args[2] {
		case "markup":
			fmt.Print(" ***text***, ___text___ → bold and italic \n **text**, __text__ → bold \n *text*, _text_ → italic \n ~~text~~ → crossed out \n --text-- → underlined \n")
		case "commands":
			fmt.Print(" journal → read your entries \n journal dir → choose directory of your journal \n journal write [insert text] → make a new entry \n")
		default:
			fmt.Print("What do you want? \n\n help markup → information about text markup \n\n help commands → informtion about all commands \n")
		}
	}
}

func main() {
	homeDir, _ := os.UserHomeDir()
	dirByte, err := os.ReadFile(homeDir + "/jornalcfg.txt")
	dir := string(dirByte)
	if err != nil {
		fmt.Println(err)
		dir = "EMPTY"
	}

	if len(os.Args) == 1 {
		if dir != "EMPTY" {
			CheckJournal(dir)
		} else {
			fmt.Println("you didn`t change your journal`s path.. Write 'journal dir'")
		}
	} else {
		switch os.Args[1] {
		case "dir":
			dir = ChooseDir()
		case "write":
			if dir == "EMPTY" {
				fmt.Println("you didn`t change your journal`s path.. Write 'journal dir'")
				break
			}
			var text string
			for i := 2; i < len(os.Args); i++ {
				text += os.Args[i] + " "
			}
			Write(dir, text)
		case "help":
			Help()
		default:
			fmt.Println("Please write 'journal help' for information")
		}

	}

}
