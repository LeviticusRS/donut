package asset

import (
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
)

const (
    ReferenceLength    = 6
    BlockHeaderLength  = 8
    BlockPayloadLength = 512
    BlockLength        = BlockHeaderLength + BlockPayloadLength
    EndOfArchive       = 0
)

type Cache struct {
    blocks        *os.File
    manifestIndex *os.File
    indexes       []*os.File
    mutex         sync.Mutex
    buffer        [BlockLength]byte
}

func OpenCache(root string, count int) (*Cache, error) {
    blocks, err := os.Open(filepath.Join(root, "main_file_cache.dat2"))
    if err != nil {
        return nil, err
    }

    manifestIndex, err := os.Open(filepath.Join(root, fmt.Sprintf("main_file_cache.idx%d", Manifests)))
    if err != nil {
        return nil, err
    }

    indexes := make([]*os.File, count)
    for i := 0; i < count; i++ {
        index, err := os.Open(filepath.Join(root, fmt.Sprintf("main_file_cache.idx%d", i)))
        if err != nil {
            return nil, err
        }
        indexes[i] = index
    }

    return &Cache{
        blocks:        blocks,
        manifestIndex: manifestIndex,
        indexes:       indexes,
        mutex:         sync.Mutex{},
    }, nil
}

func (c *Cache) Get(index uint8, id uint16) ([]byte, error) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    indexFile, err := c.getIndexFile(index)
    if err != nil {
        return nil, err
    }

    indexLength, err := fileLength(indexFile)
    if err != nil {
        return nil, err
    }

    if indexLength < ReferenceLength+ReferenceLength*int64(id) {
        return nil, errors.New("asset: archive does not exist")
    }

    if _, err := indexFile.ReadAt(c.buffer[:ReferenceLength], int64(id)*ReferenceLength); err != nil {
        return nil, err
    }

    length := uint32(c.buffer[0])<<16 | uint32(c.buffer[1])<<8 | uint32(c.buffer[2])
    block := uint32(c.buffer[3])<<16 | uint32(c.buffer[4])<<8 | uint32(c.buffer[5])

    blocksLength, err := fileLength(c.blocks)
    if err != nil {
        return nil, err
    }

    if block <= EndOfArchive || int64(block) > blocksLength/int64(BlockLength) {
        return nil, io.EOF
    }

    result := make([]byte, length)

    offset := uint32(0)
    part := uint16(0)

    for offset < length {
        if block == EndOfArchive {
            return nil, errors.New("asset: premature end of archive")
        }

        read := length - offset
        if read > BlockPayloadLength {
            read = BlockPayloadLength
        }

        if _, err := c.blocks.ReadAt(c.buffer[:BlockHeaderLength+read], int64(block)*BlockLength); err != nil {
            return nil, err
        }

        blockArchiveId := uint16(c.buffer[0])<<8 | uint16(c.buffer[1])
        blockArchiveChunk := uint16(c.buffer[2])<<8 | uint16(c.buffer[3])
        blockIndex := c.buffer[7]

        if blockArchiveId != id || blockArchiveChunk != part || blockIndex != index {
            return nil, errors.New("asset: invalid block header")
        }

        nextBlock := uint32(c.buffer[4])<<16 | uint32(c.buffer[5])<<8 | uint32(c.buffer[6])

        if int64(block) > blocksLength/int64(BlockLength) {
            return nil, io.EOF
        }

        copy(result[offset:], c.buffer[BlockHeaderLength:BlockHeaderLength+read])

        block = nextBlock
        offset += read
        part++
    }

    return result, nil
}

func (c *Cache) IndexCount() int {
    return len(c.indexes)
}

func (c *Cache) getIndexFile(index uint8) (*os.File, error) {
    if index == Manifests {
        return c.manifestIndex, nil
    }

    if len(c.indexes) < int(index) {
        return nil, fmt.Errorf("asset: cache does not contain index %d", index)
    }

    return c.indexes[index], nil
}

func fileLength(file *os.File) (int64, error) {
    info, err := file.Stat()
    if err != nil {
        return 0, err
    }
    return info.Size(), nil
}
