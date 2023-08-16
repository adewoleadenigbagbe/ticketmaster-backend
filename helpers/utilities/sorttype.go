package utilities

import "github.com/Wolechacho/ticketmaster-backend/database/entities"

type ByID[T entities.IDatabaseEntity] []T

func (byId ByID[T]) Len() int {
	return len(byId)
}

func (byId ByID[T]) Swap(i, j int) {
	byId[i], byId[j] = byId[j], byId[i]
}

func (byId ByID[T]) Less(i, j int) bool {
	return byId[i].GetId() < byId[j].GetId()
}
