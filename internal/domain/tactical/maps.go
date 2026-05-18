package tactical

type MapCatalogEntry struct {
	ID                string `json:"id"`
	UUID              string `json:"uuid"`
	Name              string `json:"name"`
	DisplayName       string `json:"displayName"`
	ImageURL          string `json:"imageUrl"`
	TacticalImageURL  string `json:"tacticalImageUrl"`
	HasTacticalLayout bool   `json:"hasTacticalLayout"`
}

func AllMaps() []MapCatalogEntry {
	maps := []struct {
		id, uuid, name string
	}{
		{"ascent", "7eaecc1b-4337-bbf6-6ab9-04b8f06b3319", "Ascent"},
		{"bind", "2c9d57ec-4431-9c5e-2939-8f9ef6dd5cba", "Bind"},
		{"breeze", "2fb9a4fd-47b8-4e7d-a969-74b4046ebd53", "Breeze"},
		{"fracture", "b529448b-4d60-346e-e89e-00a4c527a405", "Fracture"},
		{"haven", "2bee0dc9-4ffe-519b-1cbd-7fbe763a6047", "Haven"},
		{"lotus", "2fe4ed3a-450a-948b-6d6b-e89a78e680a9", "Lotus"},
		{"pearl", "fd267378-4d1d-484f-ff52-77821ed10dc2", "Pearl"},
		{"split", "d960549e-485c-e861-8d71-aa9d1aed12a2", "Split"},
		{"sunset", "92584fbe-486a-b1b2-9faa-39b0f486b498", "Sunset"},
		{"abyss", "224b0a95-48b9-f703-1bd8-67aca101a61f", "Abyss"},
	}

	out := make([]MapCatalogEntry, 0, len(maps))
	for _, item := range maps {
		out = append(out, MapCatalogEntry{
			ID:                item.id,
			UUID:              item.uuid,
			Name:              item.name,
			DisplayName:       item.name,
			ImageURL:          "https://media.valorant-api.com/maps/" + item.uuid + "/listviewicon.png",
			TacticalImageURL:  "https://media.valorant-api.com/maps/" + item.uuid + "/displayicon.png",
			HasTacticalLayout: true,
		})
	}
	return out
}

func MapByID(mapID string) (MapCatalogEntry, bool) {
	for _, item := range AllMaps() {
		if item.ID == mapID {
			return item, true
		}
	}
	return MapCatalogEntry{}, false
}
