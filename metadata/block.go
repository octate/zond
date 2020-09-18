package metadata

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/theQRL/zond/config"
	"github.com/theQRL/zond/db"
	"github.com/theQRL/zond/misc"
	"github.com/theQRL/zond/protos"
	"go.etcd.io/bbolt"
)

type BlockMetaData struct {
	pbData *protos.BlockMetaData
}

func (bm *BlockMetaData) ParentHeaderHash() []byte {
	return bm.pbData.ParentHeaderHash
}

func (bm *BlockMetaData) ChildHeaderHashes() [][]byte {
	return bm.pbData.ChildHeaderHashes
}

func (bm *BlockMetaData) FinalizedChildHeaderHash() []byte {
	return bm.pbData.FinalizedChildHeaderHash
}

func (bm *BlockMetaData) HeaderHash() []byte {
	return bm.pbData.HeaderHash
}

func (bm *BlockMetaData) SlotNumber() uint64 {
	return bm.pbData.SlotNumber
}

func (bm *BlockMetaData) Epoch() uint64 {
	return bm.pbData.SlotNumber / config.GetDevConfig().BlocksPerEpoch
}

func (bm *BlockMetaData) Serialize() ([]byte, error) {
	return proto.Marshal(bm.pbData)
}

func (bm *BlockMetaData) DeSerialize(data []byte) error {
	return proto.Unmarshal(data, bm.pbData)
}

func (bm *BlockMetaData) AddChildHeaderHash(headerHash []byte) {
	bm.pbData.ChildHeaderHashes = append(bm.pbData.ChildHeaderHashes,
		headerHash)
}

func (bm *BlockMetaData) Commit(b *bbolt.Bucket) error {
	data, err := bm.Serialize()
	if err != nil {
		return err
	}
	return b.Put(GetBlockMetaDataKey(bm.HeaderHash()), data)
}

func NewBlockMetaData(parentHeaderHash []byte, headerHash []byte, slotNumber uint64) *BlockMetaData {
	pbData := &protos.BlockMetaData {
		ParentHeaderHash: parentHeaderHash,
		HeaderHash: headerHash,
		SlotNumber: slotNumber,
	}
	return &BlockMetaData {
		pbData: pbData,
	}
}

func GetBlockMetaData(d *db.DB, headerHash []byte) (*BlockMetaData, error) {
	key := GetBlockMetaDataKey(headerHash)
	data, err := d.Get(key)
	if err != nil {
		log.Error("Error loading BlockMetaData for key ", string(key), err)
		return nil, err
	}
	bm := &BlockMetaData{
		pbData: &protos.BlockMetaData{},
	}
	return bm, bm.DeSerialize(data)
}

func GetBlockMetaDataKey(headerHash []byte) []byte {
	return []byte(fmt.Sprintf("BLOCK-META-DATA-%s", misc.Bin2HStr(headerHash)))
}

func GetBlockBucketName(blockHeaderHash []byte) []byte {
	return []byte(fmt.Sprintf("BLOCK-BUCKET-%s", misc.Bin2HStr(blockHeaderHash)))
}