package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var version = "dev"

var (
	algo   = flag.String("a", "sha256", "hash algorithm: sha256 | md5")
	width  = flag.Int("w", 17, "grid width")
	height = flag.Int("h", 9, "grid height")
	color  = flag.Bool("c", false, "enable ANSI color output")
	key    = flag.String("k", "", "read input from file (e.g. SSH public key)")
)

var ramp = []rune(" .o+=*BOX@#")

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func readInput() (string, error) {
	if *key != "" {
		data, err := os.ReadFile(*key)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	}

	if flag.NArg() > 0 {
		return strings.Join(flag.Args(), " "), nil
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeCharDevice == 0 {
		reader := bufio.NewReader(os.Stdin)
		data, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	}

	return "", fmt.Errorf("no input provided")
}

func computeHash(input string) []byte {
	switch strings.ToLower(*algo) {
	case "md5":
		sum := md5.Sum([]byte(input))
		return sum[:]
	default:
		sum := sha256.Sum256([]byte(input))
		return sum[:]
	}
}

func colorize(v int, s string) string {
	if !*color {
		return s
	}
	code := 30 + (v % 8)
	return fmt.Sprintf("\033[%dm%s\033[0m", code, s)
}

func randomArt(input string) {
	hash := computeHash(input)

	grid := make([][]int, *height)
	for i := range grid {
		grid[i] = make([]int, *width)
	}

	x, y := *width/2, *height/2
	startX, startY := x, y

	for _, b := range hash {
		for range 4 {
			move := b & 0x03
			b >>= 2

			switch move {
			case 0:
				x--
				y--
			case 1:
				x++
				y--
			case 2:
				x--
				y++
			case 3:
				x++
				y++
			}

			x = clamp(x, 0, *width-1)
			y = clamp(y, 0, *height-1)
			grid[y][x]++
		}
	}

	endX, endY := x, y

	hashHex := hex.EncodeToString(hash)
	fmt.Printf("+[%s] [%s]+\n", strings.ToUpper(*algo), hashHex[:16])

	fmt.Print("+")
	for i := 0; i < *width; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	for y := 0; y < *height; y++ {
		fmt.Print("|")
		for x := 0; x < *width; x++ {
			switch {
			case x == startX && y == startY:
				fmt.Print(colorize(9, "S"))
			case x == endX && y == endY:
				fmt.Print(colorize(9, "E"))
			default:
				v := grid[y][x]
				if v >= len(ramp) {
					fmt.Print(colorize(v, string(ramp[len(ramp)-1])))
				} else {
					fmt.Print(colorize(v, string(ramp[v])))
				}
			}
		}
		fmt.Println("|")
	}

	fmt.Print("+")
	for i := 0; i < *width; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(version)
		return
	}

	flag.Parse()

	input, err := readInput()
	if err != nil {
		fmt.Fprintln(os.Stderr, "randomart:", err)
		fmt.Fprintln(os.Stderr, "usage: randomart [options] <text> | echo <text> | randomart")
		flag.PrintDefaults()
		os.Exit(1)
	}

	randomArt(input)
}
