package main

import (
	"log"
	"fmt"
	"net"
	"io/ioutil"
	"strings"
	"bytes"
	"bufio"
	"net/http"
	"net/http/httputil"
)

func HandleUpnpMessage(upnplisten *net.UDPConn, SourceAddr *net.UDPAddr, UpnpRequest *http.Request){
	log.Println("Upnp from: ", SourceAddr)
	//log.Println("Mesage: ", string(Message))
	//Sc := httputil.NewServerConn(rw, nil)
	resp := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		ContentLength: int64(0),
		Body:          ioutil.NopCloser(strings.NewReader("")),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
	}
	resp.Header.Set("CACHE-CONTROL", "max-age=120")
	resp.Header.Set("ST", "urn:schemas-upnp-org:device:InternetGatewayDevice:1")
	resp.Header.Set("USN", "uuid:3d3cec3a-8cf0-11e0-98ee-001a6bd2d07b::urn:schemas-upnp-org:device:InternetGatewayDevice:1")
	resp.Header.Set("EXT", "")
	resp.Header.Set("SERVER", "Ubuntu/precise UPnP/1.1 GoUPNP/0.1")
	resp.Header.Set("LOCATION", "http://"+GetFirstIpOfETH("eth1")+":55455/rootDesc.xml")
	resp.Header.Set("OPT", "\"http://schemas.upnp.org/upnp/1/0/\"; ns=01")
	resp.Header.Set("01-NLS", "1")
	resp.Header.Set("BOOTID.UPNP.ORG", "1")
	resp.Header.Set("CONFIGID.UPNP.ORG", "1337")

	UpnpResponseData, _ := httputil.DumpResponse(resp, true)
	upnplisten.WriteTo(UpnpResponseData, SourceAddr)

	//log.Println("Responded to Upnp Message")
	//log.Printf("SourceAddr: %+v", SourceAddr)
	//log.Printf("UpnpRequest: %+v", UpnpRequest)
	//log.Printf("UpnpResponse: %+v", resp)
	


}

func rootDesc_xml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml; charset=\"utf-8\"")
	w.Header().Set("Connection", "close")
	//fmt.Fprintf(w, "This is rootDesc_xml it descripes my functions")
	fmt.Fprintf(w ,"<?xml version=\"1.0\"?>\n<root xmlns=\"urn:schemas-upnp-org:device-1-0\"><specVersion><major>1</major><minor>0</minor></specVersion><device><deviceType>urn:schemas-upnp-org:device:InternetGatewayDevice:1</deviceType><friendlyName>Ubuntu router</friendlyName><manufacturer>Ubuntu</manufacturer><manufacturerURL>http://www.ubuntu.com/</manufacturerURL><modelDescription>Ubuntu router</modelDescription><modelName>Ubuntu router</modelName><modelNumber>1</modelNumber><modelURL>http://www.ubuntu.com/</modelURL><serialNumber>12345678</serialNumber><UDN>uuid:3d3cec3a-8cf0-11e0-98ee-001a6bd2d07b</UDN><serviceList><service><serviceType>urn:schemas-upnp-org:service:Layer3Forwarding:1</serviceType><serviceId>urn:upnp-org:serviceId:Layer3Forwarding1</serviceId><controlURL>/ctl/L3F</controlURL><eventSubURL>/evt/L3F</eventSubURL><SCPDURL>/L3F.xml</SCPDURL></service></serviceList><deviceList><device><deviceType>urn:schemas-upnp-org:device:WANDevice:1</deviceType><friendlyName>WANDevice</friendlyName><manufacturer>MiniUPnP</manufacturer><manufacturerURL>https://github.com/callesg/upnp_portforwader</manufacturerURL><modelDescription>WAN Device</modelDescription><modelName>WAN Device</modelName><modelNumber>20140315</modelNumber><modelURL>https://github.com/callesg/upnp_portforwader</modelURL><serialNumber>12345678</serialNumber><UDN>uuid:3d3cec3a-8cf0-11e0-98ee-001a6bd2d07c</UDN><UPC>000000000000</UPC><serviceList><service><serviceType>urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1</serviceType><serviceId>urn:upnp-org:serviceId:WANCommonIFC1</serviceId><controlURL>/ctl/CmnIfCfg</controlURL><eventSubURL>/evt/CmnIfCfg</eventSubURL><SCPDURL>/WANCfg.xml</SCPDURL></service></serviceList><deviceList><device><deviceType>urn:schemas-upnp-org:device:WANConnectionDevice:1</deviceType><friendlyName>WANConnectionDevice</friendlyName><manufacturer>MiniUPnP</manufacturer><manufacturerURL>https://github.com/callesg/upnp_portforwader</manufacturerURL><modelDescription>MiniUPnP daemon</modelDescription><modelName>MiniUPnPd</modelName><modelNumber>20140315</modelNumber><modelURL>https://github.com/callesg/upnp_portforwader</modelURL><serialNumber>12345678</serialNumber><UDN>uuid:3d3cec3a-8cf0-11e0-98ee-001a6bd2d07d</UDN><UPC>000000000000</UPC><serviceList><service><serviceType>urn:schemas-upnp-org:service:WANIPConnection:1</serviceType><serviceId>urn:upnp-org:serviceId:WANIPConn1</serviceId><controlURL>/ctl/IPConn</controlURL><eventSubURL>/evt/IPConn</eventSubURL><SCPDURL>/WANIPCn.xml</SCPDURL></service></serviceList></device></deviceList></device></deviceList><presentationURL>http://"+GetFirstIpOfETH("eth1")+"/</presentationURL></device></root>")
}

