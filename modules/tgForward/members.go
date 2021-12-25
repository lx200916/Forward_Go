package tgForward

import (
	"encoding/json"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/orcaman/concurrent-map"
	"github.com/pkg/errors"
	"io/ioutil"
)

var log = utils.GetModuleLogger("external.TGForward.TGMember")

const FlagKey = "    "

// MembersMap key "    "`four spaces` in ConcurrentMap for changed flag true/false.
var MembersMap = make(map[int64]cmap.ConcurrentMap)

type SerializeMap map[int64]map[string]TGMember

func UpdateMembers() {
	data, err := ioutil.ReadFile("./data/members.json")
	if err != nil {
		log.Error(err)
		return
	}
	var SMembersMap SerializeMap
	err = json.Unmarshal(data, &SMembersMap)
	if err != nil {
		log.Error(err)
		return
	}

	for i, m := range SMembersMap {
		if entry, ok := MembersMap[i]; ok {
			for s, member := range m {
				entry.Set(s, member)
			}
			entry.Set(FlagKey, false)
		}

	}

}
func UpdateAMember(chatId int64, member TGMember) {
	name := fmt.Sprintf("%s %s", member.FirstName, member.LastName)
	if entry, ok := MembersMap[chatId]; ok {
		if !entry.Has(name) && name != FlagKey {
			entry.Set(name, member)
			entry.Set(FlagKey, true)
		}
	}
}
func GetAMember(chatId int64, name string) (TGMember, error) {
	if entry, ok := MembersMap[chatId]; ok {
		if entry.Has(name) && name != FlagKey {
			if inter, ok := entry.Get(name); ok {
				if member, ok := inter.(TGMember); ok {
					return member, nil
				}
			}
		}
	}
	return TGMember{}, errors.New("no member found")
}
func SerializeMembers() {
	serializedMembers := make(SerializeMap)
	changed := false

	for i, concurrentMap := range MembersMap {
		entry := make(map[string]TGMember)
		if flag, ok := concurrentMap.Get(FlagKey); ok && (flag.(bool)) {
			changed = true
		}
		for s, i2 := range concurrentMap.Items() {
			if s == FlagKey {
				continue
			}
			entry[s] = i2.(TGMember)
		}
		serializedMembers[i] = entry
	}
	if changed {
		marshal, err := json.Marshal(serializedMembers)
		if err != nil {
			log.Error(err)
			return
		}
		err = ioutil.WriteFile("./data/members.json", marshal, 0777)
		if err != nil {
			log.Error(err)
			return
		}
	}

}
