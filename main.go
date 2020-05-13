package main

import "fmt"
import "net"
import "bytes"
import "encoding/binary"

//                                    1  1  1  1  1  1
//      0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                      ID                       |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                    QDCOUNT                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                    ANCOUNT                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                    NSCOUNT                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                    ARCOUNT                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type DNSHeader struct {
	ID int16 // random dns id
	QR bool // question == 0, answer == 1
	Opcode [4]bool // type of question
	        // 0 == a standard query (QUERY)
		// 1 == an inverse query (IQUERY)
		// 2 == a server status request (STATUS)
		// 3-15 == reserved for future use
	AA bool // privilege. authority answer == 1, other ==0
	TC bool // TrunCation notify bit
	RD bool // Recursion query bit
	RA bool // Recursion available bit
	Z [3]bool  // Reserved for future use.  Must be zero.
	RCODE [4]bool // Respon code.
	              // 0 == No error condition
		      // 1 == Format error
		      // 2 == Server failure
		      // 3 == Name Error(Name not found from authority)
		      // 4 == Not Implemented
		      // 5 == Refused
		      // 6-15 == Reserved for future use
	QDCOUNT int16
	ANCOUNT int16
	NSCOUNT int16
	ARCOUNT int16
}

type DNSQuestion struct {
	QNAME []byte
	QTYPE int16
	QCLASS int16
}

type DNSResource struct {
	NAME []byte
	TYPE int16
	CLASS int16
	TTL int32
	RDLENGTH int16
	RDATA []byte
}

type DNSPacket struct {
	Header DNSHeader
	Question DNSQuestion
	Answer []DNSResource
	Authority []DNSResource
	Additional []DNSResource
}

func main() {
	fmt.Println("Hello, onokatio full DNS resolver.")
	conn, _ := net.Dial("udp", "198.41.0.4:53") // 198.41.0.4 == a.root-servers.net
	defer conn.Close()
	//var query = []byte{0xa9, 0x4e, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x77, 0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x29, 0x05, 0xc8, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00}
	//question := DNSQuestion{
	//	QNAME: "",
	//}
	/*
	header := DNSHeader{
		ID: 0x0f,
		QR: false,
		Opcode: [4]bool{false,false,false,false},
		AA: false,
		TC: false,
		RD: false,
		RA: false,
		Z: [3]bool{false,false,false},
		RCODE: [4]bool{false,false,false,false},
		QDCOUNT: 1,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}
	*/

	rawpacket := new(bytes.Buffer)

	binary.Write(rawpacket, binary.LittleEndian, int16(0x000f))
	binary.Write(rawpacket, binary.LittleEndian, int8(0x00 | 0x00<<1 | 0x00<<2 | 0x00<<3 | 0x00<<4 | 0x00<<5 | 0x00<<6 | 0x00<<7)) // ID + Opcode + RD
	binary.Write(rawpacket, binary.LittleEndian, int8(0x00 | 0x00<<1 | 0x00<<2 | 0x00<<3 | 0x00<<4 | 0x00<<5 | 0x00<<6 | 0x00<<7)) // RA + Z + RCODE
	binary.Write(rawpacket, binary.LittleEndian, int16(0x0001))
	binary.Write(rawpacket, binary.LittleEndian, int16(0x0000))
	binary.Write(rawpacket, binary.LittleEndian, int16(0x0000))
	binary.Write(rawpacket, binary.LittleEndian, int16(0x0000))
	err := binary.Write(rawpacket, binary.LittleEndian, []byte{0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00} )
	if err != nil {
		panic(err)
	}

	fmt.Printf("%x", rawpacket.Bytes())
	conn.Write(rawpacket.Bytes())

	response := make([]byte, 2000)
	conn.Read(response)
	fmt.Println(response)
}
