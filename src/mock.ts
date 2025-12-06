// src/mock.ts
import { Ship } from './types'

const mock: Ship[] = [
  {
    ship_id: 1,
    name: 'Ever Ace',
    description: 'самый большой в мире, двигатель Wartsila 70950 кВт',
    capacity: 23992,
    length: 400,
    width: 61.53,
    draft: 17.0,
    cranes: 6,
    containers: 11996,
    photo_url: 'ever-ace.png'  
  },
  {
    ship_id: 2,
    name: 'FESCO Diomid',
    description: 'построен в 2010 г., судно класса Ice1 (для Арктики), дизельный двигатель, используется на Северном морском пути',
    capacity: 3108,
    length: 195,
    width: 32.2,
    draft: 11.0,
    cranes: 3,
    containers: 536,
    photo_url: 'fesco-diomid.png'  
  },
  {
    ship_id: 3,
    name: 'HMM Algeciras',
    description: 'двигатель MAN B&W 11G95ME-C9.5 мощностью 64 000 кВт, двойные двигатели, система рекуперации энергии, класс DNV GL',
    capacity: 23964,
    length: 399.9,
    width: 61.0,
    draft: 16.5,
    cranes: 7,
    containers: 11982,
    photo_url: 'hmm-algeciras.png'  
  },
  {
    ship_id: 4,
    name: 'MSC Gulsun',
    description: 'первый в мире контейнеровоз, вмещающий более 23 000 TEU, двигатель MAN B&W 11G95ME-C9.5, класс DNV GL',
    capacity: 23756,
    length: 399.9,
    width: 61.4,
    draft: 16.0,
    cranes: 7,
    containers: 11878,
    photo_url: 'msc-gulsun.png'  
  }
]

export default mock