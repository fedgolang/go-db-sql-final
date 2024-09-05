package main

import (
	"database/sql"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		log.Fatalf("Failed connect to DB: %v", err) // Fatal если коннект не удаётся соединить
	}
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	if err != nil {
		log.Printf("Failed to add new parcel: %v", err)
		return
	}
	require.NotEqual(t, id, 0)

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	p, err := store.Get(id)
	if err != nil {
		log.Printf("Failed to get parcel before delete: %v", err)
		return
	}

	// В задаче не описано, нужно ли ронять тест при несоответствии данных
	// Но так как require уже был импортирован, воспользуюсь им
	require.Equal(t, p.Address, parcel.Address)
	require.Equal(t, p.Client, parcel.Client)
	require.Equal(t, p.Status, parcel.Status)
	require.Equal(t, p.CreatedAt, parcel.CreatedAt)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(id)
	if err != nil {
		log.Printf("Failed to delete parcel: %v", err)
		return
	}

	_, err = store.Get(id)
	require.Equal(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		log.Fatalf("Failed connect to DB: %v", err) // Fatal если коннект не удаётся соединить
	}
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	if err != nil {
		log.Printf("Failed to add new parcel: %v", err)
		return
	}
	require.NotEqual(t, id, 0)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)
	if err != nil {
		log.Printf("Failed to set new address: %v", err)
		return
	}

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	p, err := store.Get(id)
	if err != nil {
		log.Printf("Failed to get parcel after address change: %v", err)
		return
	}
	require.Equal(t, p.Address, newAddress)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		log.Fatalf("Failed connect to DB: %v", err) // Fatal если коннект не удаётся соединить
	}
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	if err != nil {
		log.Printf("Failed to add new parcel: %v", err)
		return
	}
	require.NotEqual(t, id, 0)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(id, ParcelStatusDelivered)
	if err != nil {
		log.Printf("Failed to set status: %v", err)
		return
	}

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	p, err := store.Get(id)
	if err != nil {
		log.Printf("Failed to get parcel after status change: %v", err)
		return
	}
	require.Equal(t, p.Status, ParcelStatusDelivered)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		log.Fatalf("Failed connect to DB: %v", err) // Fatal если коннект не удаётся соединить
	}
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		if err != nil {
			log.Printf("Failed to add new parcel: %v", err)
			return
		}
		require.NotEqual(t, id, 0)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client 6808088

	// убедитесь в отсутствии ошибки
	if err != nil {
		log.Printf("Failed to get parcels by client: %v", err)
		return
	}

	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	require.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		if p, ok := parcelMap[parcel.Number]; ok {
			// убедитесь, что значения полей полученных посылок заполнены верно
			require.Equal(t, p.Client, parcel.Client)
			require.Equal(t, p.Status, parcel.Status)
			require.Equal(t, p.Address, parcel.Address)
			require.Equal(t, p.CreatedAt, parcel.CreatedAt)
		} else {
			log.Printf("Failed to find parcel in test map parcelMap: %v", err)
			return
		}
	}
}
