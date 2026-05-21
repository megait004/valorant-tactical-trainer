import type { MapCatalogEntry } from '../../types'

const media = (uuid: string) => `https://media.valorant-api.com/maps/${uuid}/listviewicon.png`

const entry = (id: string, uuid: string, name: string): MapCatalogEntry => ({
  id,
  uuid,
  name,
  displayName: name,
  imageUrl: media(uuid),
  tacticalImageUrl: `https://media.valorant-api.com/maps/${uuid}/displayicon.png`,
  hasTacticalLayout: true,
})

export const fallbackMapCatalog: MapCatalogEntry[] = [
  entry('ascent', '7eaecc1b-4337-bbf6-6ab9-04b8f06b3319', 'Ascent'),
  entry('bind', '2c9d57ec-4431-9c5e-2939-8f9ef6dd5cba', 'Bind'),
  entry('breeze', '2fb9a4fd-47b8-4e7d-a969-74b4046ebd53', 'Breeze'),
  entry('fracture', 'b529448b-4d60-346e-e89e-00a4c527a405', 'Fracture'),
  entry('haven', '2bee0dc9-4ffe-519b-1cbd-7fbe763a6047', 'Haven'),
  entry('lotus', '2fe4ed3a-450a-948b-6d6b-e89a78e680a9', 'Lotus'),
  entry('pearl', 'fd267378-4d1d-484f-ff52-77821ed10dc2', 'Pearl'),
  entry('split', 'd960549e-485c-e861-8d71-aa9d1aed12a2', 'Split'),
  entry('sunset', '92584fbe-486a-b1b2-9faa-39b0f486b498', 'Sunset'),
  entry('abyss', '224b0a95-48b9-f703-1bd8-67aca101a61f', 'Abyss'),
]

export const mapIdFromName = (name: string) => name.trim().toLowerCase().replace(/\s+/g, '')
