package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"github.com/k0yote/privatechain/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}

type TxHasher struct{}

func (TxHasher) Hash(tx *Transaction) types.Hash {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, tx.Data)
	_ = binary.Write(buf, binary.LittleEndian, tx.To)
	_ = binary.Write(buf, binary.LittleEndian, tx.Value)
	_ = binary.Write(buf, binary.LittleEndian, tx.From)
	_ = binary.Write(buf, binary.LittleEndian, tx.Nonce)

	return types.Hash(sha256.Sum256(buf.Bytes()))
}
