/*
 * JA3 - TLS Client Hello Hash
 * Copyright (c) 2017, Salesforce.com, Inc.
 * this code was created by Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package ja3

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/gopacket/pcapgo"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// PacketSource means we can read Packets
type PacketSource interface {
	ReadPacketData() ([]byte, gopacket.CaptureInfo, error)
}

// ReadFileCSV reads the PCAP file at the given path
// and prints out all packets containing JA3 digests to the supplied io.Writer
// currently no PCAPNG support
func ReadFileCSV(file string, out io.Writer, separator string) {

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var r PacketSource

	r, err = pcapgo.NewReader(f)
	if err != nil {
		// maybe its a PCAPNG
		r, err = pcapgo.NewNgReader(f, pcapgo.DefaultNgReaderOptions)
		if err != nil {
			// nope
			panic(err)
		}
	}

	columns := []string{"timestamp", "source_ip", "source_port", "destination_ip", "destination_port", "ja3_digest", "\n"}
	out.Write([]byte(strings.Join(columns, separator)))

	for {
		// read packet data
		data, ci, err := r.ReadPacketData()
		if err == io.EOF {
			return
		} else if err != nil {
			panic(err)
		}

		var (
			// create gopacket
			p = gopacket.NewPacket(data, layers.LinkTypeEthernet, gopacket.Lazy)
			// get JA3 if possible
			digest = DigestHexPacket(p)
		)

		// check if we got a result
		if digest != "" {

			var (
				b  strings.Builder
				nl = p.NetworkLayer()
				tl = p.TransportLayer()
			)

			if tl == nil || nl == nil {
				fmt.Println("error: ", nl, tl, p.Dump())
				continue
			}

			b.WriteString(timeToString(ci.Timestamp))
			b.WriteString(separator)
			b.WriteString(nl.NetworkFlow().Src().String())
			b.WriteString(separator)
			b.WriteString(tl.TransportFlow().Src().String())
			b.WriteString(separator)
			b.WriteString(nl.NetworkFlow().Dst().String())
			b.WriteString(separator)
			b.WriteString(tl.TransportFlow().Dst().String())
			b.WriteString(separator)
			b.WriteString(digest)
			b.WriteString("\n")

			_, err := out.Write([]byte(b.String()))
			if err != nil {
				panic(err)
			}
		}
	}
}
