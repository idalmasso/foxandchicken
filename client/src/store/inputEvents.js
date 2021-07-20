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
    if (message.action === 'MOVES') {
      const p = {};
      const positions = JSON.parse(message.message);
      // debugger;
      for (const position of positions) {
        p[position.player] = position;
      }
      context.commit('SETPOSITIONS', p);
    } else {
      console.log(message);
    }
  }
};
