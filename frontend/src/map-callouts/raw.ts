/**
 * Dữ liệu callout thô cho từng map, lấy từ valorant-api.com/v1/maps.
 *
 * Mỗi map có công thức convert tọa độ Unreal Engine → minimap %:
 *   minimapX = (ueY * xMultiplier + xScalarToAdd) * 100
 *   minimapY = (ueX * yMultiplier + yScalarToAdd) * 100
 *
 * Công thức này áp dụng cho ảnh tactical chính thức (displayIcon) của Riot.
 * Mỗi ảnh user đã download nằm ở `public/maps/{mapId}-tactical.png`.
 */

import type { MapCallout } from '../types'

type UECallout = {
  region: string
  super: string
  superName: string
  ueX: number
  ueY: number
}

type MapData = {
  xMul: number
  yMul: number
  xAdd: number
  yAdd: number
  callouts: UECallout[]
}

const c = (region: string, superName: string, ueX: number, ueY: number): UECallout => ({
  region,
  super: superName,
  superName,
  ueX,
  ueY,
})

export const mapData: Record<string, MapData> = {
  ascent: {
    xMul: 7e-5,
    yMul: -7e-5,
    xAdd: 0.813895,
    yAdd: 0.573242,
    callouts: [
      c('Site', 'A', 6153.585, -6626.2114),
      c('Main', 'A', 5321.6206, -4710.1274),
      c('Lobby', 'A', 4489.032, -3014.0515),
      c('Wine', 'A', 7358.7407, -4689.2705),
      c('Garden', 'A', 3773.6653, -7551.3535),
      c('Site', 'B', -2344.065, -7548.511),
      c('Main', 'B', -1983.6713, -5840.8125),
      c('Lobby', 'B', -1490.5864, -1389.9706),
      c('Boat House', 'B', -4484.774, -7763.3584),
      c('Top', 'Mid', 2753.9297, -2129.6155),
      c('Catwalk', 'Mid', 2315.7944, -4127.2554),
      c('Market', 'Mid', 1089.1044, -7363.1914),
      c('Courtyard', 'Mid', 1222.7029, -4586.6),
      c('Spawn TC', 'Attacker', 60, 50),
      c('Spawn PT', 'Defender', 1995.2354, -9744.923),
    ],
  },
  bind: {
    xMul: 5.9e-5,
    yMul: -5.9e-5,
    xAdd: 0.576941,
    yAdd: 0.967566,
    callouts: [
      c('Site', 'A', 10747.902, 2664.4436),
      c('Short', 'A', 7983.3467, 803.96063),
      c('Lobby', 'A', 6113.239, 3158.823),
      c('Bath', 'A', 9106.541, 4449.6587),
      c('Lamps', 'A', 10649.471, 79.904434),
      c('Tower', 'A', 12872.583, 2556.7708),
      c('Site', 'B', 11108.108, -4831.4585),
      c('Long', 'B', 7666.669, -6512.8022),
      c('Short', 'B', 7424.1313, -3056.4531),
      c('Hookah', 'B', 12981.879, -4941.7544),
      c('Garden', 'B', 9144.103, -5598.1274),
      c('Window', 'B', 8826.788, -4309.4116),
      c('Spawn TC', 'Attacker', 161.64832, 77.51108),
      c('Spawn PT', 'Defender', 14641.918, -1017.6743),
    ],
  },
  breeze: {
    xMul: 7e-5,
    yMul: -7e-5,
    xAdd: 0.465123,
    yAdd: 0.833078,
    callouts: [
      c('Site', 'A', 4825, 6325),
      c('Shop', 'A', 2150, 4250),
      c('Lobby', 'A', -1250, 3400),
      c('Pyramids', 'A', 5200, 5450),
      c('Bridge', 'A', 8400, 3525),
      c('Site', 'B', 6450, -5650),
      c('Main', 'B', 3550, -4450),
      c('Window', 'B', 2225, -4175),
      c('Elbow', 'B', 4675, -2900),
      c('Hall', 'Mid', 4256.5713, 2491.0493),
      c('Pillar', 'Mid', 4175, 475),
      c('Top', 'Mid', 6175, 525),
      c('Cannon', 'Mid', 2900, -1850),
      c('Spawn TC', 'Attacker', -575, -450),
      c('Spawn PT', 'Defender', 8900, 3525),
    ],
  },
  fracture: {
    xMul: 7.8e-5,
    yMul: -7.8e-5,
    xAdd: 0.556952,
    yAdd: 1.155886,
    callouts: [
      c('Site', 'A', 8125.7627, 3373.7861),
      c('Main', 'A', 5878.792, 3450.9639),
      c('Hall', 'A', 5063.5464, 2057.6648),
      c('Dish', 'A', 11296.665, 1391.7144),
      c('Drop', 'A', 9306.803, 2826.1626),
      c('Site', 'B', 8178, -5942),
      c('Main', 'B', 5967, -5343),
      c('Arcade', 'B', 10181, -4179),
      c('Tower', 'B', 9155, -5601),
      c('Generator', 'B', 8362, -3380),
      c('Tunnel', 'B', 7402, -4058),
      c('Bridge', 'Attacker', 13204, -756),
      c('Spawn TC', 'Attacker', 4345.554, -948.4505),
      c('Spawn PT', 'Defender', 9156, -677),
    ],
  },
  haven: {
    xMul: 7.5e-5,
    yMul: -7.5e-5,
    xAdd: 1.09345,
    yAdd: 0.642728,
    callouts: [
      c('Site', 'A', 6309.3076, -9225.703),
      c('Lobby', 'A', 3438.537, -6260.409),
      c('Long', 'A', 6209.695, -6901.142),
      c('Garden', 'A', 3100.261, -4683.6016),
      c('Sewer', 'A', 3452.8735, -7915.7246),
      c('Site', 'B', 1884.706, -9231.335),
      c('Back', 'B', 1966.1608, -10664.775),
      c('Site', 'C', -2378.1328, -9010.557),
      c('Long', 'C', -3356.814, -5990.872),
      c('Garage', 'C', 180.07678, -7999.5845),
      c('Window', 'C', -10.126678, -8993.241),
      c('Mid Doors', 'Mid', 151.11594, -6262.9155),
      c('Courtyard', 'Mid', 1822.1299, -6712.6875),
      c('Spawn TC', 'Attacker', 1741.7622, -2642.7925),
      c('Spawn PT', 'Defender', 2946.3042, -12714.707),
    ],
  },
  lotus: {
    xMul: 7.2e-5,
    yMul: -7.2e-5,
    xAdd: 0.454789,
    yAdd: 0.917752,
    callouts: [
      c('Site', 'A', 7735.5396, 5557.309),
      c('Main', 'A', 5288.3022, 4159.762),
      c('Lobby', 'A', 2685.951, 2927.1755),
      c('Tree', 'A', 6149.525, 5557.309),
      c('Drop', 'A', 9516.38, 6092.8936),
      c('Site', 'B', 6368.0327, 668.18317),
      c('Main', 'B', 4876.832, -47.87195),
      c('Pillars', 'B', 3565.3691, 668.18317),
      c('Site', 'C', 6676.6636, -4265.876),
      c('Main', 'C', 5311.2646, -3148.162),
      c('Lobby', 'C', 1403.5685, -1576.5884),
      c('Hall', 'C', 7902.0615, -4265.876),
      c('Waterfall', 'C', 6719.804, -1994.2986),
      c('Spawn TC', 'Attacker', 1401.2915, 777.29834),
      c('Spawn PT', 'Defender', 9686.767, 1697.8223),
    ],
  },
  pearl: {
    xMul: 7.8e-5,
    yMul: -7.8e-5,
    xAdd: 0.480469,
    yAdd: 0.916016,
    callouts: [
      c('Site', 'A', 6613.846, 5569.5254),
      c('Main', 'A', 6368.5713, 3825),
      c('Restaurant', 'A', 4430.452, 2813.1267),
      c('Flowers', 'A', 9263.969, 2507.3403),
      c('Dugout', 'A', 7660.6597, 5854.0664),
      c('Site', 'B', 5800, -2850),
      c('Main', 'B', 4050, -4375),
      c('Hall', 'B', 7495.6177, -4954.14),
      c('Tower', 'B', 8533.423, -2851.3516),
      c('Club', 'B', 800, -1450),
      c('Plaza', 'Mid', 2750, -325),
      c('Shops', 'Mid', 800, -1450),
      c('Top', 'Mid', 2075, 725),
      c('Connector', 'Mid', 6047.0464, 1800.0436),
      c('Spawn TC', 'Attacker', -550, -600),
      c('Spawn PT', 'Defender', 11092.458, 378.79883),
    ],
  },
  split: {
    xMul: 7.8e-5,
    yMul: -7.8e-5,
    xAdd: 0.842188,
    yAdd: 0.697578,
    callouts: [
      c('Site', 'A', 6588.6597, -6761.131),
      c('Main', 'A', 6279.9795, -4492.833),
      c('Lobby', 'A', 6814.217, -2457.7468),
      c('Ramps', 'A', 4330, -4750),
      c('Tower', 'A', 4636.7925, -6748.2334),
      c('Screens', 'A', 5648.7144, -8868.611),
      c('Site', 'B', -2167.2456, -6264.7715),
      c('Main', 'B', -2716.7236, 750.4862),
      c('Lobby', 'B', -1271.6421, -1983.6248),
      c('Garage', 'B', -2190.7827, -3848.0293),
      c('Tower', 'B', 168.89589, -5290.194),
      c('Top', 'Mid', 2021.9575, -4596.936),
      c('Mail', 'Mid', 1155.3333, -4808.6436),
      c('Vent', 'Mid', 3155.1648, -5338.5215),
      c('Spawn TC', 'Attacker', 1901.97, 59.588867),
      c('Spawn PT', 'Defender', 2142.3635, -8964.969),
    ],
  },
  sunset: {
    xMul: 7.8e-5,
    yMul: -7.8e-5,
    xAdd: 0.5,
    yAdd: 0.515625,
    callouts: [
      c('Site', 'A', 1000, 3200),
      c('Main', 'A', -400, 2200),
      c('Lobby', 'A', -1800, 2000),
      c('Elbow', 'A', 200, 4200),
      c('Alley', 'A', 3400, 3600),
      c('Site', 'B', -600, -5850),
      c('Main', 'B', -2000, -5650),
      c('Market', 'B', -200, -3400),
      c('Lobby', 'B', -3400, -2600),
      c('Boba', 'B', 2200, -4800),
      c('Tiles', 'Mid', -1800, 400),
      c('Top', 'Mid', 2000, -2000),
      c('Bottom', 'Mid', -1800, -2025),
      c('Courtyard', 'Mid', -600, -1200),
      c('Spawn TC', 'Attacker', -6025, -400),
      c('Spawn PT', 'Defender', 3805.4785, -1989.0962),
    ],
  },
  abyss: {
    xMul: 8.1e-5,
    yMul: -8.1e-5,
    xAdd: 0.5,
    yAdd: 0.5,
    callouts: [
      c('Site', 'A', 4300, -200),
      c('Main', 'A', 3800, 1650),
      c('Lobby', 'A', 3250, 3400),
      c('Bridge', 'A', 5700, -375),
      c('Tower', 'A', 3025, -125),
      c('Secret', 'A', 3775, -3850),
      c('Site', 'B', -4425, -1175),
      c('Main', 'B', -4450, 1525),
      c('Lobby', 'B', -3650, 4025),
      c('Tower', 'B', -3925, -2500),
      c('Nest', 'B', -4975, 2150),
      c('Library', 'Mid', -325, -600),
      c('Top', 'Mid', 775, -2375),
      c('Catwalk', 'Mid', 600, 525),
      c('Spawn TC', 'Attacker', 950, 4950),
      c('Spawn PT', 'Defender', 950, -5275),
    ],
  },
}

