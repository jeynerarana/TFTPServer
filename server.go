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
	//reports if there is an error listening
	if err != nil {
		log.Fatal(err)
	}

	// Since this is a listener the program must runner forever
	for {
		//variable will be used for reading the response
		data := make([]byte, 512)
		//saving the size, address to variables for later use
		size, addr, err := conn.ReadFrom(data)
		//if sokmething goes wrong while reading the data
		if err != nil {
			log.Print("Error: " + err.Error())
		}
		fmt.Printf("Read %v bytes.\n", size)
		fmt.Printf("Data: %s\n", data[1])
		//This switch check the frame at position 1, where the opcode is
		// this can help decide wether we are doing RRQ or WRQ
		/*
		        2 bytes      string    1 byte     string   1 byte
            -------------------------------------------------------
           | Opcode(01/02) |  Filename  |   0  |    Mode    |   0  |
            -------------------------------------------------------

                       Figure 5-1: RRQ/WRQ packet
         */
		switch data[1] {
		case 1:
			//---------------------RRQ---------------------------
			//this frame array holds the response code for a Data Response
			/*
			         2 bytes     2 bytes      n bytes
                   -------------------------------------
                  | Opcode(03)|   Block #  |   Data     |
                   -------------------------------------

                        Figure 5-2: DATA packet
			*/
			//frame will stay static and would never change the default bloack number is set to 1
			//you can see by the last byte
			frame :=[]byte{0,3,0,1}
			//this will keep count of the block number so we can iterate it on every packet
			// count starts at 1 since we set our default block to 1
			var blk_count byte=1
			//this will be the array that we will use to send the data
			// will will encapsolate it with data or other neccesary things
			data_send :=[]byte{0,3,0,1}
			fmt.Printf("Read Request: \n")
			//lets read the file
			//reading the file that the client requested
			readFile, err := ioutil.ReadFile("files/"+string(data[2:size-10]))
			//file not found
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
				//appends the data is it isnt 512 bytes
				data_send = append(data_send, byte(runes[i]))

			}else{	//sends data if its 512 bytes
					size, err = conn.WriteTo(data_send, addr)
					if err != nil {
						log.Print("Error: " + err.Error())
					}
					fmt.Printf("Sent %v bytes.\n", size)
					//clean up the data after sending it
					data_send = nil
					//reset the data
					data_send =frame
					//increment count block
					blk_count++
					data_send[3] =blk_count
					//connect again
					size, addr, err = conn.ReadFrom(data)
					if err != nil {
						log.Print("Error: " + err.Error())
					}
				}
			}
			//we have to right one final time if the last frame is less than 512 bytes
			size, err = conn.WriteTo(data_send, addr)
			if err != nil {
				log.Print("Error: " + err.Error())
			}
			fmt.Printf("Sent %v bytes.\n", size)
		case 2:
			//---------------------WRQ---------------------------
			//this is where we will keep the recieved data
			//we need to use lsize to keep the original size of data 
			//since size will be overwritten when we collect more data
			lsize :=size
			//Acknowlegment Frame
			frame :=[]byte{0,4,0,0}
			//variable we will use to encapsolate data
			data_rcv:= make([]byte, 512)
			//this arrays will be used to filter out only the data from the fram
			var file_data []byte
			//counter
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
					//increments block counter
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
			//tells you where its saved
			fmt.Printf(string(data[2:lsize-10]) + "\nSaved in:\n files/\n")
			//creates a file with the name given in the frame from th beggining
			f, err := os.Create("files/"+string(data[2:lsize-10]))
			if err != nil{
				log.Print("Error: " + err.Error())
				}
			defer f.Close()
			//writes data we collect to the file
			f.WriteString(string(file_data))	
		default:
			fmt.Printf("File Completely sent")
		}
	}
}
