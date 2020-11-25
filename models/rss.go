package models

import (
	"crypto/md5"
	"fmt"
)

//Attribute ...
type Attribute struct {
	number int
	object *XMLTree
}

//Item ...
type Item struct {
	number        int
	object        *XMLTree
	attributesMap map[string]*Attribute
}

//Channel ...
type Channel struct {
	object   *XMLTree
	itemsMap map[string]*Item
}

//RSS ...
type RSS struct {
	RSSXMLTree *XMLTree
	RSSChannel *Channel
}

//RssFromXML ...
func RssFromXML(rootXML *XMLTree) *RSS {
	rss := new(RSS)
	rss.RSSXMLTree = rootXML
	rss.RSSChannel = new(Channel)
	channel := rss.RSSChannel
	channel.itemsMap = make(map[string]*Item)
	channelTree := rootXML.Branches[0]
	for i, channelBranch := range channelTree.Branches {
		item := new(Item)
		item.number = i
		item.object = channelBranch
		item.attributesMap = make(map[string]*Attribute)
		for j, itemBranch := range channelBranch.Branches {
			attribute := new(Attribute)
			attribute.number = j
			attribute.object = itemBranch
			item.attributesMap[HashOfString(itemBranch.Name)] = attribute
		}
		var key string
		if channelBranch.Name == "item" {
			key = item.attributesMap[HashOfString("guid")].object.Value
		} else {
			key = channelBranch.Name
		}
		channel.itemsMap[HashOfString(key)] = item
	}
	return rss
}

//HashOfString ...
func HashOfString(s string) string {
	contentBytes := []byte(s)
	return fmt.Sprintf("%x", md5.Sum(contentBytes))
}

//RssCompare ...
func RssCompare(oldRSS RSS, newRSS *RSS) {
	newRSSChannel := newRSS.RSSChannel
	oldRSSChannel := oldRSS.RSSChannel

	if newRSS.RSSXMLTree.ContentHash == oldRSS.RSSXMLTree.ContentHash {
		newRSS.RSSXMLTree.Status = OLD
		newRSSChannel.object.Status = OLD
	} else {
		newRSS.RSSXMLTree.Status = NEW
		newRSSChannel.object.Status = NEW
	}
	for id, newRSSItem := range newRSSChannel.itemsMap {
		oldRSSItem, oldItemStatus := oldRSSChannel.itemsMap[id]
		// true если данные есть в двух rss
		if oldItemStatus {
			// проверка на несовпадение содержания в двух item-x
			if oldRSSItem.object.ContentHash != newRSSItem.object.ContentHash {
				if newRSSItem.number != oldRSSItem.number {
					newRSSItem.object.Status = CHANGEDVALUEORDER
				} else {
					newRSSItem.object.Status = CHANGEDVALUE
				}
				for jd, newRSSattr := range newRSSItem.attributesMap {
					oldRSSattr, oldAttrStatus := oldRSSItem.attributesMap[jd]
					// true если данные есть в двух rss
					if oldAttrStatus {
						// проверка на совпадение содержания в двух атрибутах
						if newRSSattr.object.ContentHash != oldRSSattr.object.ContentHash {
							newRSSattr.object.Status = NEW
						} else {
							if newRSSattr.number != oldRSSattr.number {
								newRSSattr.object.Status = CHANGEDORDER
							} else {
								newRSSattr.object.Status = OLD
							}
						}
					} else {
						newRSSattr.object.Status = NEW
					}
				}
				for jd, oldRSSattr := range oldRSSItem.attributesMap {
					_, newAttrStatus := newRSSItem.attributesMap[jd]
					// если данные отсутствуют
					if !newAttrStatus {
						currentItem := newRSSItem.object
						deletedAttr := new(XMLTree)
						*deletedAttr = *oldRSSattr.object
						deletedAttr.Status = DELETED
						currentItem.Branches = append(currentItem.Branches, deletedAttr)
					}
				}
			} else {
				if newRSSItem.number != oldRSSItem.number {
					newRSSItem.object.Status = CHANGEDORDER
				} else {
					newRSSItem.object.Status = OLD
				}
			}
		} else {
			newRSSItem.object.Status = NEW
		}
	}
	//добавляем удаленные item-ы c статусом DELETED
	for id, oldRSSItem := range oldRSSChannel.itemsMap {
		_, newItemStatus := newRSSChannel.itemsMap[id]
		// если данные отсутствуют
		if !newItemStatus {
			channel := newRSS.RSSXMLTree.Branches[0]
			deletedItem := new(XMLTree)
			*deletedItem = *oldRSSItem.object
			deletedItem.Status = DELETED
			channel.Branches = append(channel.Branches, deletedItem)
		}

	}
}
