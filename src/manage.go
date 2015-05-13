package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/fiorix/go-redis/redis"
	"github.com/hashicorp/consul/api"
)

/*
AppRecord is the basic unit of application info
*/
type AppRecord struct {
	AppName      string
	CreatedAt    int64
	LastChanged  int64
	Activated    int64
	TemplateName string
	Active       bool
	Attributes   map[string]string
}

var (
	ENGINO_PATH = "engino"
)

/*
NewAppRecord creates a new application record from application name
*/
func NewAppRecord(appName string) *AppRecord {
	keys := make(map[string]string)
	return &AppRecord{AppName: appName, Attributes: keys}
}

func (ap *AppRecord) fillAppRecordFromMap(appData map[string]string) error {
	for key, value := range appData {
		switch key {
		case "appname":
			ap.AppName = value
			break
		case "template":
			ap.TemplateName = value
			break
		case "created_at":
			s, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				log.Println("Error converting date: ", value)
				return err
			} else {
				ap.CreatedAt = s
			}
			break
		case "last_changed":
			s, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				log.Println("Error converting date: ", value)
				return err
			} else {
				ap.LastChanged = s
			}
			break
		case "active":
			if strings.ToLower(value) == "true" {
				ap.Active = true
			} else {
				ap.Active = false
			}
			break
		default:
			ap.Attributes[key] = value
			break
		}

	}
	return nil
}

func (ap *AppRecord) fillAppRecordFromRedis(addr string) error {
	rc := redis.New(addr)
	defer rc.CloseAll()
	keys, err := rc.HGetAll(ap.AppName)
	if err != nil {
		return err
	}
	return ap.fillAppRecordFromMap(keys)
}

func (ap *AppRecord) fillAppRecordFromEtcd(addr string) error {
	var etcdClient = etcd.NewClient([]string{addr})

	path := fmt.Sprintf("%s/%s", ap.AppName, ENGINO_PATH)
	res, err := etcdClient.Get(path, true, false)

	if err != nil {
		return err
	}
	keys := make(map[string]string)

	for _, n := range res.Node.Nodes {
		k := strings.Split(n.Key, "/")
		key, value := strings.ToUpper(k[len(k)-1]), n.Value
		keys[key] = value
	}
	return ap.fillAppRecordFromMap(keys)

}

func (ap *AppRecord) fillAppRecordFromConsul(addr string) error {
	df := api.DefaultConfig()
	df.Address = addr
	client, _ := api.NewClient(df)
	kv := client.KV()
	path := fmt.Sprintf("%s/%s", ap.AppName, ENGINO_PATH)
	keys, _, err := kv.List(path, nil)
	if err != nil {
		return err
	}

	keymap := make(map[string]string)
	for _, kp := range keys {
		kk := strings.Split(kp.Key, "/")
		key, value := strings.ToUpper(kk[len(kk)-1]), kp.Value
		keymap[key] = string(value)
	}
	return ap.fillAppRecordFromMap(keymap)
}

func manageVHost(connection url.URL, appName string, templateDir string, throttling int, nginxConfDir string) {
	// check if boolean activated is true before looping to toggle vhost on/off
	ap := NewAppRecord(appName)
	firstTime := true

	for {
		// TODO: check *activated* for status toggle
		if firstTime || newConfigData(ap) {
			if firstTime {
				firstTime = false
			}
			switch connection.Scheme {
			case "redis":
				ap.fillAppRecordFromRedis(connection.String())
				break
			case "consul":
				ap.fillAppRecordFromConsul(connection.String())
				break
			case "etcd":
				ap.fillAppRecordFromEtcd(connection.String())
				break
			default:
				fmt.Println("Error unknown connection string")
				panic(connection)

			}
			conf := createConf(ap)
			err := updateNginxConf(appName, conf)
			if err != nil {
				log.Fatalln("Error updating nginx config ", err)
			}

			err = restartNginx()
			if err != nil {
				log.Fatalln("Error restarting nginx", err)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func newConfigData(ap *AppRecord) bool {
	return true
}
