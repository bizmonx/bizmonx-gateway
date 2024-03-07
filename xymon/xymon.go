package xymon

import (
	"fmt"
	"net"
	"os"
)

type XymonMessage struct {
	Message  string `json:"message"`
	Color    string `json:"color"`
	Test     string `json:"test"`
	Host     string `json:"host"`
	Lifetime string `json:"lifetime"`
}

type XymonHost struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (x *XymonMessage) Stringify() string {
	lifetime := map[bool]string{true: "+" + x.Lifetime, false: ""}[x.Lifetime != ""]
	m := fmt.Sprintf("status%s %s.%s %s %s ", lifetime, x.Host, x.Test, x.Color, x.Message)
	fmt.Print(m)
	return m
}

func (x *XymonMessage) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.Message + `"`), nil
}

func (x *XymonMessage) Send(h *XymonHost) error {
	conn, err := net.Dial("tcp", "localhost:1984")
	if err != nil {
		// Handle error
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	// Close the connection when the function exits
	defer conn.Close()

	// Message to send
	///message := "status newhost.go green message from go client\n"

	// Send the message
	_, err = conn.Write([]byte(x.Stringify()))
	if err != nil {
		// Handle error
		fmt.Println("Error sending message:", err.Error())
		return err
	}

	fmt.Println("Message sent successfully")
	return nil
}
