package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func help() {
	fmt.Println("Usage: engino ...")
	fmt.Println("	-c connection string. redis://localhost:6379, etcd://localhost:4001, consul://localhost:8500")
	fmt.Println("	-a application names. you can have more than one application name by separating them with a comma: -a app1,app2,app3.")
	fmt.Println("			this is going to be interpreted as a key (engino:appname)")
	fmt.Println("			on redis and folder on etc/consul with a subfolder called engino (/appname/engino/...)")
	fmt.Println("	-t template directory")
	fmt.Println("	-r throttling, number of restarts per minute")
	fmt.Println("	-n nginx config dir")
	fmt.Println("$ engino -c redis://localhost:6379 -a app1,app2 -t /opt/engino/templates")
	os.Exit(1)
}

func main() {
	connURL := flag.String("c", "", "Connection URL")
	appName := flag.String("a", "", "Application name")
	templateDir := flag.String("t", "", "Template dir")
	throttling := flag.Int("r", 4, "Connections per minute")
	nginxConfDir := flag.String("n", "/etc/nginx/", "Application name")
	flag.Usage = help
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Check your head yo - no parameters given")
		help()
		os.Exit(1)
	}

	connection, err := url.Parse(*connURL)
	if err != nil {
		panic(err)
	}

	apps := strings.Split(*appName, ",")
	for _, app := range apps {
		go manageVHost(*connection, app, *templateDir, *throttling, *nginxConfDir)
	}

}
