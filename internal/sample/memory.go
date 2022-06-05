package sample

import (
	"errors"
	"sync"
)

type memoryManager struct {
	categories map[int64]*Article
	rwlock     sync.RWMutex
	sequence   int64
}

var memmanager *memoryManager

func MemoryManagerInit() {
	memmanager = &memoryManager{
		categories: make(map[int64]*Article),
		sequence:   1,
	}
}

func CreateArticle(a *Article) error {
	memmanager.rwlock.Lock()
	defer memmanager.rwlock.Unlock()

	if _, ok := memmanager.categories[a.Id]; ok {
		return errors.New("Already exists")
	} else {
		a.Id = memmanager.sequence
		memmanager.categories[a.Id] = a
		memmanager.sequence++
	}
	return nil
}

func GetArticle(id int64) (*Article, error) {

	memmanager.rwlock.RLock()
	defer memmanager.rwlock.RUnlock()
	if val, ok := memmanager.categories[id]; ok {
		return val, nil
	} else {
		return nil, errors.New("Not Found.")
	}
}

func DeleteArticle(id int64) error {
	memmanager.rwlock.RLock()
	defer memmanager.rwlock.RUnlock()
	if _, ok := memmanager.categories[id]; ok {
		delete(memmanager.categories, id)
	} else {
		return errors.New("Not Found.")
	}
	return nil
}

func GetArticles() ([]Article, error) {
	v := make([]Article, 0, len(memmanager.categories))
	for _, value := range memmanager.categories {
		v = append(v, *value)
	}
	return v, nil
}

func UpdateArticle(a *Article) error {
	memmanager.rwlock.RLock()
	defer memmanager.rwlock.RUnlock()
	if val, ok := memmanager.categories[a.Id]; ok {
		if a.Author != "" && a.Author != val.Author{
			a.Author = val.Author
		}
		if a.Views != 0 && a.Views != val.Views{
			a.Views = val.Views
		}
		if a.Title != "" && a.Title != val.Title{
			a.Title = val.Title
		}
	} else {
		return errors.New("Not Found.")
	}
	return nil
}
