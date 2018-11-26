package output

import (
	"fmt"
	"os"

	"we.com/vera.jiang/alertbeat/conf"
)

func checkFileExist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func serviceTemplate(f *os.File, proj, group string) error {
	servicesTemplate := `
define service{
	use                     passive-service
	host_name		`
	serviceTemplate := fmt.Sprintf("%v\t\t%v\n\tservice_description\t\t%v\n}", servicesTemplate, group, proj)
	n, _ := f.Seek(0, os.SEEK_END)
	_, err := f.WriteAt([]byte(serviceTemplate), n)
	if err != nil {
		return err
	}
	return nil
}

func hostTemplate(filename, group string) (*os.File, error) {
	Template := `
define host{
	use                     passive-server
	host_name		`
	hostTemplate := fmt.Sprintf("%v\t\t%v\n}", Template, group)

	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	_, err = f.Write([]byte(hostTemplate))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func createTemplateFile(proj, ttype string) error {
	dir := conf.Config.Nagios.TemplateDir
	var file *os.File
	var err error
	var filename, group, servicename string

	if ttype == "basic" {
		group = "BasicGroup"
		filename = dir + "/BasicGroup.cfg"
	} else if ttype == "web" {
		group = "WebGroup"
		filename = dir + "/WebGroup.cfg"
	} else {
		group = "JavaBizGroup"
		filename = dir + "/JavaBizGroup.cfg"
	}

	servicename = dir + "/" + proj

	if !checkFileExist(filename) {
		file, err = hostTemplate(filename, group)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(filename, os.O_RDWR, 0666)
		if err != nil {
			return err
		}
	}

	if checkFileExist(servicename) {
		return nil
	}
	f, err := os.Create(servicename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := serviceTemplate(file, proj, group); err != nil {
		return err
	}
	defer file.Close()
	return nil
}
