package sample

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

type fileManager struct {
	filename        string
	offsets         map[int64]int64
	deleted_offsets []int64
	rwlock          sync.RWMutex
	sequence        int64
	recordsCount    int64
	current_offset  int64
}

const blockSize int64 = 1024

var manager *fileManager

func FileManagerInit(filename string) {
	manager = &fileManager{
		filename:        filename,
		offsets:         make(map[int64]int64),
		deleted_offsets: make([]int64, 0),
		sequence:        1,
	}
}

func buildOffset() error {

	g := new(errgroup.Group)

	fi, err := os.Stat(manager.filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			file, err := os.Create(manager.filename)
			if err != nil {
				return err
			}
			defer file.Close()

			err = binary.Write(file, binary.LittleEndian, manager.sequence)
			if err != nil {
				return err
			} else {
				manager.current_offset = 8
				return nil
			}
		} else {
			return err
		}
	}
	size := fi.Size()
	size -= 8

	file, err := os.Open(manager.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	current_offset, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	manager.current_offset = current_offset

	if size < StrutSizeArticles {
		return nil
	}

	rm := size % (StrutSizeArticles * blockSize)
	blocks := int64(size / (StrutSizeArticles * blockSize))

	var i int64
	for i = 0; i < blocks; i++ {
		offset := (i * StrutSizeArticles * blockSize) + 8

		g.Go(func() error {

			var category Category

			file, err := os.Open(manager.filename)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := file.Seek(offset, io.SeekStart); err != nil {
				return err
			}

			var i int64
			for i = 0; i < blockSize; i++ {
				err = binary.Read(file, binary.LittleEndian, category)
				if err != nil {
					return err
				}
				if category.Id != 0 {
					manager.offsets[category.Id] = offset + (StrutSizeArticles * i)
					manager.recordsCount++
				} else {
					manager.deleted_offsets = append(manager.deleted_offsets, offset+(StrutSizeArticles*i))
				}
			}

			return nil
		})
	}

	offset := (i * StrutSizeArticles * blockSize) + 8
	records := rm / StrutSizeArticles

	//g.Go(func() error {

	var category Category

	file, err = os.Open(manager.filename)
	//file, err := os.Open(manager.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	//		var i int64
	for i = 0; i < records; i++ {
		err = binary.Read(file, binary.LittleEndian, &category)
		if err != nil {
			return err
		}

		if category.Id != 0 {
			manager.offsets[category.Id] = offset + (StrutSizeArticles * i)
		} else {
			manager.deleted_offsets = append(manager.deleted_offsets, offset+(StrutSizeArticles*i))
		}
	}

	return nil
	//})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func GetRecord(id int64) (*Category, error) {

	category := &Category{}
	manager.rwlock.RLock()
	defer manager.rwlock.RUnlock()

	offset, found := manager.offsets[category.Id]
	if !found {
		return nil, mux.ErrNotFound
	}

	file, err := os.Open(manager.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	err = binary.Read(file, binary.LittleEndian, category)
	if err != nil {
		return nil, err
	}

	return category, nil

}

func getAllRecord(id int64) ([]Category, error) {

	articles := make([]Category, 0, manager.recordsCount)

	manager.rwlock.RLock()
	defer manager.rwlock.RUnlock()

	file, err := os.Open(manager.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for {
		category := &Category{}
		err = binary.Read(file, binary.LittleEndian, category)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		articles = append(articles, *category)
	}

	return articles, nil
}

func UpdateRecord(category *Category) error {

	manager.rwlock.Lock()
	defer manager.rwlock.Unlock()

	offset, found := manager.offsets[category.Id]
	if !found {
		return mux.ErrNotFound
	}

	file, err := os.Open(manager.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	err = binary.Write(file, binary.LittleEndian, category)
	if err != nil {
		return err
	}

	return nil

}

func DeleteRecord(id int64) error {

	category := &Category{Id: 0}
	manager.rwlock.Lock()
	defer manager.rwlock.Unlock()

	offset, found := manager.offsets[id]
	if !found {
		return mux.ErrNotFound
	}

	file, err := os.Open(manager.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	delete(manager.offsets, id)
	manager.deleted_offsets = append(manager.deleted_offsets, offset)
	err = binary.Write(file, binary.LittleEndian, category)
	if err != nil {
		return err
	}

	return nil

}

func CreateRecord(category *Category) error {

	manager.rwlock.Lock()
	defer manager.rwlock.Unlock()

	file, err := os.Create(manager.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if len(manager.deleted_offsets) != 0 {
		discarded_offset := manager.deleted_offsets[0]
		manager.deleted_offsets = manager.deleted_offsets[1:]

		if _, err := file.Seek(discarded_offset, io.SeekStart); err != nil {
			return err
		}
	} else {
		if _, err := file.Seek(manager.current_offset, io.SeekStart); err != nil {
			return err
		}
		manager.current_offset += StrutSizeArticles
	}

	category.Id = manager.sequence
	manager.sequence++

	err = binary.Write(file, binary.LittleEndian, category)
	if err != nil {
		return err
	}

	offset_temp, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	fmt.Println(offset_temp)

	err = binary.Write(file, binary.LittleEndian, &manager.sequence)
	if err != nil {
		return err
	}

	return nil

}
