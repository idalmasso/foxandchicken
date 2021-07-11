import { createStore } from 'vuex';
import inputEvents from './inputEvents';

export default createStore({
  state: {
    username: null,
    connection: null,
    error: ''
  },
  mutations: {
    ADDCONNECTION(state, conn) {
      state.connection = conn;
    },
    LOGIN(state, username) {
      state.username = username;
      state.error = '';
    },
    LOGOUT(state) {
      state.connection = null;
      state.username = '';
    },
    LOGINERROR(state, err) {
      state.error = err;
    }
  },
  actions: {
    login(context, username) {
      if (context.getters.connection == null) {
        var conn = new WebSocket('ws://localhost:3000/api/ws');
        conn.onmessage = event =>
          inputEvents.onMessageLoginReturn(event, context, username);
        conn.onerror = function(event) {
          console.log('Error');
          console.log(event);
        };
        conn.onopen = event => {
          console.log(event);
          console.log('Successfully connected to the echo websocket server...');
          conn.send(JSON.stringify({ username: username }));
        };
        conn.onclose = event => {
          console.log(event);
          console.log('Connection closed');
          context.commit('LOGOUT');
        };
        context.commit('ADDCONNECTION', conn);
      } else {
        conn.onmessage = event =>
          inputEvents.onMessageLoginReturn(event, context, username);
        context.getters.connection.send(JSON.stringify({ username: username }));
      }
    }
  },
  getters: {
    connection(state) {
      return state.connection;
    }
  },
  modules: {}
});
