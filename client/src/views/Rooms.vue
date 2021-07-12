<template>
  <div class="rooms">
    <h1>This is the rooms page</h1>
    <text-adder :textRequest="'New room name'" @text-added="createRoomRequest"></text-adder>
    <table class="list-rooms-table">
      <tr>
        <th>Room</th>
        <th>#Players</th>
        <th>Join</th>
      </tr>
      <tr v-for="room in rooms" :key="room.name">
        <td> {{ room.name }}</td><td>{{ room.players }}</td><td><button type="submit" @click="joinRoom(room.name)">JOIN</button></td>
        </tr>
    </table>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import TextAdder from '../components/ui/TextAdder.vue';
export default {
  components: { TextAdder },
  computed: {
    ...mapGetters({
      rooms: 'getRooms'
    })
  },
  methods: {
    ...mapActions({
      addRoom: 'addRoom',
      joinRoom: 'joinRoom'
    }),
    createRoomRequest(text) {
      this.addRoom(text);
    }
  }

};
</script>

<style scoped>
button {
  margin-top: 1rem;
  width: 100%;
  border-radius: 8px;
  background-color: darksalmon;
  padding: 8px;
}
</style>
