package msg

import (
	"github.com/golang/protobuf/proto"
	message "gitlab.com/pangold/goim/msg/protobuf"
	"hash/crc32"
	"log"
	"sync"
)

type Package = message.Package
type Message = message.Message
type MessageType = message.Message_MessageType

const (
	MSG_JSON     = 1
	MSG_PROTOBUF = 2
	//
	PACKAGE_SIZE = 1024
)

var (
	msgId int64  = 0
	mutex        = sync.Mutex{}
	table32      = crc32.MakeTable(0xD5828281)

	// [id][index]
	// sentPackage [][]*Package //
	// append as soon as it's been sent, and wait for reply
	sentPackages map[int64][]*Package // maybe map[int64][]byte should be better.
	// receivedPackage [][]*Package
	// append as soon as it's been received, callback complete message or resend the missed packages
	receivedPackages map[int64][]*Package
)

type Buffer struct {

}

// procedure:
// A: send to B -> records it into sentList
// B: received message, reply to A -> records it into receivedList -> check if completed -> remove from receivedList
// A: received replied, remove from sentList, or resent if timeout

func getNextId() int64 {
	mutex.Lock()
	msgId++
	id := msgId
	mutex.Unlock()
	return id
}

//
func Serialize(m *Message) (result [][]byte) {
	body, err := proto.Marshal(m)

	if err != nil {
		log.Printf("serialize message error: %v", err)
	} else {
		result = split(body)
	}
	return result
}

//
func SerializePackage(p *Package) ([]byte, error) {
	*p.Checksum = crc32.Checksum(p.Body, table32)
	buf, err := proto.Marshal(p)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// split
func split(body []byte) (result [][]byte) {
	total := len(body) / (PACKAGE_SIZE + 1) + 1
	for index := 0; index < total; index++ {
		p := &Package{}
		*p.Name = "Message" // protobuf message name in Body
		// packages segment if too large
		begin := index * PACKAGE_SIZE
		end := (index + 1) * PACKAGE_SIZE
		if end > len(body) {
			end = len(body)
		}
		// id = atomic.AddInt64(&id, 1)
		*p.Index = int32(index)
		*p.Total = int32(total)
		*p.Count = int32(end - begin)
		p.Body = body[begin:end]
		buf, err := SerializePackage(p)
		if err != nil {
			log.Printf("serialize package(%d/%d) error: %v", index, total, err)
			return nil
		}
		result = append(result, buf)
		// TODO: records sent packages in a collection[n],
		//       clear package from collection[i] when received reply
		// TODO: resend mechanism, individually of course...
	}
	return result
}

func DeserializePackage(data []byte) *Package {
	p := &Package{}
	if err := proto.Unmarshal(data, p); err != nil {
		log.Printf("unmarshal error: %v", err)
		return nil
	}
	if crc32.Checksum(p.Body, table32) != *p.Checksum {
		log.Printf("invalid checksum")
		return nil
	}
	return p
}

func Deserialize(data []byte) (result []*Message) {
	var curIds []int64 = nil
	for pos := 0; pos * PACKAGE_SIZE < len(data); pos += PACKAGE_SIZE {
		p := DeserializePackage(data[pos:])
		// records received package
		receivedPackages[p.GetId()] = append(receivedPackages[p.GetId()], p)
		curIds = append(curIds, p.GetId())
	}
	// TODO:
	// then check package table
	for _, id := range curIds {
		if len(receivedPackages[id]) == int(receivedPackages[id][0].GetTotal()) {
			// sort
			// combine
			var buf []byte = nil
			for _, p := range receivedPackages[id] {
				buf = append(buf, p.GetBody()...)
			}
			// TODO: defer remove receivedPackages[id]
			m := &Message{}
			if err := proto.Unmarshal(buf, m); err != nil {
				log.Printf("unmarshal error: %v", err)
			}
			result = append(result, m)
		}
	}
	return result
}
