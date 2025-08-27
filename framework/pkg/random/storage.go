package random

import (
	"bufio"
	"os"
)

// Storage 负责随机码的存储
type Storage struct {
	filename string
	codes    map[string]struct{}
}

// NewStorage 创建一个新的Storage实例并从文件加载数据
func NewStorage(filename string) (*Storage, error) {
	storage := &Storage{
		filename: filename,
		codes:    make(map[string]struct{}),
	}
	if err := storage.load(); err != nil {
		return nil, err
	}
	return storage, nil
}

// load 从文件加载已存在的随机码
func (s *Storage) load() error {
	file, err := os.Open(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，表示这是第一次运行
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		code := scanner.Text()
		s.codes[code] = struct{}{}
	}

	return scanner.Err()
}

// Save 将新的随机码保存到文件并添加到内存中
func (s *Storage) Save(code string) error {
	file, err := os.OpenFile(s.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(code + "\n"); err != nil {
		return err
	}

	s.codes[code] = struct{}{}
	return nil
}

// Exists 检查随机码是否已经存在
func (s *Storage) Exists(code string) bool {
	_, exists := s.codes[code]
	return exists
}
