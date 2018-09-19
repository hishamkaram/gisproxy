package gisproxy

import (
	"encoding/xml"
	"fmt"
)

//ServiceException WMS ServiceException
type ServiceException struct {
	XMLName xml.Name `xml:"ServiceException"`
	Message string   `xml:",chardata"`
}

//ServiceExceptionReport WMS ServiceExceptionReport
type ServiceExceptionReport struct {
	XMLName    xml.Name           `xml:"ServiceExceptionReport"`
	Version    string             `xml:"version,attr"`
	Exceptions []ServiceException `xml:"ServiceException"`
}

// GenerateExceptionReport Get Exception as XML
func (proxyServer *GISProxy) GenerateExceptionReport(report *ServiceExceptionReport) (xmlBody []byte, err error) {
	if xmlBody, err = xml.MarshalIndent(report, "		", "	"); err == nil {
		header := fmt.Sprintf("%s%s\n%s\n", xml.Header, []byte(`<!DOCTYPE ServiceExceptionReport SYSTEM "http://www.digitalearth.gov/wmt/xml/exception_1_1_1.dtd">`), xmlBody)
		xmlBody = []byte(header)
	}
	return
}
