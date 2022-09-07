import { atom } from 'recoil'

export const moonsState = atom({
  key: 'moonsState', // unique ID (with respect to other atoms/selectors)
  default: {trammel:0, felucca:0}, // default value (aka initial value)
})

export const windState = atom({
  key: 'windState', 
  default: null,
})

export const boardState = atom({
  key: 'boardState', 
  default: [],
})

export const gemPeerState = atom({
  key: 'gemPeerState', 
  default: null, 
})

export const messagesState = atom({
  key: 'messagesState', 
  default: [], 
}) 

export const promptState = atom({
  key: 'promptState', 
  default: 'default',
}) 

export const extraBoardClassesState = atom({
  key: 'extraBoardClassesState', 
  default: [],
}) 

export const nameState = atom({
  key: 'nameState', 
  default: '',
}) 

export const foodState = atom({
  key: 'foodState', 
  default: 0,
})

export const healthState = atom({
  key: 'healthState', 
  default: 0,
})

export const gemState = atom({
  key: 'gemState', 
  default: 0,
})

export const goldState = atom({
  key: 'goldState', 
  default: 0,
})

export const silverState = atom({
  key: 'silverState', 
  default: 0,
})