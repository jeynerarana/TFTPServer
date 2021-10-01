package main

import (
	"fmt"
	"log"
	"net"
	"io/ioutil"
	"os"
)
func banner(){
	fmt.Println("-----------Welcome to Jeyner's TFTP Server-----------")
	for i:=0; i < 5; i++{
		fmt.Println("-----------------------------------------------------")
	}

}

func main() {
	//print my banner
	banner()

	// socket is listening on port 69
	conn, err := net.ListenPacket("udp", ":69")
	if err != nil {
		log.Fatal(err)
	}

	// Since this is a listener the program must runner forever
	for {
		data := make([]byte, 512)
		// blocking call, ReadFrom conn, store result in data
		// addr is who sent us the data
		// size of data is also returned from the function
		size, addr, err := conn.ReadFrom(data)
		if err != nil {
			log.Print("Error: " + err.Error())
		}
		fmt.Printf("Read %v bytes.\n", size)
		fmt.Printf("Data: %s\n", data[1])
		//seeing the opcode
		switch data[1] {
		case 1:
			//---------------------RRQ---------------------------
			//var frame[]byte
			frame :=[]byte{0,3,0,1}
			var blk_count byte=1
			data_send :=[]byte{0,3,0,1}
			fmt.Printf("Read Request: \n")
			//lets read the file
			readFile, err := ioutil.ReadFile("files/"+string(data[2:size-10]))
			
			if err != nil {
	    		fmt.Println("error file not Found")
			}
			sFile := string(readFile)
			//prints out an error if there is nothing 
			//sending files
			runes:= []rune(sFile)
			//This part seperates the bytes into 512 chunks and sends them to the client
			for i := 0; i < len(runes); i++ {
				if (i%513 !=0 || i==0){
				data_send = append(data_send, byte(runes[i]))

			}else{
					size, err = conn.WriteTo(data_send, addr)
					if err != nil {
						log.Print("Error: " + err.Error())
					}
					fmt.Printf("Sent %v bytes.\n", size)
					data_send = nil
					//addr =nil
					data_send =frame
					blk_count++
					data_send[3] =blk_count
					//connect again
					size, addr, err = conn.ReadFrom(data)
					if err != nil {
						log.Print("Error: " + err.Error())
					}
				}
			}

			size, err = conn.WriteTo(data_send, addr)
			if err != nil {
				log.Print("Error: " + err.Error())
			}
			fmt.Printf("Sent %v bytes.\n", size)
		case 2:
			//---------------------WRQ---------------------------
			//this is where we will keeo the recieved data
			lsize :=size
			frame :=[]byte{0,4,0,0}
			data_rcv:= make([]byte, 512)
			var file_data []byte
			var blk_count byte=1
			ack_packet := []byte{0,4,0,1}//it would be better if the last byte is just grabbed from sender
			fmt.Printf("Write Request")
			//first Acknowlgemnet packet
			size, err = conn.WriteTo(frame, addr)
			if err != nil {
				log.Print("Error: " + err.Error())
			}
			//increase counter
			//no we start reading
				for{
					size1, addr, err := conn.ReadFrom(data_rcv)
					if err != nil {
						log.Print("Error: " + err.Error())
						}
					//save the data
					file_data = append(file_data,data_rcv[4:size1]...)
					size, err = conn.WriteTo(ack_packet, addr)
					if err != nil {
						log.Print("Error: " + err.Error())
						}
					ack_packet =frame
					blk_count++
					ack_packet[3] =blk_count
					//check for terminations
					if(size1 <512 && size1 > 0){
						break
					}

				}
			//now lets put tha file
			//print it
			fmt.Println("\nSaving....\n")
			fmt.Printf(string(data[2:lsize-10]) + "\nSaved in:\n files/\n")
			f, err := os.Create("files/"+string(data[2:lsize-10]))
			if err != nil{
				log.Print("Error: " + err.Error())
				}
			defer f.Close()
			f.WriteString(string(file_data))	
		default:
			fmt.Printf("File Completely sent")
		}
	}
}
