import { createStore } from 'vuex';
import router from '../router';
import inputEvents from './inputEvents';

export default createStore({
  state: {
    username: null,
    connection: null,
    actualRoom: '',
    error: '',
    rooms: [],
    positions: null
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
    },
    SETROOMS(state, rooms) {
      state.rooms = rooms;
    },
    JOINEDROOM(state, room) {
      state.actualRoom = room;
    },
    LEAVEROOM(state) {
      state.actualRoom = '';
    },
    SETPOSITIONS(state, positions) {
      state.positions = positions;
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
          router.push({ name: 'Home' });
        };
        context.commit('ADDCONNECTION', conn);
      } else {
        context.getters.connection.onmessage = event =>
          inputEvents.onMessageLoginReturn(event, context, username);
        context.getters.connection.send(JSON.stringify({ username: username }));
      }
    },
    async getRooms(context) {
      if (context.getters.connection != null) {
        context.getters.connection.onmessage = event =>
          inputEvents.onMessageGetRoomsEvent(event, context);
        context.getters.connection.send(
          JSON.stringify({ action: 'LISTROOMS' })
        );
      }
    },
    addRoom(context, roomName) {
      if (
        roomName !== '' &&
        context.state.connection != null &&
        context.state.connection.username !== ''
      ) {
        context.getters.connection.onmessage = event =>
          inputEvents.onMessageCreateJoinRoomEvent(event, context, roomName);
        context.getters.connection.send(
          JSON.stringify({ action: 'CREATEROOM', message: roomName })
        );
      }
    },
    joinRoom(context, roomName) {
      if (
        roomName !== '' &&
        context.state.connection != null &&
        context.state.connection.username !== ''
      ) {
        context.getters.connection.onmessage = event =>
          inputEvents.onMessageCreateJoinRoomEvent(event, context, roomName);
        context.getters.connection.send(
          JSON.stringify({ action: 'JOINROOM', message: roomName })
        );
      }
    },
    joinedRoom(context, roomName) {
      if (
        roomName !== '' &&
        context.state.connection != null &&
        context.state.username !== ''
      ) {
        context.commit('JOINEDROOM', roomName);
        context.getters.connection.onmessage = event =>
          inputEvents.onMessagePositionEvent(event, context);
      }
    },
    setAcceleration(context, { accelX, accelY }) {
      if (
        context.state.connection != null &&
        context.state.username !== '' &&
        context.state.actualRoom !== ''
      ) {
        context.state.connection.send(
          JSON.stringify({
            action: 'MOVEMENT',
            message: {
              a_x: accelX,
              a_y: accelY
            }
          })
        );
      }
    }
  },
  getters: {
    connection(state) {
      return state.connection;
    },
    getRooms(state) {
      return state.rooms;
    },
    authenticated(state) {
      return state.username !== null && state.username !== '';
    },
    positions(state) {
      return state.positions;
    }
  },
  modules: {}
});
