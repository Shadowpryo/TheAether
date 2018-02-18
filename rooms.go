package main

import (
	"sync"
	"errors"
)

type Rooms struct {
	db          *DBHandler
	querylocker sync.RWMutex
}

type Room struct {

	ID string `storm:"id"` // primary key

	Name 				string
	ParentID			string

	RoleID				string


	// Connecting Room ID's
	UpID				string
	UpItemID			[]string

	DownID				string
	DownItemID			[]string

	NorthID				string
	NorthItemID			[]string

	NorthEastID			string
	NorthEastItemID		[]string

	EastID				string
	EastItemID			[]string

	SouthEastID			string
	SouthEastItemID		[]string

	SouthID				string
	SouthItemID			[]string

	SouthWestID			string
	SouthWestItemID		[]string

	WestID				string
	WestItemID			[]string

	NorthWestID 		string
	NorthWestItemID		[]string


}



func (h *Rooms) SaveRoomToDB(room Room) (err error) {
	h.querylocker.Lock()
	defer h.querylocker.Unlock()

	db := h.db.rawdb.From("Rooms")
	err = db.Save(&room)
	return err
}

func (h *Rooms) RemoveRoomFromDB(room Room) (err error) {
	h.querylocker.Lock()
	defer h.querylocker.Unlock()

	db := h.db.rawdb.From("Rooms")
	err = db.DeleteStruct(&room)
	return err
}

func (h *Rooms) RemoveRoomByID(roomID string) (err error) {

	room, err := h.GetRoomByID(roomID)
	if err != nil {
		return err
	}

	err = h.RemoveRoomFromDB(room)
	if err != nil {
		return err
	}

	return nil
}

func (h *Rooms) GetRoomByID(roomID string) (room Room, err error) {

	rooms, err := h.GetAllRooms()
	if err != nil{
		return room, err
	}

	for _, i := range rooms {
		if i.ID == roomID{
			return i, nil
		}
	}

	return room, errors.New("No record found")
}


// GetAllRooms function
func (h *Rooms) GetAllRooms() (roomlist []Room, err error) {
	h.querylocker.Lock()
	defer h.querylocker.Unlock()

	db := h.db.rawdb.From("Rooms")
	err = db.All(&roomlist)
	if err != nil {
		return roomlist, err
	}

	return roomlist, nil
}
