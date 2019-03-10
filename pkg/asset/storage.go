package asset

import "sync"

type Storage struct {
    cache    *Cache
    archives map[uint8]map[uint16][]byte
    mutex    sync.Mutex
}

func NewStorage(cache *Cache) (*Storage, error) {
    storage := Storage{
        cache:    cache,
        archives: make(map[uint8]map[uint16][]byte, cache.IndexCount()),
        mutex:    sync.Mutex{},
    }

    storage.archives[Manifests] = make(map[uint16][]byte)

    for i := 0; i < cache.IndexCount(); i++ {
        storage.archives[uint8(i)] = make(map[uint16][]byte)
    }

    // The cache never contains the release manifest, so generate it and store it.
    // TODO(hadyn): Consider moving out of this functionality.
    manifest, err := CreateReleaseManifest(cache)
    if err != nil {
        return nil, err
    }

    packed, err := CompressArchive(Uncompressed, manifest)
    if err != nil {
        return nil, err
    }
    storage.archives[Manifests][255] = packed

    return &storage, nil
}

func (s *Storage) Get(index uint8, id uint16) ([]byte, error) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    if b, exists := s.archives[index][id]; exists {
        return b, nil
    }

    return s.retrieve(index, id)
}

func (s *Storage) retrieve(index uint8, id uint16) ([]byte, error) {
    b, err := s.cache.Get(index, id)
    if err != nil {
        return nil, err
    }

    // Trimming the archive fixes an issue with caches containing a footer which is used to hold to archive's version
    // to compare against the manifest.
    b, err = TrimArchive(b)
    if err != nil {
        return nil, err
    }

    s.archives[index][id] = b

    return b, nil
}