const calloutKind = (region: string, superName: string): MapCallout['kind'] => {
  if (region === 'Site') return 'site'
  if (region.startsWith('Spawn')) return 'spawn'
  if (superName === 'Mid') return 'mid'
  return 'lane'
}

const calloutLabel = (region: string, superName: string): string => {
  if (region === 'Site') return superName // A, B, C
  if (region.startsWith('Spawn')) return region === 'Spawn TC' ? 'SPAWN TẤN CÔNG' : 'SPAWN PHÒNG THỦ'
  if (superName === 'Mid') return `${region.toUpperCase()}`
  if (superName === 'Attacker' || superName === 'Defender') return region.toUpperCase()
  return `${superName} ${region}`.toUpperCase()
}

/** Convert UE coordinates → minimap % theo công thức của Riot. */
export const buildMapCallouts = (mapId: string): MapCallout[] => {
  const data = mapData[mapId]
  if (!data) return []

  return data.callouts.map((item, index) => {
    const x = (item.ueY * data.xMul + data.xAdd) * 100
    const y = (item.ueX * data.yMul + data.yAdd) * 100
    return {
      id: `${mapId}-${index}-${item.super}-${item.region}`.toLowerCase().replace(/\s+/g, '-'),
      label: calloutLabel(item.region, item.superName),
      x: Math.max(0, Math.min(100, Number(x.toFixed(2)))),
      y: Math.max(0, Math.min(100, Number(y.toFixed(2)))),
      kind: calloutKind(item.region, item.superName),
    }
  })
}
