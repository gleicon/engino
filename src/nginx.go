package main

func createConf(ap *AppRecord) string {
	return "oi"
}

func activateSite(appname string) {

}

func updateNginxConf(appname string, conf string) error {
	// conferir a saida de nginx -t
	// usar lock
	return nil
}

func restartNginx() error {
	// TODO: embed throttling and restart batching
	// nginx -s reload
	// usar rwlock com ref count ou semaforo
	return nil
}
