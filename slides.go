package main

import (
	"bufio"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

var (
	lines   []string
	linum   = 1
	newpage = true
)

func main() {
	in := os.Stdin
	if len(os.Args) > 1 {
		var err error
		in, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Println(header)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " \n\t\v")
		switch {
		case line == "":
			printPage()
		case strings.HasPrefix(line, "\\"):
			// escaped lines are inserted verbatim
			lines = append(lines, line[1:])
		case strings.HasPrefix(line, "@"):
			opts["image"].Set(line[1:])
		case strings.HasPrefix(line, "#+"):
			// "command" comments are parsed seperatly
			sub := cmdRe.FindStringSubmatch(line)
			opt, ok := opts[strings.ToLower(sub[1])]
			if !ok {
				fmt.Fprintf(os.Stderr, "unknown command %q in line %d\n",
					sub[1], linum)
				break
			}
			opt.Set(sub[3])
			if sub[2] != "" {
				opt.Push()
			}
		case strings.HasPrefix(line, "#"):
			// regular comments are ignored
		default:
			lines = append(lines, line)
		}
		linum++
	}
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, scanner.Err().Error())
	}
	printPage()
}
