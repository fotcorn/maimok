<template>
  <v-data-table :headers="headers" :items="vms" :items-per-page="100">
    <template v-slot:item.name="{ item }">
      <v-icon small color="success" v-if="item.running">fa-play</v-icon>
      <v-icon small color="default" v-else>fa-pause</v-icon>
      <span class="ml-2">{{ item.name }}</span>
    </template>
    <template v-slot:item.memory="{ item }">
      {{ formatMemory(item.memory) }} GB
    </template>
  </v-data-table>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";

@Component
export default class VMList extends Vue {
  readonly headers = [
    { text: "Name", value: "name" },
    { text: "Memory", value: "memory" }
  ];

  vms = [];
  async mounted() {
    const response = await fetch("http://localhost:7000/vms");
    this.vms = await response.json();
  }

  formatMemory(memory: number) {
    return +(memory / 1024 / 1024).toFixed(2);
  }
}
</script>
