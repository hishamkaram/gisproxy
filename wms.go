package gisproxy

import (
	"encoding/xml"
	"log"
)

//OnlineResource tag
type OnlineResource struct {
	XMLName xml.Name `xml:"OnlineResource"`
	Type    string   `xml:"http://www.w3.org/1999/xlink type,attr,omitempty"`
	Href    string   `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
}

//KeywordList tag
type KeywordList struct {
	XMLName xml.Name  `xml:"KeywordList"`
	Keyword []*string `xml:"Keyword"`
}

//Get tag
type Get struct {
	XMLName        xml.Name       `xml:"Get"`
	OnlineResource OnlineResource `xml:"OnlineResource"`
}

//Post tag
type Post struct {
	XMLName        xml.Name       `xml:"Post"`
	OnlineResource OnlineResource `xml:"OnlineResource"`
}

//HTTP tag
type HTTP struct {
	XMLName xml.Name `xml:"HTTP"`
	Get     *Get     `xml:"Get"`
	Post    *Post    `xml:"Post"`
}

//DCPType tag
type DCPType struct {
	XMLName xml.Name `xml:"DCPType"`
	HTTP    HTTP     `xml:"HTTP"`
}

//RequestEntry tag
type RequestEntry struct {
	Format  []*string `xml:"Format"`
	DCPType DCPType   `xml:"DCPType"`
}

//Request tag
type Request struct {
	XMLName          xml.Name     `xml:"Request"`
	GetCapabilities  RequestEntry `xml:"GetCapabilities"`
	GetMap           RequestEntry `xml:"GetMap"`
	GetFeatureInfo   RequestEntry `xml:"GetFeatureInfo"`
	DescribeLayer    RequestEntry `xml:"DescribeLayer"`
	GetLegendGraphic RequestEntry `xml:"GetLegendGraphic"`
	GetStyles        RequestEntry `xml:"GetStyles"`
}

//UserDefinedSymbolization tag
type UserDefinedSymbolization struct {
	XMLName    xml.Name `xml:"UserDefinedSymbolization"`
	SupportSLD string   `xml:"SupportSLD,attr"`
	UserLayer  string   `xml:"UserLayer,attr"`
	UserStyle  string   `xml:"UserStyle,attr"`
	RemoteWFS  string   `xml:"RemoteWFS,attr"`
}

//Exception tag
type Exception struct {
	XMLName xml.Name  `xml:"Exception"`
	Format  []*string `xml:"Format"`
}

//AuthorityURL tag
type AuthorityURL struct {
	XMLName        xml.Name       `xml:"AuthorityURL"`
	Name           string         `xml:"name,attr"`
	OnlineResource OnlineResource `xml:"OnlineResource"`
}

//LatLonBoundingBox tag
type LatLonBoundingBox struct {
	XMLName xml.Name `xml:"LatLonBoundingBox"`
	MinX    float64  `xml:"minx,attr"`
	MinY    float64  `xml:"miny,attr"`
	MaxX    float64  `xml:"maxx,attr"`
	MaxY    float64  `xml:"maxy,attr"`
}

//Layer tag
type Layer struct {
	XMLName           xml.Name          `xml:"Layer"`
	Title             string            `xml:"Title"`
	Abstract          string            `xml:"Abstract"`
	SRS               []*string         `xml:"SRS"`
	LatLonBoundingBox LatLonBoundingBox `xml:"LatLonBoundingBox"`
	AuthorityURL      AuthorityURL      `xml:"AuthorityURL"`
}

//Capability tag
type Capability struct {
	XMLName                  xml.Name                 `xml:"Capability"`
	Request                  Request                  `xml:"Request"`
	Exception                Exception                `xml:"Exception"`
	UserDefinedSymbolization UserDefinedSymbolization `xml:"UserDefinedSymbolization"`
	Layer                    Layer                    `xml:"Layer"`
}

//Service tag
type Service struct {
	XMLName           xml.Name       `xml:"Service"`
	Name              string         `xml:"Name"`
	Title             string         `xml:"Title"`
	KeywordList       KeywordList    `xml:"KeywordList"`
	OnlineResource    OnlineResource `xml:"OnlineResource"`
	Fees              string         `xml:"Fees"`
	AccessConstraints string         `xml:"AccessConstraints"`
}

//WMSCapabilities parent tag
type WMSCapabilities struct {
	XMLName        xml.Name   `xml:"WMT_MS_Capabilities"`
	Version        string     `xml:"version,attr,omitempty"`
	UpdateSequence string     `xml:"updateSequence,attr,omitempty"`
	Service        Service    `xml:"Service"`
	Capability     Capability `xml:"Capability"`
}

//ParseWMSCapabilities read wms capabilities
func ParseWMSCapabilities(xmlByte []byte) *WMSCapabilities {
	var cap WMSCapabilities
	err := xml.Unmarshal(xmlByte, &cap)
	if err != nil {
		log.Fatal(err)
	}
	return &cap
}