func GetFirstIpOfETH(eth string) (string){
	ExternalETH, err := net.InterfaceByName(eth)
	if err != nil {
		panic(err)
	}
	Addresses, err := ExternalETH.Addrs()
	if err != nil {
		panic(err)
	}
	if(len(Addresses) < 1){
		log.Printf("Less that one ip address on external interface",)
		return ""
	}
	IPParts := strings.Split(Addresses[0].String(), "/")
	return IPParts[0]
}
func DoSoapAction(Action string, xml string) (ReturnData string){
	xmlnsbs := "urn:schemas-upnp-org:service:WANIPConnection:1"

	//Pre XML
	ReturnData = "<?xml version=\"1.0\"?>"
	ReturnData += "<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\" s:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\"><s:Body>"

	if Action == "GetExternalIPAddress" {
		ExternalIp := GetFirstIpOfETH("eth0")
		log.Println("Reporing External Ip:", ExternalIp)

		ReturnData += "<u:"+Action+"Response xmlns:u=\""+xmlnsbs+"\">"
		ReturnData += "<NewExternalIPAddress>"+ExternalIp+"</NewExternalIPAddress>"
		ReturnData += "</u:"+Action+"Response>"

	}else if Action == "AddPortMapping" {
		ReturnData += "<u:"+Action+"Response xmlns:u=\""+xmlnsbs+"\"/>"
		log.Printf("XML: %s", xml)
		
	}else{
		log.Printf("Unknown Action: %s", Action)
		log.Printf("XML: %s", xml)
	}

	//Post XML
	ReturnData +="</s:Body></s:Envelope>"
	return ReturnData
}

func ctl_IPConn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml; charset=\"utf-8\"")
	w.Header().Set("Connection", "close")

	if val, ok := r.Header["Soapaction"]; ok {
		SoapAccPart := strings.Split(val[0], "#")
		if len(SoapAccPart) != 2 {
			fmt.Fprintf(w, "To many #")
			log.Printf("To many #")
		}else{
			SoapAction := strings.TrimRight(SoapAccPart[1],"\"")
			log.Printf("Client requested action: %s", SoapAction)
			PostData, err := ioutil.ReadAll(r.Body)
			if err == nil {
				fmt.Fprintf(w ,DoSoapAction(SoapAction, string(PostData)))
			}else{
				log.Printf("error reciving post data")
			}
		}
	}else{
		fmt.Fprintf(w, "No Soapaction")
		log.Printf("No Soapaction")
	}

	//fmt.Fprintf(w, "This is ctl_IPConn here i will take request for diffrent stuff")
	//fmt.Fprintf(w, "Sush as: urn:schemas-upnp-org:service:WANIPConnection:1#GetExternalIPAddress")
	//fmt.Fprintf(w, "Sush as: urn:schemas-upnp-org:service:WANIPConnection:1#AddPortMapping")

}


func main() {

	addr, _ := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	ComputerInterface, err := net.InterfaceByName("eth1")
	upnplisten, err := net.ListenMulticastUDP("udp4", ComputerInterface, addr) 
	if err != nil {
		log.Println("Could not listen to upnp port udp 1900")
	}
	defer upnplisten.Close()


	http.HandleFunc("/rootDesc.xml", rootDesc_xml)
	http.HandleFunc("/ctl/IPConn", ctl_IPConn)

	go http.ListenAndServe(":55455", nil)

	for {
		upnpmes := make([]byte, 2048)
		readlen, SourceAddr, err := upnplisten.ReadFromUDP(upnpmes)
		if err != nil {
			log.Println("udp read fail: len:  %i, mes: %s", err,readlen)
		}else{
			if(readlen == len(upnpmes)){
				log.Println("filled the udp message buffer incoming message to log can not handle the message")
			}else{
				//log.Println("upnp data", string(upnpmes[:readlen]))
				UpnpHttpreader := httputil.NewServerConn(nil, bufio.NewReader(bytes.NewReader(upnpmes[:readlen])))
				defer UpnpHttpreader.Close()

				UpnpRequest, err := UpnpHttpreader.Read()
				if err != nil {
					continue
				}
				HandleUpnpMessage(upnplisten, SourceAddr, UpnpRequest)
			}
		}
	}
}
