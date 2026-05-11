package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MaplesMcDepth/mcdepth-tools/internal/jwt"
	"github.com/MaplesMcDepth/mcdepth-tools/internal/uuid"
)

func usage() {
	fmt.Print(`mcdepth - Swiss Army Knife for developers

Usage: mcdepth <command> [args]

Commands:
  fmt           Format/pretty-print JSON
  b64           Base64 encode/decode
  jwt           JWT decode (no verification)
  pass          Generate password
  url           URL encode/decode
  time          Timestamp conversions
  uuid          Generate UUID v4
  hash          MD5/SHA256 of string or file
  qr            Generate ASCII QR code (text mode)
  ip            Show local IP addresses

Examples:
  mcdepth fmt '{"a":1}'
  mcdepth b64 encode "hello world"
  mcdepth jwt decode eyJhbGciOiJIUzI1NiIs...
  mcdepth pass -l 20 -s
  mcdepth url encode "hello world"
  mcdepth time now
  mcdepth hash sha256 file.txt
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "fmt":
		runFmt(args)
	case "b64":
		runBase64(args)
	case "jwt":
		runJWT(args)
	case "pass":
		runPass(args)
	case "url":
		runURL(args)
	case "time":
		runTime(args)
	case "uuid":
		runUUID()
	case "hash":
		runHash(args)
	case "ip":
		runIP()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		usage()
		os.Exit(1)
	}
}

func runFmt(args []string) {
	var input string
	if len(args) > 0 {
		input = strings.Join(args, " ")
	} else {
		// Read from stdin
		data, err := os.ReadFile("/dev/stdin")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
			os.Exit(1)
		}
		input = string(data)
	}

	var v interface{}
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		fmt.Fprintln(os.Stderr, "Invalid JSON:", err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}

func runBase64(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: mcdepth b64 <encode|decode> <text>")
		os.Exit(1)
	}

	action := args[0]
	text := args[1]

	switch action {
	case "encode", "e":
		fmt.Println(base64.StdEncoding.EncodeToString([]byte(text)))
	case "decode", "d":
		data, err := base64.StdEncoding.DecodeString(text)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Decode error:", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
	default:
		fmt.Println("Usage: mcdepth b64 <encode|decode> <text>")
	}
}

func runJWT(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: mcdepth jwt decode <token>")
		os.Exit(1)
	}

	action := args[0]
	token := args[1]

	if action != "decode" && action != "d" {
		fmt.Println("Usage: mcdepth jwt decode <token>")
		os.Exit(1)
	}

	header, payload, err := jwt.Decode(token)
	if err != nil {
		fmt.Fprintln(os.Stderr, "JWT decode error:", err)
		os.Exit(1)
	}

	fmt.Println("Header:")
	fmt.Println(prettyJSON(header))
	fmt.Println("\nPayload:")
	fmt.Println(prettyJSON(payload))
}

func prettyJSON(data []byte) string {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return string(data)
	}
	out, _ := json.MarshalIndent(v, "", "  ")
	return string(out)
}

func runPass(args []string) {
	length := 16
	useSpecial := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-l", "--length":
			if i+1 < len(args) {
				l, err := strconv.Atoi(args[i+1])
				if err == nil {
					length = l
				}
				i++
			}
		case "-s", "--special":
			useSpecial = true
		}
	}

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if useSpecial {
		charset += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	password := make([]byte, length)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		password[i] = charset[n.Int64()]
	}
	fmt.Println(string(password))
}

func runURL(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: mcdepth url <encode|decode> <text>")
		os.Exit(1)
	}

	action := args[0]
	text := args[1]

	switch action {
	case "encode", "e":
		fmt.Println(url.QueryEscape(text))
	case "decode", "d":
		decoded, err := url.QueryUnescape(text)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Decode error:", err)
			os.Exit(1)
		}
		fmt.Println(decoded)
	default:
		fmt.Println("Usage: mcdepth url <encode|decode> <text>")
	}
}

func runTime(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: mcdepth time <now|unix|parse> [value]")
		os.Exit(1)
	}

	action := args[0]
	switch action {
	case "now":
		now := time.Now()
		fmt.Printf("Local:    %s\n", now.Format(time.RFC3339))
		fmt.Printf("Unix:     %d\n", now.Unix())
		fmt.Printf("Unix ms:  %d\n", now.UnixMilli())
		fmt.Printf("UTC:      %s\n", now.UTC().Format(time.RFC3339))
	case "unix":
		if len(args) < 2 {
			fmt.Println("Usage: mcdepth time unix <timestamp>")
			os.Exit(1)
		}
		ts, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid timestamp:", err)
			os.Exit(1)
		}
		t := time.Unix(ts, 0)
		fmt.Printf("Local:  %s\n", t.Format(time.RFC3339))
		fmt.Printf("UTC:    %s\n", t.UTC().Format(time.RFC3339))
	case "parse":
		if len(args) < 2 {
			fmt.Println("Usage: mcdepth time parse <RFC3339 string>")
			os.Exit(1)
		}
		t, err := time.Parse(time.RFC3339, args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Parse error:", err)
			os.Exit(1)
		}
		fmt.Printf("Unix:   %d\n", t.Unix())
		fmt.Printf("Local:  %s\n", t.Local().Format(time.RFC3339))
		fmt.Printf("UTC:    %s\n", t.UTC().Format(time.RFC3339))
	default:
		fmt.Println("Usage: mcdepth time <now|unix|parse> [value]")
	}
}

func runUUID() {
	u, err := uuid.GenerateV4()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Println(u)
}

func runHash(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: mcdepth hash <md5|sha256> <file|string>")
		os.Exit(1)
	}

	algo := args[0]
	target := args[1]

	var data []byte
	if _, err := os.Stat(target); err == nil {
		data, err = os.ReadFile(target)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Read error:", err)
			os.Exit(1)
		}
	} else {
		data = []byte(target)
	}

	switch algo {
	case "md5":
		fmt.Println(hex.EncodeToString(md5Sum(data)))
	case "sha256":
		fmt.Println(hex.EncodeToString(sha256Sum(data)))
	default:
		fmt.Println("Supported: md5, sha256")
	}
}

func md5Sum(data []byte) []byte {
	sum := md5.Sum(data)
	return sum[:]
}

func sha256Sum(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

func runIP() {
	fmt.Println("Network interfaces:")
	fmt.Println("  (Use 'ip addr' or 'ifconfig' for full details)")
}
