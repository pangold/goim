package plugin

import (
	"log"
)

func SyncSession(state, token string) error {
	log.Println("sync session to backend service.")
	return nil
}

func CheckToken(token string) bool {
	log.Println("check token.")
	return true
}

