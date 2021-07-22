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
    positions: {},
    upPressed: false,
    downPressed: false,
    rightPressed: false,
    leftPressed: false,
    accelX: 0,
    accelY: 0
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
    },
    SETACCELERATION(state, acceleration) {
      state.accelX = acceleration.x;
      state.accelY = acceleration.y;
    },
    SETUPBUTTONPRESSED(state, pressed) {
      state.upPressed = pressed;
    },
    SETDOWNBUTTONPRESSED(state, pressed) {
      state.downPressed = pressed;
    },
    SETLEFTBUTTONPRESSED(state, pressed) {
      state.leftPressed = pressed;
    },
    SETRIGHTBUTTONPRESSED(state, pressed) {
      state.rightPressed = pressed;
    }
  },
  actions: {
    login(context, username) {
      if (context.getters.connection == null) {
        var conn = new WebSocket(
          'ws://' + window.location.hostname + ':' + location.port + '/api/ws'
        );
        conn.onmessage = event =>
          inputEvents.onMessageLoginReturn(event, context, username);
        conn.onerror = function(event) {
          this.$showLog && console.log('Error');
          this.$showLog && console.log(event);
        };
        conn.onopen = event => {
          this.$showLog && console.log(event);
          this.$showLog &&
            console.log(
              'Successfully connected to the echo websocket server...'
            );
          conn.send(JSON.stringify({ username: username }));
        };
        conn.onclose = event => {
          this.$showLog && console.log(event);
          this.$showLog && console.log('Connection closed');
          this.$showLog && context.commit('LOGOUT');
          router.push({ name: 'Home' });
        };
        context.commit('ADDCONNECTION', conn);
      } else {
        context.getters.connection.onmessage = event =>
          inputEvents.onMessageLoginReturn(event, context, username);
        context.getters.connection.send(JSON.stringify({ username: username }));
      }
    },
    async retrieveListRooms(context) {
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
        context.state.username !== ''
      ) {
        context.getters.connection.onmessage = event =>
          inputEvents.onMessageCreateJoinRoomEvent(event, context, roomName);
        context.getters.connection.send(
          JSON.stringify({ action: 'CREATEROOM', message: roomName })
        );
      }
    },
    leaveRoom(context) {
      if (
        context.state.actualRoom !== '' &&
        context.state.connection != null &&
        context.state.username !== ''
      ) {
        context.getters.connection.send(
          JSON.stringify({ action: 'LEAVEROOM' })
        );
      }
    },
    joinRoom(context, roomName) {
      if (
        roomName !== '' &&
        context.state.connection != null &&
        context.state.username !== ''
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
    async changeButtonState(context, { button, isPressed }) {
      const calcAccel = () => {
        const acceleration = { x: 0, y: 0 };
        if (context.getters.buttonsPressed.up) {
          acceleration.y += 1;
        }
        if (context.getters.buttonsPressed.down) {
          acceleration.y -= 1;
        }
        if (context.getters.buttonsPressed.left) {
          acceleration.x -= 1;
        }
        if (context.getters.buttonsPressed.right) {
          acceleration.x += 1;
        }
        const modul = Math.sqrt(
          Math.pow(acceleration.x, 2) + Math.pow(acceleration.y, 2)
        );
        if (modul > 1) {
          acceleration.x /= modul;
          acceleration.y /= modul;
        }
        context.commit('SETACCELERATION', acceleration);
        context.state.connection.send(
          JSON.stringify({
            action: 'MOVEMENT',
            message: {
              a_x: acceleration.x,
              a_y: acceleration.y
            }
          })
        );
      };
      if (
        context.state.connection != null &&
        context.state.username !== '' &&
        context.state.actualRoom !== ''
      ) {
        switch (button) {
          case 'up':
            if (context.getters.buttonsPressed.up !== isPressed) {
              context.commit('SETUPBUTTONPRESSED', isPressed);
              calcAccel();
            }
            break;
          case 'down':
            if (context.getters.buttonsPressed.down !== isPressed) {
              context.commit('SETDOWNBUTTONPRESSED', isPressed);
              calcAccel();
            }
            break;
          case 'left':
            if (context.getters.buttonsPressed.left !== isPressed) {
              context.commit('SETLEFTBUTTONPRESSED', isPressed);
              calcAccel();
            }
            break;
          case 'right':
            if (context.getters.buttonsPressed.right !== isPressed) {
              context.commit('SETRIGHTBUTTONPRESSED', isPressed);
              calcAccel();
            }
            break;
        }
      }
    }
  },
  getters: {
    username(state) {
      return state.username;
    },
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
    },
    accelaration(state) {
      return { x: state.accelX, y: state.accelY };
    },
    buttonsPressed(state) {
      return {
        up: state.upPressed,
        down: state.downPressed,
        right: state.rightPressed,
        left: state.leftPressed
      };
    }
  },
  modules: {}
});
