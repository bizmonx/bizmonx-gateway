package xymon

import (
	"fmt"
	"net"
)

type StatusMessage struct {
	Message  string `json:"message"`
	Color    string `json:"color"`
	Test     string `json:"test"`
	Host     string `json:"host"`
	Lifetime string `json:"lifetime"`
}

type DataMessage struct {
	Message  string `json:"message"`
	Host     string `json:"host"`
	DataName string `json:"data_name"`
}

type XymonHost struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (x *StatusMessage) Stringify() string {
	lifetime := map[bool]string{true: "+" + x.Lifetime, false: ""}[x.Lifetime != ""]
	m := fmt.Sprintf("status%s %s.%s %s %s ", lifetime, x.Host, x.Test, x.Color, x.Message)
	fmt.Println(m)
	return m
}

func (x *DataMessage) Stringify() string {
	m := fmt.Sprintf("data %s.%s\n%s ", x.Host, x.DataName, x.Message)
	fmt.Println(m)
	return m
}

func (x *DataMessage) Send(host *XymonHost) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host.Host, host.Port))
	if err != nil {
		fmt.Printf("Error connecting to %s:%d: %s\n", host.Host, host.Port, err.Error())
		return err
	}
	defer conn.Close()

	// Send the message
	_, err = conn.Write([]byte(x.Stringify()))
	if err != nil {
		fmt.Println("Error sending message:", err.Error())
		return err
	}

	fmt.Println("Message successfully sent to Xymon.")
	return nil
}

func (x *StatusMessage) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.Message + `"`), nil
}

func (x *StatusMessage) Send(h *XymonHost) error {
	conn, err := net.Dial("tcp", "localhost:1984")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		//os.Exit(1)
		return err
	}
	// Close the connection when the function exits
	defer conn.Close()

	// Send the message
	_, err = conn.Write([]byte(x.Stringify()))
	if err != nil {
		// Handle error
		fmt.Println("Error sending message:", err.Error())
		return err
	}

	fmt.Println("Message successfully sent to Xymon.")
	return nil
}
