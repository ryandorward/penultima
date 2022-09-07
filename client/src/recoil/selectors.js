import { selector } from 'recoil'
import { messagesState } from './atoms'

export const getLastMessage = selector({
  key: 'getLastMessage1',
  get: ({get}) => {
    // const messages = get(messagesState)   
    return get(messagesState)?.slice(-1)[0]
  },
});