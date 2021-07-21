import router from '../router';

export default {
  onMessageLoginReturn(event, context, username) {
    const message = JSON.parse(event.data);
    if (message.message === 'OK') {
      context.commit('LOGIN', username);
      context.dispatch('getRooms');
      router.push({ name: 'Rooms' });
    } else {
      context.commit('LOGINERROR', message.message);
    }
  },
  onMessageGetRoomsEvent(event, context) {
    const message = JSON.parse(event.data);
    context.commit('SETROOMS', message);
  },
  onMessageCreateJoinRoomEvent(event, context, roomName) {
    const message = JSON.parse(event.data);
    if (message.message === 'OK') {
      context.dispatch('joinedRoom', roomName);
      router.push({ name: 'Game' });
    } else {
      context.commit('LOGINERROR', message.message);
    }
  },
  onMessagePositionEvent(event, context) {
    const message = JSON.parse(event.data);
    const p = {};
    var messageValue;
    switch (message.action) {
      case 'MOVES':
        messageValue = JSON.parse(message.message);
        for (const position of messageValue) {
          p[position.player] = position;
        }
        context.commit('SETPOSITIONS', p);
        break;
      case 'LEAVEROOMRESPONSE':
        messageValue = message.message;
        if (messageValue.username === context.getters.username) {
          context.commit('LEAVEROOM');
          router.push({ name: 'Rooms' });
        }
        break;
      default:
        console.log(message);
    }
  }
};
