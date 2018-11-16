/*
 * JA3 - TLS Client Hello Hash
 * Copyright (c) 2018 Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dreadl0ck/ja3"
)

var (
	flagJSON      = flag.Bool("json", true, "print as JSON array")
	flagCSV       = flag.Bool("csv", false, "print as CSV")
	flagTSV       = flag.Bool("tsv", false, "print as TAB separated values")
	flagSeparator = flag.String("separator", ",", "set a custom separator")
	flagInput     = flag.String("read", "", "read PCAP file")
)

func main() {

	flag.Parse()

	if *flagInput == "" {
		fmt.Println("use the -read flag to supply an input file.")
		os.Exit(1)
	}

	if *flagTSV {
		ja3.ReadFileCSV(*flagInput, os.Stdout, "\t")
		return
	}

	if *flagCSV {
		ja3.ReadFileCSV(*flagInput, os.Stdout, *flagSeparator)
		return
	}

	if *flagJSON {
		ja3.ReadFileJSON(*flagInput, os.Stdout)
	}
}
