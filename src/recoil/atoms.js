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
  default: [], 
})

export const messagesState = atom({
  key: 'messagesState', 
  default: [], 
}) 