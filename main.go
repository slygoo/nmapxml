package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const PORT_OPEN = "open"

type NmapXML struct {
	XMLName          xml.Name `xml:"nmaprun"`
	Text             string   `xml:",chardata"`
	Scanner          string   `xml:"scanner,attr"`
	Args             string   `xml:"args,attr"`
	Start            string   `xml:"start,attr"`
	Startstr         string   `xml:"startstr,attr"`
	Version          string   `xml:"version,attr"`
	Xmloutputversion string   `xml:"xmloutputversion,attr"`
	Scaninfo         struct {
		Text        string `xml:",chardata"`
		Type        string `xml:"type,attr"`
		Protocol    string `xml:"protocol,attr"`
		Numservices string `xml:"numservices,attr"`
		Services    string `xml:"services,attr"`
	} `xml:"scaninfo"`
	Verbose struct {
		Text  string `xml:",chardata"`
		Level string `xml:"level,attr"`
	} `xml:"verbose"`
	Debugging struct {
		Text  string `xml:",chardata"`
		Level string `xml:"level,attr"`
	} `xml:"debugging"`
	Hosthint []struct {
		Text   string `xml:",chardata"`
		Status struct {
			Text      string `xml:",chardata"`
			State     string `xml:"state,attr"`
			Reason    string `xml:"reason,attr"`
			ReasonTtl string `xml:"reason_ttl,attr"`
		} `xml:"status"`
		Address []struct {
			Text     string `xml:",chardata"`
			Addr     string `xml:"addr,attr"`
			Addrtype string `xml:"addrtype,attr"`
			Vendor   string `xml:"vendor,attr"`
		} `xml:"address"`
		Hostnames string `xml:"hostnames"`
	} `xml:"hosthint"`
	Host []struct {
		Text      string `xml:",chardata"`
		Starttime string `xml:"starttime,attr"`
		Endtime   string `xml:"endtime,attr"`
		Status    struct {
			Text      string `xml:",chardata"`
			State     string `xml:"state,attr"`
			Reason    string `xml:"reason,attr"`
			ReasonTtl string `xml:"reason_ttl,attr"`
		} `xml:"status"`
		Address []struct {
			Text     string `xml:",chardata"`
			Addr     string `xml:"addr,attr"`
			Addrtype string `xml:"addrtype,attr"`
			Vendor   string `xml:"vendor,attr"`
		} `xml:"address"`
		Hostnames string `xml:"hostnames"`
		Ports     struct {
			Text       string `xml:",chardata"`
			Extraports struct {
				Text         string `xml:",chardata"`
				State        string `xml:"state,attr"`
				Count        string `xml:"count,attr"`
				Extrareasons struct {
					Text   string `xml:",chardata"`
					Reason string `xml:"reason,attr"`
					Count  string `xml:"count,attr"`
					Proto  string `xml:"proto,attr"`
					Ports  string `xml:"ports,attr"`
				} `xml:"extrareasons"`
			} `xml:"extraports"`
			Port []struct {
				Text     string `xml:",chardata"`
				Protocol string `xml:"protocol,attr"`
				Portid   string `xml:"portid,attr"`
				State    struct {
					Text      string `xml:",chardata"`
					State     string `xml:"state,attr"`
					Reason    string `xml:"reason,attr"`
					ReasonTtl string `xml:"reason_ttl,attr"`
				} `xml:"state"`
				Service struct {
					Text   string `xml:",chardata"`
					Name   string `xml:"name,attr"`
					Method string `xml:"method,attr"`
					Conf   string `xml:"conf,attr"`
				} `xml:"service"`
			} `xml:"port"`
		} `xml:"ports"`
		Times struct {
			Text   string `xml:",chardata"`
			Srtt   string `xml:"srtt,attr"`
			Rttvar string `xml:"rttvar,attr"`
			To     string `xml:"to,attr"`
		} `xml:"times"`
	} `xml:"host"`
	Runstats struct {
		Text     string `xml:",chardata"`
		Finished struct {
			Text    string `xml:",chardata"`
			Time    string `xml:"time,attr"`
			Timestr string `xml:"timestr,attr"`
			Summary string `xml:"summary,attr"`
			Elapsed string `xml:"elapsed,attr"`
			Exit    string `xml:"exit,attr"`
		} `xml:"finished"`
		Hosts struct {
			Text  string `xml:",chardata"`
			Up    string `xml:"up,attr"`
			Down  string `xml:"down,attr"`
			Total string `xml:"total,attr"`
		} `xml:"hosts"`
	} `xml:"runstats"`
}

type HostAndPorts struct {
	Host  string
	Ports []string
}

var AllHostsAndPorts []HostAndPorts

func main() {
	f := flag.String("f", "", "filename, if no filename is given it will load all the xmls in the current directory")
	flag.Parse()
	filepath := *f
	AllNmaps := []NmapXML{}
	if filepath == "" {
		xmls, err := GetAllXMLsInCurrentDirectory()
		if err != nil {
			log.Println(err)
			return
		}
		for i := 0; i < len(xmls); i++ {
			NmapXML, err := LoadXMLFile(xmls[i])
			if err != nil {
				log.Println(err)
				continue
			}
			AllNmaps = append(AllNmaps, NmapXML)
		}
	} else {
		NmapXML, err := LoadXMLFile(filepath)
		if err != nil {
			log.Println(err)
			return
		}
		AllNmaps = append(AllNmaps, NmapXML)
	}
	for i := 0; i < len(AllNmaps); i++ {
		ParseNmapXML(AllNmaps[i])
	}
	var newcsv string
	if filepath == "" {
		newcsv = "NmapResults.csv"
	} else {
		newcsv = strings.TrimSuffix(filepath, ".xml")
		newcsv = newcsv + ".csv"
	}
	err := OutputToCSV(AllHostsAndPorts, newcsv)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("[+] Successfully Outputed To " + newcsv)
}

func GetAllXMLsInCurrentDirectory() ([]string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return []string{}, err
	}

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return []string{}, err
	}

	xmls := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".xml" {
			xmls = append(xmls, entry.Name())
		}
	}
	return xmls, nil
}

func ParseNmapXML(NmapXML NmapXML) {
	for i := 0; i < len(NmapXML.Host); i++ {
		hap := HostAndPorts{}
		for j := 0; j < len(NmapXML.Host[i].Ports.Port); j++ {
			if NmapXML.Host[i].Ports.Port[j].State.State == PORT_OPEN && NmapXML.Host[i].Ports.Port[j].Protocol == "tcp" {
				hap.Ports = append(hap.Ports, NmapXML.Host[i].Ports.Port[j].Portid)
			}
		}
		if len(hap.Ports) == 0 {
			continue
		}
		hap.Host = NmapXML.Host[i].Address[0].Addr //pray not ipv6
		AllHostsAndPorts = append(AllHostsAndPorts, hap)
	}
}

func LoadXMLFile(filepath string) (NmapXML, error) {
	xmlFile, err := os.Open(filepath)
	if err != nil {
		return NmapXML{}, err
	}
	defer xmlFile.Close()
	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return NmapXML{}, err
	}
	var NmapXML NmapXML
	err = xml.Unmarshal(byteValue, &NmapXML)
	if err != nil {
		return NmapXML, err
	}
	return NmapXML, nil
}

func OutputToCSV(hap []HostAndPorts, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"IP", "Ports"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, entry := range hap {
		ports := strings.Join(entry.Ports, " ") // Join ports with a semicolon or comma
		row := []string{entry.Host, ports}

		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
