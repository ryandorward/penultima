export const catchDirection = ({e, directionCallback}) => {
  let move = false
  let dirLabel
  // console.log(e.key)
  switch (e.key) {
    case 'ArrowUp':
      dirLabel = 'North'
      move = {x: 0, y: -1}         
      break;
    case 'ArrowDown':      
      dirLabel = 'South' 
      move = {x: 0, y: 1}       
      break;
    case 'ArrowLeft':
      dirLabel = 'West'
      move = {x: -1, y: 0}       
      break;
    case 'ArrowRight':
      dirLabel = 'East'
      move = {x: 1, y: 0}       
      break;
  }
  directionCallback && move && directionCallback({dir: move, dirLabel: dirLabel})   
  return move  
}

export const messageDelay = 300