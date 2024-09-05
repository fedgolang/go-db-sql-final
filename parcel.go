package main

import (
	"database/sql"
	"log"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	add, err := s.db.Exec("INSERT INTO parcel(client, status, address, created_at) values(:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		log.Printf("Fail to add data: %v", err)
		return 0, err
	}
	id, err := add.LastInsertId()
	if err != nil {
		log.Printf("Fail to return LastInsertId: %v", err)
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	stmt := s.db.QueryRow("select * from parcel where number = :number",
		sql.Named("number", number),
	)
	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := stmt.Scan(&number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		log.Printf("Failed to get parcel by number: %v", err)
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	stmt, err := s.db.Query("select * from parcel where client = :client",
		sql.Named("client", client),
	)
	if err != nil {
		log.Printf("Failed to get parcel by client: %v", err)
		return nil, err
	}
	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for stmt.Next() {
		p := Parcel{}
		err = stmt.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			log.Printf("Failed to get parcel by client: %v", err)
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number),
	)
	if err != nil {
		log.Printf("Failed to set status: %v", err)
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		log.Printf("Failed to set status: %v", err)
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		log.Printf("Failed to delete parcel: %v", err)
		return err
	}
	return nil
}
