import router from '../router';

export default {
  onMessageLoginReturn(event, context, username) {
    // console.log('Message received');
    const message = JSON.parse(event.data);
    if (message.message === 'OK') {
      context.commit('LOGIN', username);
      router.push({ name: 'Rooms' });
    } else {
      context.commit('LOGINERROR', message.message);
    }
  }
};
